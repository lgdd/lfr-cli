package project

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/schollz/progressbar/v3"
)

// Metadata represents the basic informations associated with a Liferay project
type Metadata struct {
	Product        string
	BundleUrl      string
	TomcatVersion  string
	TargetPlatform string
	DockerImage    string
	GroupId        string
	ArtifactId     string
	Name           string
}

type GitHubAsset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}

// Build options
const (
	Gradle = "gradle"
	Maven  = "maven"
)

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

		byteValue, _ := ioutil.ReadAll(pomParent)

		var pom Pom
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
func NewMetadata(base, version string) (*Metadata, error) {
	switch version {
	case "7.4":

		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Fetching latest info from GitHub"),
			progressbar.OptionSpinnerType(11))

		resp, err := http.Get("https://api.github.com/repos/liferay/liferay-portal/releases/latest")

		release := &GitHubRelease{
			TagName: "7.4.3.30-ga30",
			Assets: []GitHubAsset{
				{
					BrowserDownloadURL: "https://github.com/liferay/liferay-portal/releases/download/7.4.3.30-ga30/liferay-ce-portal-tomcat-7.4.3.30-ga30-20220622172832884.tar.gz",
				},
			},
		}
		downloadURL := release.Assets[0].BrowserDownloadURL

		if err != nil {
			bar.Clear()
			printutil.Warning("Can not fetch info from GitHub\n")
			printutil.Warning("Start offline mode process\n\n")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			io.Copy(io.MultiWriter(bar), resp.Body)

			defer resp.Body.Close()

			if err := json.Unmarshal(body, &release); err != nil {
				bar.Clear()
				fmt.Println("Can not unmarshal GitHub release response")
				fmt.Println(body)
				fmt.Println(err.Error())
				fmt.Println("Start offline mode process")
			} else {

				for _, asset := range release.Assets {
					if strings.Contains(asset.BrowserDownloadURL, "tomcat") && strings.Contains(asset.BrowserDownloadURL, "tar.gz") {
						downloadURL = asset.BrowserDownloadURL
						break
					}
				}
			}

		}

		bar.Clear()

		gaUpdateSplit := strings.Split(release.TagName, "-")
		gaUpdate := gaUpdateSplit[len(gaUpdateSplit)-1]

		return &Metadata{
			Product:        strings.Join([]string{"portal", version, gaUpdate}, "-"),
			BundleUrl:      downloadURL,
			TargetPlatform: strings.Split(release.TagName, "-")[0],
			DockerImage:    strings.Join([]string{"liferay/portal", release.TagName}, ":"),
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.3":
		return &Metadata{
			Product:        "portal-7.3-ga8",
			BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.3.7-ga8/liferay-ce-portal-tomcat-7.3.7-ga8-20210610183559721.tar.gz",
			TargetPlatform: "7.3.7",
			DockerImage:    "liferay/portal:7.3.7-ga8",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.2":
		return &Metadata{
			Product:        "portal-7.2-ga2",
			BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.tar.gz",
			TargetPlatform: "7.2.1-1",
			DockerImage:    "liferay/portal:7.2.1-ga2",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.1":
		return &Metadata{
			Product:        "portal-7.1-ga4",
			BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.tar.gz",
			TargetPlatform: "7.1.3-1",
			DockerImage:    "liferay/portal:7.1.3-ga4",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.0":
		return &Metadata{
			Product:        "portal-7.0-ga7",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.0.6-ga7/liferay-ce-portal-tomcat-7.0-ga7-20180507111753223.zip",
			TargetPlatform: "7.0.6-2",
			DockerImage:    "liferay/portal:7.0.6-ga7",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	}
	return nil, fmt.Errorf("invalid Liferay version")
}
