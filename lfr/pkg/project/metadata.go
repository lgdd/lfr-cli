package project

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
)

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

const (
	Gradle = "gradle"
	Maven  = "maven"
)

var PackageName string

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

func NewMetadata(base, version string) (*Metadata, error) {
	switch version {
	case "7.4":
		return &Metadata{
			Product:        "portal-7.4-ga2",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.4.1-ga2/liferay-ce-portal-tomcat-7.4.1-ga2-20210609223456272.tar.gz",
			TomcatVersion:  "9.0.43",
			TargetPlatform: "7.4.1-1",
			DockerImage:    "liferay/portal:7.4.1-ga2",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.3":
		return &Metadata{
			Product:        "portal-7.3-ga8",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.3.7-ga8/liferay-ce-portal-tomcat-7.3.7-ga8-20210610183559721.tar.gz",
			TomcatVersion:  "9.0.43",
			TargetPlatform: "7.3.7",
			DockerImage:    "liferay/portal:7.3.6-ga8",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.2":
		return &Metadata{
			Product:        "portal-7.2-ga2",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.tar.gz",
			TomcatVersion:  "9.0.17",
			TargetPlatform: "7.2.1-1",
			DockerImage:    "liferay/portal:7.2.1-ga2",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	case "7.1":
		return &Metadata{
			Product:        "portal-7.1-ga4",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.tar.gz",
			TomcatVersion:  "9.0.17",
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
			TomcatVersion:  "8.0.32",
			TargetPlatform: "7.0.6-2",
			DockerImage:    "liferay/portal:7.0.6-ga7",
			GroupId:        strcase.ToDelimited(PackageName, '.'),
			ArtifactId:     strcase.ToKebab(strings.ToLower(base)),
			Name:           strcase.ToCamel(strings.ToLower(base)),
		}, nil
	}
	return nil, fmt.Errorf("invalid Liferay version")
}
