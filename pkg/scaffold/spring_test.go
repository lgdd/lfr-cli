package scaffold

import (
	"path/filepath"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/metadata"
)

func Test_CreateModuleSpring_WithThymeleaf_ShouldCreateExpectedFiles(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	chdirWorkspace(t, liferayWorkspace)

	name := "example-spring"
	if err := CreateModuleSpring(name, "thymeleaf"); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(liferayWorkspace, "modules", name)
	viewsPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "views")

	assertPathExists(t, modulePath)
	assertPathExists(t, filepath.Join(modulePath, "build.gradle"))
	assertPathExists(t, filepath.Join(viewsPath, "user.html"))
	assertPathExists(t, filepath.Join(viewsPath, "greeting.html"))
	assertPathAbsent(t, filepath.Join(viewsPath, "user.jspx"))
	assertPathAbsent(t, filepath.Join(viewsPath, "greeting.jspx"))
}

func Test_CreateModuleSpring_WithJSP_ShouldCreateJspxViews(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	chdirWorkspace(t, liferayWorkspace)

	name := "example-spring"
	if err := CreateModuleSpring(name, "jsp"); err != nil {
		t.Fatal(err)
	}

	viewsPath := filepath.Join(liferayWorkspace, "modules", name, "src", "main", "webapp", "WEB-INF", "views")

	assertPathExists(t, filepath.Join(viewsPath, "user.jspx"))
	assertPathExists(t, filepath.Join(viewsPath, "greeting.jspx"))
	assertPathAbsent(t, filepath.Join(viewsPath, "user.html"))
	assertPathAbsent(t, filepath.Join(viewsPath, "greeting.html"))
}

func Test_CreateModuleSpring_WithInvalidEngine_ShouldReturnError(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	if err := CreateWorkspace(liferayWorkspace, metadata.Gradle, "7.3", "portal"); err != nil {
		t.Fatal(err)
	}

	chdirWorkspace(t, liferayWorkspace)

	err := CreateModuleSpring("example-spring", "freemarker")
	if err == nil {
		t.Fatal("expected an error for unsupported template engine")
	}
}
