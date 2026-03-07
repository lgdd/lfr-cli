// Package metadata fetches Liferay release information (bundle URLs, Docker
// images, target platform versions) from lgdd/liferay-product-info on GitHub
// and exposes it as WorkspaceData for use by the scaffold package. Hardcoded
// offline fallback values are used when GitHub is unreachable.
package metadata

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// WorkspaceData holds the release metadata injected into workspace templates.
type WorkspaceData struct {
	Edition              string
	Product              string
	BundleUrl            string
	GithubBundleUrl      string
	TomcatVersion        string
	TargetPlatform       string
	DockerImage          string
	GradleWrapperVersion string
	GroupId              string
	ArtifactId           string
	Name                 string
}

// Release represents a single Liferay product release entry from the releases JSON feed.
type Release struct {
	Product               string            `json:"product"`
	ProductGroupVersion   string            `json:"productGroupVersion"`
	ProductVersion        string            `json:"productVersion"`
	Promoted              string            `json:"promoted"`
	ReleaseKey            string            `json:"releaseKey"`
	TargetPlatformVersion string            `json:"targetPlatformVersion"`
	URL                   string            `json:"url"`
	ReleaseProperties     ReleaseProperties `json:"releaseProperties"`
}

// ReleaseProperties contains the detailed properties of a Liferay release,
// including bundle URLs, Docker image tags, and version identifiers.
type ReleaseProperties struct {
	URL                    string `json:"url"`
	AppServerTomcatVersion string `json:"appServerTomcatVersion"`
	BuildTimestamp         string `json:"buildTimestamp"`
	BundleChecksumSha512   string `json:"bundleChecksumSha512"`
	BundleURL              string `json:"bundleURL"`
	GithubBundleURL        string `json:"githubBundleURL"`
	GitHashLiferayDocker   string `json:"gitHashLiferayDocker"`
	GitHasLiferayPortalEE  string `json:"gitHashLiferayPortalEE"`
	LiferayDockerImage     string `json:"liferayDockerImage"`
	LiferayDockerTags      string `json:"liferayDockerTags"`
	LiferayProductVersion  string `json:"liferayProductVersion"`
	ReleaseDate            string `json:"releaseDate"`
	TargetPlatformVersion  string `json:"targetPlatformVersion"`
}

// QuarterlyRelease represents a single entry from the top-level releases.json feed
// used for DXP quarterly releases.
type QuarterlyRelease struct {
	Product               string   `json:"product"`
	ProductGroupVersion   string   `json:"productGroupVersion"`
	Promoted              string   `json:"promoted"`
	ReleaseKey            string   `json:"releaseKey"`
	TargetPlatformVersion string   `json:"targetPlatformVersion"`
	Tags                  []string `json:"tags"`
	URL                   string   `json:"url"`
}

// Build tool and edition options for a Liferay workspace.
const (
	// Gradle is the build tool identifier for a Gradle workspace.
	Gradle = "gradle"
	// Maven is the build tool identifier for a Maven workspace.
	Maven = "maven"
	// DXP is the edition identifier for Liferay DXP (commercial).
	DXP = "dxp"
	// Portal is the edition identifier for Liferay Portal (community).
	Portal = "portal"
)

// ErrUnkownEdition is returned when the provided edition is neither "dxp" nor "portal".
var ErrUnkownEdition = errors.New("unknown edition (it should be 'dxp' or 'portal')")

// ErrUnsupportedVersion is returned when the provided Liferay version is not supported.
var ErrUnsupportedVersion = errors.New("invalid or unsupported Liferay version")

// PackageName is the Java package name used for generated modules, defaulting to "org.acme".
var PackageName string

// IsQuarterlyVersion reports whether version is a DXP quarterly string, e.g. "2023.q3".
func IsQuarterlyVersion(version string) bool {
	matched, _ := regexp.MatchString(`^\d{4}\.q[1-4]$`, version)
	return matched
}

// IsPortalGAVersion reports whether version is a specific portal GA string, e.g. "7.4.3.112-ga112".
func IsPortalGAVersion(version string) bool {
	matched, _ := regexp.MatchString(`^7\.4\.3\.\d+-ga\d+$`, version)
	return matched
}

