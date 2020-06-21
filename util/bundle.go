package util

import (
	"fmt"
)

type Project struct {
	BundleUrl      string
	TomcatVersion  string
	TargetPlatform string
	DockerImage    string
	GroupId        string
	ArtifactId     string
	Name           string
}

func GetProject(version string) (*Project, error) {
	switch version {
	case "7.3":
		return &Project{
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.3.2-ga3/liferay-ce-portal-tomcat-7.3.2-ga3-20200519164024819.tar.gz",
			TomcatVersion:  "9.0.33",
			TargetPlatform: "7.3.2",
			DockerImage:    "liferay/portal:7.3.2-ga3",
		}, nil
	case "7.2":
		return &Project{
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.2.1-ga2/liferay-ce-portal-tomcat-7.2.1-ga2-20191111141448326.tar.gz",
			TomcatVersion:  "9.0.17",
			TargetPlatform: "7.2.1",
			DockerImage:    "liferay/portal:7.2.1-ga2",
		}, nil
	case "7.1":
		return &Project{
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.1.3-ga4/liferay-ce-portal-tomcat-7.1.3-ga4-20190508171117552.tar.gz",
			TomcatVersion:  "9.0.17",
			TargetPlatform: "7.1.3",
			DockerImage:    "liferay/portal:7.1.3-ga4",
		}, nil
	case "7.0":
		return &Project{
			BundleUrl:      "https://releases-cdn.liferay.com/portal/7.0.6-ga7/liferay-ce-portal-tomcat-7.0-ga7-20180507111753223.zip",
			TomcatVersion:  "8.0.32",
			TargetPlatform: "7.0.6",
			DockerImage:    "liferay/portal:7.0.6-ga7",
		}, nil
	}
	return nil, fmt.Errorf("invalid Liferay version")
}
