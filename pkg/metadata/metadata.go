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
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
)

// WorkspaceData represents the basic informations associated with a Liferay workspace
type WorkspaceData struct {
	Edition         string
	Product         string
	BundleUrl       string
	GithubBundleUrl string
	TomcatVersion   string
	TargetPlatform  string
	DockerImage     string
	GroupId         string
	ArtifactId      string
	Name            string
}

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

// Build & Edition options
const (
	Gradle = "gradle"
	Maven  = "maven"
	DXP    = "dxp"
	Portal = "portal"
)

// Expected errors
var ErrUnkownEdition = errors.New("unknown edition (it should be 'dxp' or 'portal')")
var ErrUnsupportedVersion = errors.New("invalid or unsupported Liferay version")

// Package name to use for the project, default is org.acme
var PackageName string

// Get the group ID (base package name) associated with the Liferay workspace
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

// Returns metadata for a given project and the chosen Liferay version
func NewWorkspaceData(base, version, edition string) (*WorkspaceData, error) {
	// workaround timeout on this release
	if edition == Portal && version == "7.0" {
		return getOfflineWorkspaceData(base, version, edition)
	}

	releases, err := fetchReleases(version, edition)

	if err != nil {
		if err == ErrUnkownEdition || err == ErrUnsupportedVersion {
			return nil, err
		} else {
			printutil.Warning(fmt.Sprintf("%s\n", err.Error()))
			return getOfflineWorkspaceData(base, version, edition)
		}
	}

	latestRelease := releases[0]

	// workaround issue in releases.json for 7.1
	if edition == Portal && version == "7.1" {
		latestRelease = releases[1]
	}

	latestRelease.BuildGithubBundleURL()

	return &WorkspaceData{
		Edition:         edition,
		Product:         latestRelease.ReleaseKey,
		BundleUrl:       latestRelease.ReleaseProperties.BundleURL,
		GithubBundleUrl: latestRelease.ReleaseProperties.GithubBundleURL,
		TargetPlatform:  latestRelease.ReleaseProperties.TargetPlatformVersion,
		DockerImage:     latestRelease.ReleaseProperties.LiferayDockerImage,
		GroupId:         strcase.ToDelimited(PackageName, '.'),
		ArtifactId:      strcase.ToKebab(strings.ToLower(base)),
		Name:            strcase.ToCamel(strings.ToLower(base)),
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
	bundleURLSplit := strings.Split(release.ReleaseProperties.BundleURL, "/")
	bundleName := bundleURLSplit[len(bundleURLSplit)-1]

	githubBundleURLBuilder.WriteString(githubBaseURL)
	githubBundleURLBuilder.WriteString(release.ReleaseKey)
	githubBundleURLBuilder.WriteString("/")
	githubBundleURLBuilder.WriteString(bundleName)

	release.ReleaseProperties.GithubBundleURL = githubBundleURLBuilder.String()
}

func getOfflineWorkspaceData(base, version, edition string) (*WorkspaceData, error) {
	if edition == DXP {
		switch version {
		case "7.4":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "dxp-2024.q1.5",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/2024.q1.5/liferay-dxp-tomcat-2024.q1.5-1712566347.7z",
				TargetPlatform: "2024.q1.5",
				DockerImage:    "liferay/dxp:2024.q1.5",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.3":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "dxp-7.3-u36",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.3.10-u36/liferay-dxp-tomcat-7.3.10-u36-1706652128.7z",
				TargetPlatform: "7.3.10.u36",
				DockerImage:    "liferay/dxp:7.3.10-u36",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.2":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "dxp-7.2.8",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.2.10.8/liferay-dxp-tomcat-7.2.10.8-sp8-20220912234451782.7z",
				TargetPlatform: "7.2.10.8",
				DockerImage:    "liferay/dxp:7.2.10-sp8",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.1":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "dxp-7.1-dxp-28",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.1.10-dxp-28/liferay-dxp-tomcat-7.1.10-dxp-28-20220823192814876.7z",
				TargetPlatform: "7.1.10.8",
				DockerImage:    "liferay/dxp:7.1.10-dxp-28",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.0":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "dxp-7.0.17",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.0.10.17/liferay-dxp-digital-enterprise-tomcat-7.0.10.17-sp17-slim-20211014075354439.7z",
				TargetPlatform: "7.0.10.17",
				DockerImage:    "liferay/dxp:7.0.10-sp17",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		}
	} else {
		switch version {
		case "7.4":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "portal-7.4-ga112",
				BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.4.3.112-ga112/liferay-ce-portal-tomcat-7.4.3.112-ga112-20240226061339195.7z",
				TargetPlatform: "7.4.3.112",
				DockerImage:    "liferay/portal:7.4.3.112-ga112",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.3":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "portal-7.3-ga8",
				BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.3.7-ga8/liferay-ce-portal-tomcat-7.3.7-ga8-20210610183559721.7z",
				TargetPlatform: "7.3.7",
				DockerImage:    "liferay/portal:7.3.7-ga8",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.2":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "portal-7.2-ga2",
				BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.7z",
				TargetPlatform: "7.2.1-1",
				DockerImage:    "liferay/portal:7.2.1-ga2",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.1":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "portal-7.1-ga4",
				BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.7z",
				TargetPlatform: "7.1.3-1",
				DockerImage:    "liferay/portal:7.1.3-ga4",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.0":
			return &WorkspaceData{
				Edition:        edition,
				Product:        "portal-7.0-ga7",
				BundleUrl:      "https://releases-cdn.liferay.com/portal/7.0.6-ga7/liferay-ce-portal-tomcat-7.0-ga7-20180507111753223.zip",
				TargetPlatform: "7.0.6-2",
				DockerImage:    "liferay/portal:7.0.6-ga7",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid or unsupported Liferay version")
}