// GetGroupId returns the group ID (base Java package name) defined in the
// current Liferay workspace's build file (build.gradle or pom.xml).
func GetGroupId() (string, error) {
	workspacePath, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return "", err
	}
	if fileutil.IsMavenWorkspace(workspacePath) {
		pomParentPath := filepath.Join(workspacePath, "pom.xml")
		pomParent, err := os.Open(pomParentPath)

		if err != nil {
			return "", err
		}

		byteValue, _ := io.ReadAll(pomParent)

		var pom fileutil.Pom
		err = xml.Unmarshal(byteValue, &pom)

		if err != nil {
			return "", err
		}

		if pom.GroupId == "" {
			PackageName = "org.acme"
		} else {
			PackageName = strcase.ToDelimited(pom.GroupId, '.')
		}
	} else if fileutil.IsGradleWorkspace(workspacePath) {
		file, err := os.Open(filepath.Join(workspacePath, "build.gradle"))

		if err != nil {
			return "", err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "group") {
				groupId := strings.Split(scanner.Text(), "=")[1]
				groupId = strings.TrimSpace(groupId)
				groupId = strings.ReplaceAll(groupId, "'", "")
				groupId = strings.ReplaceAll(groupId, "\"", "")
				PackageName = groupId
				break
			}
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}
	} else {
		return "", errors.New("unknown build tool used for this workspace")
	}
	return PackageName, nil
}

// NewWorkspaceData returns the WorkspaceData for the given project base name,
// Liferay version, and edition by fetching the latest release metadata from
// GitHub. It falls back to hardcoded offline values when GitHub is unreachable.
func NewWorkspaceData(base, version, edition string) (*WorkspaceData, error) {
	// DXP quarterly releases (e.g. "2023.q3", "2025.q1")
	if edition == DXP && IsQuarterlyVersion(version) {
		data, err := fetchQuarterlyDXPWorkspaceData(base, version)
		if err != nil {
			if err != ErrUnsupportedVersion {
				logger.PrintWarn(err.Error())
			}
			return getOfflineDXPQuarterlyData(base, version)
		}
		return data, nil
	}

	// Specific portal GA releases (e.g. "7.4.3.112-ga112")
	if edition == Portal && IsPortalGAVersion(version) {
		data, err := fetchPortalGAWorkspaceData(base, version)
		if err != nil {
			logger.PrintWarn(err.Error())
			return getOfflinePortalGAData(base, version)
		}
		return data, nil
	}

	// workaround timeout on this release
	if edition == Portal && version == "7.0" {
		return getOfflineWorkspaceData(base, version, edition)
	}

	releases, err := fetchReleases(version, edition)

	if err != nil {
		if err == ErrUnkownEdition || err == ErrUnsupportedVersion {
			return nil, err
		} else {
			logger.PrintWarn(err.Error())
			return getOfflineWorkspaceData(base, version, edition)
		}
	}

	latestRelease := releases[0]

	gradleWrapperVersion := "7.6.4"

	// workaround issue in releases.json for 7.1
	if edition == Portal && version == "7.1" {
		latestRelease = releases[1]
	}

	if version == "7.4" {
		gradleWrapperVersion = "8.5"
	}

	latestRelease.BuildGithubBundleURL()

	return &WorkspaceData{
		Edition:              edition,
		Product:              latestRelease.ReleaseKey,
		BundleUrl:            latestRelease.ReleaseProperties.BundleURL,
		GithubBundleUrl:      latestRelease.ReleaseProperties.GithubBundleURL,
		TargetPlatform:       latestRelease.TargetPlatformVersion,
		DockerImage:          latestRelease.ReleaseProperties.LiferayDockerImage,
		GradleWrapperVersion: gradleWrapperVersion,
		GroupId:              strcase.ToDelimited(PackageName, '.'),
		ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
		Name:                 strcase.ToCamel(strings.ToLower(base)),
	}, nil
}

