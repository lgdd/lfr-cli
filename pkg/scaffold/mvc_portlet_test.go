package scaffold

import (
	"path/filepath"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/metadata"
)

func Test_CreateModuleMVC_WithGradle_ShouldCreateExpectedFiles(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	chdirWorkspace(t, liferayWorkspace)

	name := "example-mvc"
	if err := CreateModuleMVC(name); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(liferayWorkspace, "modules", name)
	assertPathExists(t, modulePath)
	assertPathExists(t, filepath.Join(modulePath, ".gitignore"))
	assertPathExists(t, filepath.Join(modulePath, "bnd.bnd"))
	assertPathExists(t, filepath.Join(modulePath, "build.gradle"))
	assertPathExists(t, filepath.Join(modulePath, "src", "main", "resources", "META-INF", "resources", "view.jsp"))
	assertPathAbsent(t, filepath.Join(modulePath, "pom.xml"))
}

func Test_CreateModuleMVC_WithMaven_ShouldCreateExpectedFilesAndUpdateParentPom(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Maven, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	chdirWorkspace(t, liferayWorkspace)

	name := "example-mvc"
	if err := CreateModuleMVC(name); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(liferayWorkspace, "modules", name)
	assertPathExists(t, modulePath)
	assertPathExists(t, filepath.Join(modulePath, "pom.xml"))
	assertPathAbsent(t, filepath.Join(modulePath, "build.gradle"))
	assertParentPomContainsModule(t, filepath.Join(liferayWorkspace, "modules", "pom.xml"), name)
}
