package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lgdd/liferay-cli/lfr/pkg/generate/workspace"
	"github.com/lgdd/liferay-cli/lfr/pkg/project"
)

func Test_GenerateDocker_ShouldCreateDockerfileAndCompose(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = Generate(liferayWorkspace, false, 8)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(liferayWorkspace, "Dockerfile")); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(liferayWorkspace, "docker-compose.yml")); err != nil {
		t.Fatal(err)
	}
}

func Test_GenerateDocker_WithMultiStageAndGradle_ShouldContainExpectedStages(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = Generate(liferayWorkspace, true, 8)
	if err != nil {
		t.Fatal(err)
	}
	fileBytes, err := os.ReadFile(filepath.Join(liferayWorkspace, "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	fileContent := string(fileBytes)
	firstLine := strings.Split(fileContent, "\n")[0]
	expectedFirstLine := "FROM azul/zulu-openjdk-alpine:8 AS builder"
	if firstLine != expectedFirstLine {
		t.Fatalf("Found: '%s'\nExpected: '%s'",
			firstLine, expectedFirstLine)
	}
}

func Test_GenerateDocker_WithMultiStageAndMaven_ShouldContainExpectedStages(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Maven, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = Generate(liferayWorkspace, true, 11)
	if err != nil {
		t.Fatal(err)
	}
	fileBytes, err := os.ReadFile(filepath.Join(liferayWorkspace, "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	fileContent := string(fileBytes)
	firstLine := strings.Split(fileContent, "\n")[0]
	expectedFirstLine := "FROM azul/zulu-openjdk-alpine:11 AS builder"
	if firstLine != expectedFirstLine {
		t.Fatalf("Found: '%s'\nExpected: '%s'",
			firstLine, expectedFirstLine)
	}
}

func Test_GenerateDocker_WithWrongJavaVersion_ShouldFail(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = Generate(liferayWorkspace, true, 16)
	if err == nil {
		t.Fatal("wrong java version (!= 8 || 11) should fail")
	}
}