func fetchReleases(version, edition string) ([]Release, error) {
	var releasesURLBuilder strings.Builder
	releasesURLBuilder.WriteString("https://raw.githubusercontent.com/lgdd/liferay-product-info/main/releases/")
	releaseVersion := strings.ReplaceAll(version, ".", "")

	if edition != DXP && edition != Portal {
		return []Release{}, ErrUnkownEdition
	}

	if version != "7.4" && version != "7.3" && version != "7.2" && version != "7.1" && version != "7.0" {
		return []Release{}, ErrUnsupportedVersion
	}

	releasesURLBuilder.WriteString(edition)
	releasesURLBuilder.WriteString("_")
	releasesURLBuilder.WriteString(releaseVersion)
	releasesURLBuilder.WriteString("_releases.json")

	resp, err := http.Get(releasesURLBuilder.String())

	if err != nil {
		return []Release{}, err
	}

	defer resp.Body.Close()
	var releases []Release
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &releases)

	if err != nil {
		return []Release{}, err
	}

	return releases, nil
}

func (release *Release) BuildGithubBundleURL() {
	var githubBundleURLBuilder strings.Builder
	githubBaseURL := "https://github.com/lgdd/liferay-dxp-releases/releases/download/"
	githubReleaseName := release.ReleaseKey

	if release.Product == "portal" {
		githubBaseURL = "https://github.com/liferay/liferay-portal/releases/download/"
		githubReleaseName = strings.Split(release.ReleaseProperties.LiferayDockerImage, ":")[1]
	}

	bundleURLSplit := strings.Split(release.ReleaseProperties.BundleURL, "/")
	bundleName := bundleURLSplit[len(bundleURLSplit)-1]

	githubBundleURLBuilder.WriteString(githubBaseURL)
	githubBundleURLBuilder.WriteString(githubReleaseName)
	githubBundleURLBuilder.WriteString("/")
	githubBundleURLBuilder.WriteString(bundleName)

	release.ReleaseProperties.GithubBundleURL = githubBundleURLBuilder.String()
}

// fetchQuarterlyDXPWorkspaceData fetches workspace data for a DXP quarterly release
// (e.g. "2024.q1") by loading releases.json from ~/.lfr/ (with a background refresh)
// or fetching it from GitHub on first use.
func fetchQuarterlyDXPWorkspaceData(base, quarter string) (*WorkspaceData, error) {
	allReleases, err := loadOrFetchReleasesJSON()
	if err != nil {
		return nil, err
	}

	var candidates []QuarterlyRelease
	for _, r := range allReleases {
		if r.Product == "dxp" && r.ProductGroupVersion == quarter {
			candidates = append(candidates, r)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrUnsupportedVersion
	}

	best := pickBestRelease(candidates)

	props, err := fetchReleaseProperties(best.URL)
	if err != nil {
		return nil, err
	}

	return buildWorkspaceData(base, best.ReleaseKey, DXP, props, "8.5"), nil
}

// loadOrFetchReleasesJSON loads the quarterly releases index from ~/.lfr/releases.json
// if it already exists (triggering a background refresh), or fetches it from GitHub
// and saves it on first use.
func loadOrFetchReleasesJSON() ([]QuarterlyRelease, error) {
	releasesPath := filepath.Join(conf.GetConfigPath(), "releases.json")

	if _, err := os.Stat(releasesPath); err == nil {
		releases, err := readReleasesJSONFile(releasesPath)
		if err == nil {
			go refreshReleasesJSON(releasesPath)
			return releases, nil
		}
	}

	return fetchAndWriteReleasesJSON(releasesPath)
}

func fetchAndWriteReleasesJSON(path string) ([]QuarterlyRelease, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/lgdd/liferay-product-info/main/releases.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []QuarterlyRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	_ = os.WriteFile(path, body, 0644)

	return releases, nil
}

func readReleasesJSONFile(path string) ([]QuarterlyRelease, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var releases []QuarterlyRelease
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil, err
	}
	return releases, nil
}

func refreshReleasesJSON(path string) {
	_, _ = fetchAndWriteReleasesJSON(path)
}

// fetchPortalGAWorkspaceData fetches workspace data for a specific portal GA release
// (e.g. "7.4.3.112-ga112") from the Liferay CDN.
func fetchPortalGAWorkspaceData(base, version string) (*WorkspaceData, error) {
	cdnURL := "https://releases-cdn.liferay.com/portal/" + version
	props, err := fetchReleaseProperties(cdnURL)
	if err != nil {
		return nil, err
	}
	releaseKey := "portal-7.4-ga" + extractGANumber(version)
	return buildWorkspaceData(base, releaseKey, Portal, props, "8.5"), nil
}

