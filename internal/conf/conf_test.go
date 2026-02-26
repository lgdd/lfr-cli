package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitWithPath_FirstInstall_ShouldCreateConfigWithDefaults(t *testing.T) {
	configPath := t.TempDir()

	t.Cleanup(func() { viper.Reset() })

	initWithPath(configPath)

	configFilePath := filepath.Join(configPath, "config.toml")
	if _, err := os.Stat(configFilePath); err != nil {
		t.Fatalf("config file was not created at %s: %v", configFilePath, err)
	}

	defaults := map[string]any{
		WorkspaceEdition: "portal",
		WorkspaceVersion: "7.4",
		WorkspaceBuild:   "gradle",
		WorkspaceInit:    false,
		LogsFollow:       false,
		DeployClean:      false,
		ModulePackage:    "org.acme",
		DockerMultistage: false,
		DockerJDK:        11,
		OutputNoColor:    false,
		OutputAccessible: false,
	}

	for key, expected := range defaults {
		got := viper.Get(key)
		if got != expected {
			t.Errorf("key %q: expected %v (%T), got %v (%T)", key, expected, expected, got, got)
		}
	}
}
