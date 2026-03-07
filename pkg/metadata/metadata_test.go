package metadata

import (
	"strings"
	"testing"
)

func Test_NewWorkspaceData_WithDXP73_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.3.10.u36"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.3", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXP72_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.2.10.8"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.2", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXP71_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.1.10.fp28"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.1", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXP70_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.0.10.17"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.0", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortal73_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.3.7"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortal72_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.2.1-1"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.2", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortal71_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.1.3-1"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.1", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortal70_ShouldReturnLatestValid(t *testing.T) {
	expectedTargetPlatform := "7.0.6-2"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.0", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXPQuarterly2024Q1_ShouldReturnLatestValid(t *testing.T) {
	metadata, err := NewWorkspaceData("liferay-workspace", "2024.q1", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(metadata.TargetPlatform, "2024.q1") {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform to start with: '2024.q1'",
			metadata.TargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXPQuarterly2025Q4_ShouldReturnLatestValid(t *testing.T) {
	metadata, err := NewWorkspaceData("liferay-workspace", "2025.q4", "dxp")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(metadata.TargetPlatform, "2025.q4") {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform to start with: '2025.q4'",
			metadata.TargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortalGA112_ShouldReturnValid(t *testing.T) {
	expectedTargetPlatform := "7.4.3.112"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.4.3.112-ga112", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithPortalGA132_ShouldReturnValid(t *testing.T) {
	expectedTargetPlatform := "7.4.3.132"
	metadata, err := NewWorkspaceData("liferay-workspace", "7.4.3.132-ga132", "portal")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.TargetPlatform != expectedTargetPlatform {
		t.Fatalf("Found Target Platform: '%s'\nExpected Target Platform: '%s'",
			metadata.TargetPlatform, expectedTargetPlatform)
	}
}

func Test_NewWorkspaceData_WithDXPUnsupportedQuarter_ShouldFail(t *testing.T) {
	_, err := NewWorkspaceData("liferay-workspace", "2022.q1", "dxp")
	if err == nil {
		t.Fatal("metadata with unsupported DXP quarterly version should fail")
	}
}

func Test_NewWorkspaceData_WithPortalUnknownGA_ShouldFail(t *testing.T) {
	_, err := NewWorkspaceData("liferay-workspace", "7.4.3.999-ga999", "portal")
	if err == nil {
		t.Fatal("metadata with unknown portal GA version should fail")
	}
}

func Test_NewWorkspaceData_WithWrongVersion_ShouldFail(t *testing.T) {
	_, err := NewWorkspaceData("liferay-workspace", "6.2", "dxp")
	if err == nil {
		t.Fatal("metadata with wrong Liferay major version should fail")
	}
}

func Test_NewWorkspaceData_WithWrongEdition_ShouldFail(t *testing.T) {
	_, err := NewWorkspaceData("liferay-workspace", "7.4", "opensource")
	if err == nil {
		t.Fatal("metadata with wrong Liferay major version should fail")
	}
}

func Test_NewWorkspaceData_WithGivenName_ShouldReturnFormattedData(t *testing.T) {
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

func getMetadataArrayForAllVersion(name string) ([]*WorkspaceData, error) {
	var metadataArray []*WorkspaceData
	metadata73, err := NewWorkspaceData(name, "7.3", "portal")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata73)
	metadata72, err := NewWorkspaceData(name, "7.2", "portal")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata72)
	metadata71, err := NewWorkspaceData(name, "7.1", "portal")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata71)
	metadata70, err := NewWorkspaceData(name, "7.0", "portal")
	if err != nil {
		return nil, err
	}
	metadataArray = append(metadataArray, metadata70)
	return metadataArray, nil
}