// fetchReleaseProperties fetches and parses the release.properties file from a CDN base URL.
func fetchReleaseProperties(cdnURL string) (map[string]string, error) {
	resp, err := http.Get(cdnURL + "/release.properties")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s/release.properties", resp.StatusCode, cdnURL)
	}

	props := make(map[string]string)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		props[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if props["bundle.url"] == "" {
		return nil, fmt.Errorf("no bundle.url in release.properties at %s", cdnURL)
	}
	return props, nil
}

// buildWorkspaceData constructs a WorkspaceData from a props map and metadata.
func buildWorkspaceData(base, releaseKey, edition string, props map[string]string, gradleWrapperVersion string) *WorkspaceData {
	bundleURL := props["bundle.url"]
	dockerImage := props["liferay.docker.image"]
	target := props["target.platform.version"]
	if target == "" {
		target = deriveTargetPlatform(releaseKey)
	}
	return &WorkspaceData{
		Edition:              edition,
		Product:              releaseKey,
		BundleUrl:            bundleURL,
		GithubBundleUrl:      buildGithubBundleURL(edition, releaseKey, bundleURL),
		TargetPlatform:       target,
		DockerImage:          dockerImage,
		GradleWrapperVersion: gradleWrapperVersion,
		GroupId:              strcase.ToDelimited(PackageName, '.'),
		ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
		Name:                 strcase.ToCamel(strings.ToLower(base)),
	}
}

// buildGithubBundleURL constructs the GitHub mirror bundle URL for a release.
func buildGithubBundleURL(edition, releaseKey, bundleURL string) string {
	parts := strings.Split(bundleURL, "/")
	bundleName := parts[len(parts)-1]

	if edition == DXP {
		return "https://github.com/lgdd/liferay-dxp-releases/releases/download/" + releaseKey + "/" + bundleName
	}
	// Portal: version tag is the second-to-last segment of the CDN bundle URL path.
	versionTag := parts[len(parts)-2]
	return "https://github.com/liferay/liferay-portal/releases/download/" + versionTag + "/" + bundleName
}

// pickBestRelease selects the best release from candidates, preferring "recommended"
// tag, then promoted == "true", falling back to the first entry.
func pickBestRelease(candidates []QuarterlyRelease) QuarterlyRelease {
	for _, r := range candidates {
		for _, tag := range r.Tags {
			if tag == "recommended" {
				return r
			}
		}
	}
	for _, r := range candidates {
		if r.Promoted == "true" {
			return r
		}
	}
	return candidates[0]
}

