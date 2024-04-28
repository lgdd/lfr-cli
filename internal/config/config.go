package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config keys
const (
	WorkspaceEdition = "workspace.edition"
	WorkspaceVersion = "workspace.version"
	WorkspaceBuild   = "workspace.build"
	WorkspaceInit    = "workspace.init"
	LogsFollow       = "logs.follow"
	ModulePackage    = "module.package"
	DockerMultistage = "docker.multistage"
	DockerJDK        = "docker.jdk"
	OutputNoColor    = "output.no_color"
)

// NoColor allows to disable colors for printed messages, default is false
var NoColor bool

func Init() {
	configFile := "config.toml"
	configName := strings.Split(configFile, ".")[0]
	configType := strings.Split(configFile, ".")[1]

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(GetConfigPath())

	configFilePath := filepath.Join(GetConfigPath(), configFile)
	configFileInfo, _ := os.Stat(configFilePath)

	// if empty, remove to set defaults
	if configFileInfo != nil && configFileInfo.Size() == 0 {
		os.Remove(configFilePath)
	}

	setDefaults()
	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err.Error())
	}

	return filepath.Join(homeDir, ".lfr")
}

func setDefaults() {
	viper.SetDefault(WorkspaceEdition, "portal")
	viper.SetDefault(WorkspaceVersion, "7.4")
	viper.SetDefault(WorkspaceBuild, "gradle")
	viper.SetDefault(WorkspaceInit, false)
	viper.SetDefault(LogsFollow, false)
	viper.SetDefault(ModulePackage, "org.acme")
	viper.SetDefault(DockerMultistage, false)
	viper.SetDefault(DockerJDK, 11)
	viper.SetDefault(OutputNoColor, false)
}
