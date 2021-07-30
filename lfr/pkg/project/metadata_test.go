package project

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func Test_NewMetadata_With73_ShouldReturnLatestValid(t *testing.T) {
	liferayProductsInfo := getProductsInfo(t)
	expectedTargetPlatform := liferayProductsInfo.Portal73Ga8.TargetPlatformVersion
	expectedTomcatVersion := liferayProductsInfo.Portal73Ga8.TomcatVersion
	metadata, err := NewMetadata("liferay-workspace", "7.3")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
	if metadata.TomcatVersion != expectedTomcatVersion {
		t.Fatalf("Found Tomcat Version: '%s'\nExpected Tomcat Version: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewMetadata_With72_ShouldReturnLatestValid(t *testing.T) {
	liferayProductsInfo := getProductsInfo(t)
	expectedTargetPlatform := liferayProductsInfo.Portal72Ga2.TargetPlatformVersion
	expectedTomcatVersion := liferayProductsInfo.Portal72Ga2.TomcatVersion
	metadata, err := NewMetadata("liferay-workspace", "7.2")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
	if metadata.TomcatVersion != expectedTomcatVersion {
		t.Fatalf("Found Tomcat Version: '%s'\nExpected Tomcat Version: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewMetadata_With71_ShouldReturnLatestValid(t *testing.T) {
	liferayProductsInfo := getProductsInfo(t)
	expectedTargetPlatform := liferayProductsInfo.Portal71Ga4.TargetPlatformVersion
	expectedTomcatVersion := liferayProductsInfo.Portal71Ga4.TomcatVersion
	metadata, err := NewMetadata("liferay-workspace", "7.1")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
	if metadata.TomcatVersion != expectedTomcatVersion {
		t.Fatalf("Found Tomcat Version: '%s'\nExpected Tomcat Version: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewMetadata_With70_ShouldReturnLatestValid(t *testing.T) {
	liferayProductsInfo := getProductsInfo(t)
	expectedTargetPlatform := liferayProductsInfo.Portal70Ga7.TargetPlatformVersion
	expectedTomcatVersion := liferayProductsInfo.Portal70Ga7.TomcatVersion
	metadata, err := NewMetadata("liferay-workspace", "7.0")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
	if metadata.TomcatVersion != expectedTomcatVersion {
		t.Fatalf("Found Tomcat Version: '%s'\nExpected Tomcat Version: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewMetadata_WithWrongVersion_ShouldFail(t *testing.T) {
	_, err := NewMetadata("liferay-workspace", "6.2")
	if err == nil {
		t.Fatal("metadata with wrong Liferay major version should fail")
	}
}

func Test_NewMetadata_WithGivenName_ShouldReturnFormattedData(t *testing.T) {
	name := "ExAmPle-WorkSpace"
	expectedName := "ExampleWorkspace"
	expectedArtifactID := "example-workspace"
	metadataArray, err := getMetadataArrayForAllVersion(name)
	if err != nil {
		t.Fatal(err)
	}
	for _, metadata := range metadataArray {
		if metadata.Name != expectedName {
			t.Fatalf("Metadata %v\nFound Name: '%s'\nExpected Name: '%s'",
				metadata.Product, metadata.Name, expectedName)
		}
		if metadata.ArtifactId != expectedArtifactID {
			t.Fatalf("Metadata %v\nFound ArtifactId: '%s'\nExpected ArtifactId: '%s'",
				metadata.Product, metadata.ArtifactId, expectedArtifactID)
		}
	}
}

func getMetadataArrayForAllVersion(name string) ([]*Metadata, error) {
	var metadataArray []*Metadata
	metadata73, err := NewMetadata(name, "7.3")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata73)
	metadata72, err := NewMetadata(name, "7.2")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata72)
	metadata71, err := NewMetadata(name, "7.1")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata71)
	metadata70, err := NewMetadata(name, "7.0")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata70)
	return metadataArray, nil
}

func getProductsInfo(t *testing.T) LiferayProductsInfo {
	resp, err := http.Get("https://releases-cdn.liferay.com/tools/workspace/.product_info.json")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var liferayProductsInfo LiferayProductsInfo
	err = json.Unmarshal(body, &liferayProductsInfo)
	if err != nil {
		t.Fatal(err)
	}
	return liferayProductsInfo
}

type LiferayProductsInfo struct {
	Dxp70Sp12 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp12"`
	Commerce102 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"commerce-1.0.2"`
	Dxp72Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.2-ga1"`
	Dxp71Sp4 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-sp4"`
	Dxp72Sp1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.2-sp1"`
	Dxp71Sp3 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-sp3"`
	Dxp71Sp2 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-sp2"`
	Dxp71Sp1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-sp1"`
	Dxp70Sp13 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp13"`
	Dxp70Sp11 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp11"`
	Dxp70Sp10 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp10"`
	Commerce206 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"commerce-2.0.6"`
	Dxp71Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-ga1"`
	Dxp70Sp8 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp8"`
	Commerce20772 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"commerce-2.0.7-7.2"`
	Commerce20771 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"commerce-2.0.7-7.1"`
	Portal73Ga3 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga3"`
	Portal73Ga2 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga2"`
	Portal73Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga1"`
	Portal72Ga2 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.2-ga2"`
	Portal72Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.2-ga1"`
	Portal71Ga4 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.1-ga4"`
	Portal71Ga3 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.1-ga3"`
	Portal71Ga2 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.1-ga2"`
	Portal71Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.1-ga1"`
	Portal70Ga7 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.0-ga7"`
	Dxp72Sp2 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.2-sp2"`
	Portal73Ga4 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga4"`
	Dxp73Ep4 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.3-ep4"`
	Dxp73Ep3 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.3-ep3"`
	Dxp70Sp14 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp14"`
	Portal73Ga5 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga5"`
	Dxp73Ep5 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.3-ep5"`
	Dxp72Fp7 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.2-fp7"`
	Dxp71Fp19 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-fp19"`
	Dxp72Sp3 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-sp3"`
	Dxp73Ga1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.3-ga1"`
	Portal73Ga6 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"portal-7.3-ga6"`
	Dxp70Sp15 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.0-sp15"`
	Dxp72Fp9 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.2-fp9"`
	Dxp71Sp5 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
	} `json:"dxp-7.1-sp5"`
	Dxp72Fp10 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-fp10"`
	Dxp71Fp21 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.1-fp21"`
	Dxp70De97 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.0-de97"`
	Dxp72Fp11 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-fp11"`
	Dxp72Sp4 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-sp4"`
	Dxp73Sp1 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.3-sp1"`
	Portal73Ga8 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"portal-7.3-ga8"`
	Dxp71Fp22 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.1-fp22"`
	Dxp70De98 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.0-de98"`
	Dxp72Fp5 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-fp5"`
	Dxp72Fp12 struct {
		TomcatVersion         string `json:"appServerTomcatVersion"`
		BundleURL             string `json:"bundleUrl"`
		BundleChecksumMD5     string `json:"bundleChecksumMD5"`
		BundleChecksumMD5URL  string `json:"bundleChecksumMD5Url"`
		DockerImage           string `json:"liferayDockerImage"`
		ProductVersion        string `json:"liferayProductVersion"`
		Promoted              string `json:"promoted"`
		ReleaseDate           string `json:"releaseDate"`
		TargetPlatformVersion string `json:"targetPlatformVersion"`
	} `json:"dxp-7.2-fp12"`
}