// extractGANumber extracts the numeric GA suffix from a portal version string.
// e.g. "7.4.3.112-ga112" → "112"
func extractGANumber(version string) string {
	parts := strings.Split(version, "-ga")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// deriveTargetPlatform derives the target platform version from a release key
// when target.platform.version is absent in release.properties.
// e.g. "portal-7.4-ga108" → "7.4.3.108"
func deriveTargetPlatform(releaseKey string) string {
	const prefix = "portal-7.4-ga"
	if strings.HasPrefix(releaseKey, prefix) {
		n := strings.TrimPrefix(releaseKey, prefix)
		return "7.4.3." + n
	}
	return ""
}

// getOfflineDXPQuarterlyData returns hardcoded WorkspaceData for a DXP quarterly version
// when the live releases.json is unreachable.
func getOfflineDXPQuarterlyData(base, quarter string) (*WorkspaceData, error) {
	type entry struct {
		releaseKey  string
		bundleURL   string
		dockerImage string
		target      string
	}

	offlineData := map[string]entry{
		"2023.q3": {
			releaseKey:  "dxp-2023.q3.10",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2023.q3.10/liferay-dxp-tomcat-2023.q3.10-1740439074.7z",
			dockerImage: "liferay/dxp:2023.q3.10",
			target:      "2023.q3.10",
		},
		"2023.q4": {
			releaseKey:  "dxp-2023.q4.10",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2023.q4.10/liferay-dxp-tomcat-2023.q4.10-1737579304.7z",
			dockerImage: "liferay/dxp:2023.q4.10",
			target:      "2023.q4.10",
		},
		"2024.q1": {
			releaseKey:  "dxp-2024.q1.25",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2024.q1.25/liferay-dxp-tomcat-2024.q1.25-1770226216.7z",
			dockerImage: "liferay/dxp:2024.q1.25",
			target:      "2024.q1.25",
		},
		"2024.q2": {
			releaseKey:  "dxp-2024.q2.13",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2024.q2.13/liferay-dxp-tomcat-2024.q2.13-1739794596.7z",
			dockerImage: "liferay/dxp:2024.q2.13",
			target:      "2024.q2.13",
		},
		"2024.q3": {
			releaseKey:  "dxp-2024.q3.13",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2024.q3.13/liferay-dxp-tomcat-2024.q3.13-1734359202.7z",
			dockerImage: "liferay/dxp:2024.q3.13",
			target:      "2024.q3.13",
		},
		"2024.q4": {
			releaseKey:  "dxp-2024.q4.7",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2024.q4.7/liferay-dxp-tomcat-2024.q4.7-1739190803.7z",
			dockerImage: "liferay/dxp:2024.q4.7",
			target:      "2024.q4.7",
		},
		"2025.q1": {
			releaseKey:  "dxp-2025.q1.20-lts",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2025.q1.20-lts/liferay-dxp-tomcat-2025.q1.20-lts-1765372184.7z",
			dockerImage: "liferay/dxp:2025.q1.20-lts",
			target:      "2025.q1.20",
		},
		"2025.q2": {
			releaseKey:  "dxp-2025.q2.12",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2025.q2.12/liferay-dxp-tomcat-2025.q2.12-1756092679.7z",
			dockerImage: "liferay/dxp:2025.q2.12",
			target:      "2025.q2.12",
		},
		"2025.q3": {
			releaseKey:  "dxp-2025.q3.10",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2025.q3.10/liferay-dxp-tomcat-2025.q3.10-1763478154.7z",
			dockerImage: "liferay/dxp:2025.q3.10",
			target:      "2025.q3.10",
		},
		"2025.q4": {
			releaseKey:  "dxp-2025.q4.9",
			bundleURL:   "https://releases-cdn.liferay.com/dxp/2025.q4.9/liferay-dxp-tomcat-2025.q4.9-1771236885.7z",
			dockerImage: "liferay/dxp:2025.q4.9",
			target:      "2025.q4.9",
		},
	}

	d, ok := offlineData[quarter]
	if !ok {
		return nil, ErrUnsupportedVersion
	}

	return &WorkspaceData{
		Edition:              DXP,
		Product:              d.releaseKey,
		BundleUrl:            d.bundleURL,
		GithubBundleUrl:      buildGithubBundleURL(DXP, d.releaseKey, d.bundleURL),
		TargetPlatform:       d.target,
		DockerImage:          d.dockerImage,
		GradleWrapperVersion: "8.5",
		GroupId:              strcase.ToDelimited(PackageName, '.'),
		ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
		Name:                 strcase.ToCamel(strings.ToLower(base)),
	}, nil
}

// getOfflinePortalGAData returns hardcoded WorkspaceData for a specific portal GA version
// when the CDN is unreachable.
func getOfflinePortalGAData(base, version string) (*WorkspaceData, error) {
	type entry struct {
		bundleURL   string
		dockerImage string
		target      string
	}

	offlineData := map[string]entry{
		"7.4.3.98-ga98": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.98-ga98/liferay-ce-portal-tomcat-7.4.3.98-ga98-20231012155046762.7z",
			dockerImage: "liferay/portal:7.4.3.98-ga98",
			target:      "7.4.3.98",
		},
		"7.4.3.99-ga99": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.99-ga99/liferay-ce-portal-tomcat-7.4.3.99-ga99-20231019093906822.7z",
			dockerImage: "liferay/portal:7.4.3.99-ga99",
			target:      "7.4.3.99",
		},
		"7.4.3.100-ga100": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.100-ga100/liferay-ce-portal-tomcat-7.4.3.100-ga100-20231028043213993.7z",
			dockerImage: "liferay/portal:7.4.3.100-ga100",
			target:      "7.4.3.100",
		},
		"7.4.3.101-ga101": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.101-ga101/liferay-ce-portal-tomcat-7.4.3.101-ga101-20231102201111554.7z",
			dockerImage: "liferay/portal:7.4.3.101-ga101",
			target:      "7.4.3.101",
		},
		"7.4.3.102-ga102": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.102-ga102/liferay-ce-portal-tomcat-7.4.3.102-ga102-20231109165213885.7z",
			dockerImage: "liferay/portal:7.4.3.102-ga102",
			target:      "7.4.3.102",
		},
		"7.4.3.103-ga103": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.103-ga103/liferay-ce-portal-tomcat-7.4.3.103-ga103-20231116132758925.7z",
			dockerImage: "liferay/portal:7.4.3.103-ga103",
			target:      "7.4.3.103",
		},
		"7.4.3.104-ga104": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.104-ga104/liferay-ce-portal-tomcat-7.4.3.104-ga104-20231124094833182.7z",
			dockerImage: "liferay/portal:7.4.3.104-ga104",
			target:      "7.4.3.104",
		},
		"7.4.3.105-ga105": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.105-ga105/liferay-ce-portal-tomcat-7.4.3.105-ga105-20231201054529989.7z",
			dockerImage: "liferay/portal:7.4.3.105-ga105",
			target:      "7.4.3.105",
		},
		"7.4.3.106-ga106": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.106-ga106/liferay-ce-portal-tomcat-7.4.3.106-ga106-20231207073813307.7z",
			dockerImage: "liferay/portal:7.4.3.106-ga106",
			target:      "7.4.3.106",
		},
		"7.4.3.107-ga107": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.107-ga107/liferay-ce-portal-tomcat-7.4.3.107-ga107-20231214043229049.7z",
			dockerImage: "liferay/portal:7.4.3.107-ga107",
			target:      "7.4.3.107",
		},
		"7.4.3.108-ga108": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.108-ga108/liferay-ce-portal-tomcat-7.4.3.108-ga108-20231229104558156.7z",
			dockerImage: "liferay/portal:7.4.3.108-ga108",
			target:      "7.4.3.108",
		},
		"7.4.3.109-ga109": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.109-ga109/liferay-ce-portal-tomcat-7.4.3.109-ga109-20240103085835525.7z",
			dockerImage: "liferay/portal:7.4.3.109-ga109",
			target:      "7.4.3.109",
		},
		"7.4.3.112-ga112": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.112-ga112/liferay-ce-portal-tomcat-7.4.3.112-ga112-20240226061339195.7z",
			dockerImage: "liferay/portal:7.4.3.112-ga112",
			target:      "7.4.3.112",
		},
		"7.4.3.120-ga120": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.120-ga120/liferay-portal-tomcat-7.4.3.120-ga120-1718225443.7z",
			dockerImage: "liferay/portal:7.4.3.120-ga120",
			target:      "7.4.3.120",
		},
		"7.4.3.125-ga125": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.125-ga125/liferay-portal-tomcat-7.4.3.125-ga125-1726242956.7z",
			dockerImage: "liferay/portal:7.4.3.125-ga125",
			target:      "7.4.3.125",
		},
		"7.4.3.129-ga129": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.129-ga129/liferay-portal-tomcat-7.4.3.129-ga129-1733783976.7z",
			dockerImage: "liferay/portal:7.4.3.129-ga129",
			target:      "7.4.3.129",
		},
		"7.4.3.132-ga132": {
			bundleURL:   "https://releases-cdn.liferay.com/portal/7.4.3.132-ga132/liferay-portal-tomcat-7.4.3.132-ga132-1739912568.7z",
			dockerImage: "liferay/portal:7.4.3.132-ga132",
			target:      "7.4.3.132",
		},
	}

	d, ok := offlineData[version]
	if !ok {
		return nil, ErrUnsupportedVersion
	}

	releaseKey := "portal-7.4-ga" + extractGANumber(version)
	return &WorkspaceData{
		Edition:              Portal,
		Product:              releaseKey,
		BundleUrl:            d.bundleURL,
		GithubBundleUrl:      buildGithubBundleURL(Portal, releaseKey, d.bundleURL),
		TargetPlatform:       d.target,
		DockerImage:          d.dockerImage,
		GradleWrapperVersion: "8.5",
		GroupId:              strcase.ToDelimited(PackageName, '.'),
		ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
		Name:                 strcase.ToCamel(strings.ToLower(base)),
	}, nil
}

