package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/helper"
)

func Test_HandleClientExtensionsOffline_ShouldIncludeGitFolderForFutureUpdates(t *testing.T) {
	configTestFolderName := "config_test"
	configTestPath := filepath.Join(t.TempDir(), configTestFolderName)
	if err := os.Mkdir(configTestPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	helper.HandleClientExtensionsOffline(configTestPath)
	if _, err := os.Stat(filepath.Join(configTestPath, conf.ClientExtensionSampleProjectName, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(configTestPath); err != nil {
		t.Fatal(err)
	}
}
