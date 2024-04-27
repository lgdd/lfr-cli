package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/metadata"
)

func Test_CreateDockerFiles_ShouldCreateDockerfileAndCompose(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = CreateDockerFiles(liferayWorkspace, false, 8)
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

func Test_CreateDockerFiles_WithMultiStageAndGradle_ShouldContainExpectedStages(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = CreateDockerFiles(liferayWorkspace, true, 8)
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

func Test_CreateDockerFiles_WithMultiStageAndMaven_ShouldContainExpectedStages(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(liferayWorkspace, metadata.Maven, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = CreateDockerFiles(liferayWorkspace, true, 11)
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

func Test_CreateDockerFiles_WithWrongJavaVersion_ShouldFail(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal")
	if err != nil {
		t.Fatal(err)
	}
	err = CreateDockerFiles(liferayWorkspace, true, 16)
	if err == nil {
		t.Fatal("wrong java version (!= 8 || 11) should fail")
	}
}