func getOfflineWorkspaceData(base, version, edition string) (*WorkspaceData, error) {
	if edition == DXP {
		switch version {
		case "7.4":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "dxp-2024.q1.5",
				BundleUrl:            "https://releases-cdn.liferay.com/dxp/2024.q1.5/liferay-dxp-tomcat-2024.q1.5-1712566347.7z",
				TargetPlatform:       "2024.q1.5",
				DockerImage:          "liferay/dxp:2024.q1.5",
				GradleWrapperVersion: "8.5",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.3":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "dxp-7.3-u36",
				BundleUrl:            "https://releases-cdn.liferay.com/dxp/7.3.10-u36/liferay-dxp-tomcat-7.3.10-u36-1706652128.7z",
				TargetPlatform:       "7.3.10.u36",
				DockerImage:          "liferay/dxp:7.3.10-u36",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.2":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "dxp-7.2.8",
				BundleUrl:            "https://releases-cdn.liferay.com/dxp/7.2.10.8/liferay-dxp-tomcat-7.2.10.8-sp8-20220912234451782.7z",
				TargetPlatform:       "7.2.10.8",
				DockerImage:          "liferay/dxp:7.2.10-sp8",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.1":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "dxp-7.1-dxp-28",
				BundleUrl:            "https://releases-cdn.liferay.com/dxp/7.1.10-dxp-28/liferay-dxp-tomcat-7.1.10-dxp-28-20220823192814876.7z",
				TargetPlatform:       "7.1.10.8",
				DockerImage:          "liferay/dxp:7.1.10-dxp-28",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.0":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "dxp-7.0.17",
				BundleUrl:            "https://releases-cdn.liferay.com/dxp/7.0.10.17/liferay-dxp-digital-enterprise-tomcat-7.0.10.17-sp17-slim-20211014075354439.7z",
				TargetPlatform:       "7.0.10.17",
				DockerImage:          "liferay/dxp:7.0.10-sp17",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		}
	} else {
		switch version {
		case "7.4":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "portal-7.4-ga112",
				BundleUrl:            "https://github.com/liferay/liferay-portal/releases/download/7.4.3.112-ga112/liferay-ce-portal-tomcat-7.4.3.112-ga112-20240226061339195.7z",
				TargetPlatform:       "7.4.3.112",
				DockerImage:          "liferay/portal:7.4.3.112-ga112",
				GradleWrapperVersion: "8.5",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.3":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "portal-7.3-ga8",
				BundleUrl:            "https://github.com/liferay/liferay-portal/releases/download/7.3.7-ga8/liferay-ce-portal-tomcat-7.3.7-ga8-20210610183559721.7z",
				TargetPlatform:       "7.3.7",
				DockerImage:          "liferay/portal:7.3.7-ga8",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.2":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "portal-7.2-ga2",
				BundleUrl:            "https://github.com/liferay/liferay-portal/releases/download/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.7z",
				TargetPlatform:       "7.2.1-1",
				DockerImage:          "liferay/portal:7.2.1-ga2",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.1":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "portal-7.1-ga4",
				BundleUrl:            "https://github.com/liferay/liferay-portal/releases/download/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.7z",
				TargetPlatform:       "7.1.3-1",
				DockerImage:          "liferay/portal:7.1.3-ga4",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.0":
			return &WorkspaceData{
				Edition:              edition,
				Product:              "portal-7.0-ga7",
				BundleUrl:            "https://releases-cdn.liferay.com/portal/7.0.6-ga7/liferay-ce-portal-tomcat-7.0-ga7-20180507111753223.zip",
				TargetPlatform:       "7.0.6-2",
				DockerImage:          "liferay/portal:7.0.6-ga7",
				GradleWrapperVersion: "7.6.4",
				GroupId:              strcase.ToDelimited(PackageName, '.'),
				ArtifactId:           strcase.ToKebab(strings.ToLower(base)),
				Name:                 strcase.ToCamel(strings.ToLower(base)),
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid or unsupported Liferay version")
}
