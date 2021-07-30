package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lgdd/liferay-cli/lfr/pkg/project"
)

func TestGenerate_Gradle_ShouldHaveExpectedFiles(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := Generate(workspaceDir, "gradle", "7.3")

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

	metadata, err := project.NewMetadata(workspaceDir, "7.3")

	if err != nil {
		t.Fatal(err)
	}

	hasExpectedValue := strings.Contains(string(gradleProps), metadata.Product)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.Product)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.DockerImage)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.DockerImage)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.BundleUrl)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.BundleUrl)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.TomcatVersion)

	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.TomcatVersion)
	}

	hasExpectedValue = strings.Contains(string(gradleProps), metadata.TargetPlatform)
	if !hasExpectedValue {
		t.Fatalf("gradle.properties doesn't contain %v", metadata.TargetPlatform)
	}
}

func TestGenerate_Maven_ShouldHaveExpectedFiles(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := Generate(workspaceDir, "maven", "7.3")

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

	metadata, err := project.NewMetadata(workspaceDir, "7.3")

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

func TestGenerateWorkspace_WithWrongVersion_ShouldFail(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := Generate(workspaceDir, "gradle", "6.2")

	if err == nil {
		t.Fatal("workspace with wrong version should fail")
	}

	_, err = os.Stat(workspaceDir)

	if err == nil {
		t.Fatal("workspace should not be created")
	}
}

func TestGenerateWorkspace_WithWrongBuildTool_ShouldFail(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	err := Generate(workspaceDir, "ant", "7.3")

	if err == nil {
		t.Fatal("workspace with wrong build tool should fail")
	}

	_, err = os.Stat(workspaceDir)

	if err == nil {
		t.Fatal("workspace should not be created")
	}
}

func TestGenerateWorkspace_WithMaven_ShouldHaveDefaultPackageName(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	project.PackageName = "org.acme"
	err := Generate(workspaceDir, "maven", "7.3")
	if err != nil {
		t.Fatal(err)
	}
	pomXml, err := os.ReadFile(filepath.Join(workspaceDir, "pom.xml"))
	if !strings.Contains(string(pomXml), "org.acme") {
		t.Fatal("workspace pom.xml doesn't contain default packag 'org.acme'")
	}
}

func TestGenerateWorkspace_WithGradle_ShouldHaveDefaultPackageName(t *testing.T) {
	workspaceDir := filepath.Join(t.TempDir(), "liferay-workspace")
	project.PackageName = "org.acme"
	err := Generate(workspaceDir, "gradle", "7.3")
	if err != nil {
		t.Fatal(err)
	}
	buildGradle, err := os.ReadFile(filepath.Join(workspaceDir, "build.gradle"))
	if !strings.Contains(string(buildGradle), "org.acme") {
		t.Fatal("workspace build.gradle doesn't contain default packag 'org.acme'")
	}
}
