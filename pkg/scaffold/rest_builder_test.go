package scaffold

import (
	"path/filepath"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/metadata"
)

func Test_CreateModuleRESTBuilder_WithGradle_ShouldCreateExpectedDirs(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	name := "example-rb"
	if err := CreateModuleRESTBuilder(liferayWorkspace, name); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(liferayWorkspace, "modules", name)
	apiPath := filepath.Join(modulePath, name+"-api")
	implPath := filepath.Join(modulePath, name+"-impl")

	assertPathExists(t, modulePath)
	assertPathExists(t, apiPath)
	assertPathExists(t, implPath)
	assertPathExists(t, filepath.Join(implPath, "rest-openapi.yaml"))
	assertPathAbsent(t, filepath.Join(modulePath, "pom.xml"))
	assertPathAbsent(t, filepath.Join(apiPath, "pom.xml"))
	assertPathAbsent(t, filepath.Join(implPath, "pom.xml"))
}

func Test_CreateModuleRESTBuilder_WithMaven_ShouldCreateExpectedDirsAndUpdateParentPom(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Maven, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	name := "example-rb"
	if err := CreateModuleRESTBuilder(liferayWorkspace, name); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(liferayWorkspace, "modules", name)
	apiPath := filepath.Join(modulePath, name+"-api")
	implPath := filepath.Join(modulePath, name+"-impl")

	assertPathExists(t, modulePath)
	assertPathExists(t, apiPath)
	assertPathExists(t, implPath)
	assertPathExists(t, filepath.Join(modulePath, "pom.xml"))
	assertPathExists(t, filepath.Join(apiPath, "pom.xml"))
	assertPathExists(t, filepath.Join(implPath, "pom.xml"))
	assertPathAbsent(t, filepath.Join(apiPath, "build.gradle"))
	assertPathAbsent(t, filepath.Join(implPath, "build.gradle"))
	assertParentPomContainsModule(t, filepath.Join(liferayWorkspace, "modules", "pom.xml"), name)
}
