package cx

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_HandleClientExtensionsOffline_ShouldIncludeGitFolderForFutureUpdates(t *testing.T) {
	configTestFolderName := "config_test"
	configTestPath := filepath.Join(t.TempDir(), configTestFolderName)
	if err := os.Mkdir(configTestPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	HandleClientExtensionsOffline(configTestPath)
	if _, err := os.Stat(filepath.Join(configTestPath, "liferay-portal", ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(configTestPath); err != nil {
		t.Fatal(err)
	}
}
