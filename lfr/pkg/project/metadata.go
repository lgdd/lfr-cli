package project

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
	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
	progressbar "github.com/schollz/progressbar/v3"
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

type ProductInfo struct {
	AppServerTomcatVersion string `json:"appServerTomcatVersion"`
	BundleURL              string `json:"bundleUrl"`
	BundleChecksumMD5      string `json:"bundleChecksumMD5"`
	BundleChecksumMD5URL   string `json:"bundleChecksumMD5Url"`
	LiferayDockerImage     string `json:"liferayDockerImage"`
	LiferayProductVersion  string `json:"liferayProductVersion"`
	Promoted               string `json:"promoted"`
	ReleaseDate            string `json:"releaseDate"`
	TargetPlatformVersion  string `json:"targetPlatformVersion"`
	Name                   string `json:"name"`
}

// Build & Edition options
const (
	Gradle = "gradle"
	Maven  = "maven"
	DXP    = "dxp"
	Portal = "portal"
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

		byteValue, _ := io.ReadAll(pomParent)

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
func NewMetadata(base, version, edition string) (*Metadata, error) {
	var productInfoURLBuilder strings.Builder
	productInfoURLBuilder.WriteString("https://raw.githubusercontent.com/lgdd/liferay-product-info/main/")
	productInfoVersion := strings.ReplaceAll(version, ".", "")

	if edition != DXP && edition != Portal {
		return nil, errors.New("unknown edition (it should be 'dxp' or 'portal')")
	}

	if version != "7.4" && version != "7.3" && version != "7.2" && version != "7.1" && version != "7.0" {
		return nil, fmt.Errorf("invalid or unsupported Liferay version")
	}

	productInfoURLBuilder.WriteString(edition)
	productInfoURLBuilder.WriteString("_")
	productInfoURLBuilder.WriteString(productInfoVersion)
	productInfoURLBuilder.WriteString("_product_info.json")

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Fetching latest info from https://github.com/lgdd/liferay-product-info"),
		progressbar.OptionSpinnerType(11))

	resp, err := http.Get(productInfoURLBuilder.String())

	if err != nil {
		bar.Clear()
		printutil.Warning(fmt.Sprintf("%s\n", err.Error()))
		return getOfflineMetadata(base, version, edition)
	}

	var productInfoList []ProductInfo
	body, _ := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err := json.Unmarshal(body, &productInfoList); err != nil {
		bar.Clear()
		fmt.Println("Can not unmarshal response")
		fmt.Println(body)
		fmt.Println(err.Error())
		fmt.Println("Get offline data")
		return getOfflineMetadata(base, version, edition)
	}

	latestProductInfo := productInfoList[len(productInfoList)-1]
	bar.Clear()
	return &Metadata{
		Product:        latestProductInfo.Name,
		BundleUrl:      latestProductInfo.BundleURL,
		TargetPlatform: latestProductInfo.TargetPlatformVersion,
		DockerImage:    latestProductInfo.LiferayDockerImage,
		GroupId:        strcase.ToDelimited(PackageName, '.'),
		ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
		Name:           strcase.ToCamel(strings.ToLower(base)),
	}, nil
}

func getOfflineMetadata(base, version, edition string) (*Metadata, error) {
	if edition == DXP {
		switch version {
		case "7.4":
			return &Metadata{
				Product:        "7.4.13.u102",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.4.13-u102/liferay-dxp-tomcat-7.4.13.u102-20231109153600206.tar.gz",
				TargetPlatform: "7.4.13.u102",
				DockerImage:    "liferay/dxp:7.4.13-u102",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.3":
			return &Metadata{
				Product:        "dxp-7.3-u35",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.3.10-u35/liferay-dxp-tomcat-7.3.10.u35-20231114110531823.tar.gz",
				TargetPlatform: "7.3.10.u35",
				DockerImage:    "liferay/dxp:7.3.10-u35",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.2":
			return &Metadata{
				Product:        "dxp-7.2-sp8",
				BundleUrl:      "https://api.liferay.com/downloads/portal/7.2.10.8/liferay-dxp-tomcat-7.2.10.8-sp8-slim-20220912234451782.tar.gz",
				TargetPlatform: "7.2.10.8",
				DockerImage:    "liferay/dxp:7.2.10-sp8",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.1":
			return &Metadata{
				Product:        "dxp-7.1-sp8",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.1.10.8/liferay-dxp-tomcat-7.1.10.8-sp8-slim-20220926154152962.tar.gz",
				TargetPlatform: "7.1.10.8",
				DockerImage:    "liferay/dxp:7.1.10-sp8",
				GroupId:        strcase.ToDelimited(PackageName, '.'),
				ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
				Name:           strcase.ToCamel(strings.ToLower(base)),
			}, nil
		case "7.0":
			return &Metadata{
				Product:        "dxp-7.0-sp17",
				BundleUrl:      "https://releases-cdn.liferay.com/dxp/7.0.10.17/liferay-dxp-digital-enterprise-tomcat-7.0.10.17-sp17-slim-20211014075354439.tar.gz",
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
			return &Metadata{
				Product:        "portal-7.4-ga102",
				BundleUrl:      "https://github.com/liferay/liferay-portal/releases/download/7.4.3.102-ga102/liferay-ce-portal-tomcat-7.4.3.102-ga102-20231109165213885.tar.gz",
				TargetPlatform: "7.4.3.102",
				DockerImage:    "liferay/portal:7.4.3.102-ga102",
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
	}
	return nil, fmt.Errorf("invalid or unsupported Liferay version")
}
