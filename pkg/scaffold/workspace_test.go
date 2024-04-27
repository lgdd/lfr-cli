package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/metadata"
)

func TestCreateWorkspace_Gradle_ShouldHaveExpectedFiles(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(workspaceDir, "gradle", "7.3", "portal")

	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		"configs",
		"gradle",
		"modules",
		"themes",
		"wars",
		".gitignore",
		"gradle.properties",
		"gradle-local.properties",
		"gradlew",
		"gradlew.bat",
		"platform.bndrun",
		"settings.gradle",
	}

	for _, file := range expectedFiles {
		_, err = os.Stat(filepath.Join(workspaceDir, file))

		if err != nil {
			t.Fatal(err)
		}
	}

	gradleProps, err := os.ReadFile(filepath.Join(workspaceDir, "gradle.properties"))

	if err != nil {
		t.Fatal(err)
	}

	metadata, err := metadata.NewWorkspaceData(workspaceDir, "7.3", "portal")

	if err != nil {
		t.Fatal(err)
	}

	hasExpectedValue := strings.Contains(string(gradleProps), metadata.Product)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.Product)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.GithubBundleUrl)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.GithubBundleUrl)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.DockerImage)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.DockerImage)
	}

}

func TestCreateWorkspace_Maven_ShouldHaveExpectedFiles(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(workspaceDir, "maven", "7.3", "portal")

	defer t.Cleanup(func() {
		err = os.RemoveAll(workspaceDir)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		"configs",
		".mvn",
		"modules",
		"themes",
		"wars",
		".gitignore",
		"mvnw",
		"mvnw.cmd",
		"platform.bndrun",
		"pom.xml",
	}

	for _, file := range expectedFiles {
		_, err = os.Stat(filepath.Join(workspaceDir, file))

		if err != nil {
			t.Fatal(err)
		}
	}

	pomXml, err := os.ReadFile(filepath.Join(workspaceDir, "pom.xml"))

	if err != nil {
		t.Fatal(err)
	}

	metadata, err := metadata.NewWorkspaceData(workspaceDir, "7.3", "portal")

	if err != nil {
		t.Fatal(err)
	}

	hasExpectedValue := strings.Contains(string(pomXml), metadata.TargetPlatform)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.TargetPlatform)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.DockerImage)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.DockerImage)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.BundleUrl)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.BundleUrl)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.GithubBundleUrl)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.GithubBundleUrl)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.GroupId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.GroupId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.ArtifactId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.ArtifactId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.Name)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.Name)
	}

	pomXml, err = os.ReadFile(filepath.Join(workspaceDir, "modules", "pom.xml"))

	hasExpectedValue = strings.Contains(string(pomXml), metadata.GroupId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.GroupId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.ArtifactId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.ArtifactId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.Name)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.Name)
	}

	pomXml, err = os.ReadFile(filepath.Join(workspaceDir, "themes", "pom.xml"))

	hasExpectedValue = strings.Contains(string(pomXml), metadata.GroupId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.GroupId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.ArtifactId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.ArtifactId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.Name)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.Name)
	}

	pomXml, err = os.ReadFile(filepath.Join(workspaceDir, "wars", "pom.xml"))

	hasExpectedValue = strings.Contains(string(pomXml), metadata.GroupId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.GroupId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.ArtifactId)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.ArtifactId)
	}

	hasExpectedValue = strings.Contains(string(pomXml), metadata.Name)

	if !hasExpectedValue {
		t.Fatalf("pom.xml doesn't contain %v", metadata.Name)
	}
}

func TestCreateWorkspaceWorkspace_WithWrongVersion_ShouldFail(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(workspaceDir, "gradle", "6.2", "portal")

	if err == nil {
		t.Fatal("workspace with wrong version should fail")
	}

	_, err = os.Stat(workspaceDir)

	if err == nil {
		t.Fatal("workspace should not be created")
	}
}

func TestCreateWorkspaceWorkspace_WithWrongBuildTool_ShouldFail(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(workspaceDir, "ant", "7.3", "portal")

	if err == nil {
		t.Fatal("workspace with wrong build tool should fail")
	}

	_, err = os.Stat(workspaceDir)

	if err == nil {
		t.Fatal("workspace should not be created")
	}
}

func TestCreateWorkspaceWorkspace_WithMaven_ShouldHaveDefaultPackageName(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	metadata.PackageName = "org.acme"
	err := CreateWorkspace(workspaceDir, "maven", "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	pomXml, err := os.ReadFile(filepath.Join(workspaceDir, "pom.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pomXml), "org.acme") {
		t.Fatal("workspace pom.xml doesn't contain default packag 'org.acme'")
	}
}

func TestCreateWorkspaceWorkspace_WithGradle_ShouldHaveDefaultPackageName(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	metadata.PackageName = "org.acme"
	err := CreateWorkspace(workspaceDir, "gradle", "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	buildGradle, err := os.ReadFile(filepath.Join(workspaceDir, "build.gradle"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(buildGradle), "org.acme") {
		t.Fatal("workspace build.gradle doesn't contain default packag 'org.acme'")
	}
}
