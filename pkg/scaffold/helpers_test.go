package scaffold

import (
	"encoding/xml"
	"io"
	"os"
	"testing"

	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

// chdirWorkspace changes the working directory to dir and restores it after
// the test. Required for scaffold functions that call GetLiferayWorkspacePath.
func chdirWorkspace(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

// assertPathExists fails the test if the given path does not exist.
func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path to exist: %s (%v)", path, err)
	}
}

// assertPathAbsent fails the test if the given path exists.
func assertPathAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected path to be absent: %s", path)
	}
}

// assertParentPomContainsModule reads the pom.xml at pomPath and fails the
// test if moduleName is not listed as a module.
func assertParentPomContainsModule(t *testing.T, pomPath, moduleName string) {
	t.Helper()
	f, err := os.Open(pomPath)
	if err != nil {
		t.Fatalf("could not open pom.xml at %s: %v", pomPath, err)
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	var pom fileutil.Pom
	if err = xml.Unmarshal(data, &pom); err != nil {
		t.Fatalf("could not parse pom.xml at %s: %v", pomPath, err)
	}
	for _, m := range pom.Modules.Module {
		if m == moduleName {
			return
		}
	}
	t.Fatalf("module %q not found in %s", moduleName, pomPath)
}
