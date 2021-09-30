package project

import "encoding/xml"

// XMLHeader is the first line to be written in an XML file
const (
	XMLHeader = `<?xml version="1.0"?>` + "\n"
)

// Pom represents the common structure of a Maven pom file
type Pom struct {
	XMLName        xml.Name `xml:"project"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xmlns:schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	Parent         struct {
		GroupId      string `xml:"groupId"`
		ArtifactId   string `xml:"artifactId"`
		Version      string `xml:"version"`
		RelativePath string `xml:"relativePath"`
	} `xml:"parent"`
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Name       string `xml:"name"`
	Packaging  string `xml:"packaging"`
	Modules    struct {
		Module []string `xml:"module"`
	} `xml:"modules"`
}

// WorkspacePom represents the common structure of a parent Maven pom file in a Liferay Workspace
type WorkspacePom struct {
	XMLName        xml.Name `xml:"project"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xmlns:schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	Parent         struct {
		GroupId      string `xml:"groupId"`
		ArtifactId   string `xml:"artifactId"`
		Version      string `xml:"version"`
		RelativePath string `xml:"relativePath"`
	} `xml:"parent"`
	ArtifactId string `xml:"artifactId"`
	Name       string `xml:"name"`
	Packaging  string `xml:"packaging"`
	Modules    struct {
		Module []string `xml:"module"`
	} `xml:"modules"`
	Properties struct {
		LiferayBomVersion          string `xml:"liferay.bom.version"`
		LiferayDockerImage         string `xml:"liferay.docker.image"`
		LiferayWorkspaceBundleURL  string `xml:"liferay.workspace.bundle.url"`
		LiferayRepositoryURL       string `xml:"liferay.repository.url"`
		ProjectBuildSourceEncoding string `xml:"project.build.sourceEncoding"`
	} `xml:"properties"`
}
