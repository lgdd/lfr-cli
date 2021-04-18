package project

import (
	"fmt"

	"github.com/iancoleman/strcase"
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

func NewMetadata(base, version string) (*Metadata, error) {
	switch version {
	case "7.3":
		return &Metadata{
			Product:        "portal-7.3-ga7",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.3.6-ga7/liferay-ce-portal-tomcat-7.3.6-ga7-20210301155526191.tar.gz",
			TomcatVersion:  "9.0.40",
			TargetPlatform: "7.3.6",
			DockerImage:    "liferay/portal:7.3.6-ga7",
			GroupId:        strcase.ToDelimited(base, '.'),
			ArtifactId:     base,
			Name:           strcase.ToCamel(base),
		}, nil
	case "7.2":
		return &Metadata{
			Product:        "portal-7.2-ga2",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.tar.gz",
			TomcatVersion:  "9.0.17",
			TargetPlatform: "7.2.1",
			DockerImage:    "liferay/portal:7.2.1-ga2",
			GroupId:        strcase.ToDelimited(base, '.'),
			ArtifactId:     base,
			Name:           strcase.ToCamel(base),
		}, nil
	case "7.1":
		return &Metadata{
			Product:        "portal-7.1-ga4",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.tar.gz",
			TomcatVersion:  "9.0.17",
			TargetPlatform: "7.1.3",
			DockerImage:    "liferay/portal:7.1.3-ga4",
			GroupId:        strcase.ToDelimited(base, '.'),
			ArtifactId:     base,
			Name:           strcase.ToCamel(base),
		}, nil
	case "7.0":
		return &Metadata{
			Product:        "portal-7.0-ga7",
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.0.6-ga7/liferay-ce-portal-tomcat-7.0-ga7-20180507111753223.zip",
			TomcatVersion:  "8.0.32",
			TargetPlatform: "7.0.6-2",
			DockerImage:    "liferay/portal:7.0.6-ga7",
			GroupId:        strcase.ToDelimited(base, '.'),
			ArtifactId:     base,
			Name:           strcase.ToCamel(base),
		}, nil
	}
	return nil, fmt.Errorf("invalid Liferay version")
}
