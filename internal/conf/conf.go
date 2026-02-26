package conf

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	ClientExtensionSamplePrefix           = "liferay-sample-"
	ClientExtensionExtraSamplePrefix      = "sample-"
	ClientExtensionSampleProjectName      = "liferay-client-extensions-samples"
	ClientExtensionExtraSampleProjectName = "liferay-client-extensions-extra-samples"
	// Config file keys
	WorkspaceEdition = "workspace.edition"
	WorkspaceVersion = "workspace.version"
	WorkspaceBuild   = "workspace.build"
	WorkspaceInit    = "workspace.init"
	LogsFollow       = "logs.follow"
	DeployClean      = "deploy.clean"
	ModulePackage    = "module.package"
	DockerMultistage = "docker.multistage"
	DockerJDK        = "docker.jdk"
	OutputNoColor    = "output.no_color"
	OutputAccessible = "output.accessible"
)

// NoColor allows to disable colors for printed messages, default is false
var NoColor bool

func Init() {
	initWithPath(GetConfigPath())
}

func initWithPath(configPath string) {
	configFile := "config.toml"
	configName := strings.Split(configFile, ".")[0]
	configType := strings.Split(configFile, ".")[1]

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	createConfigFolderAt(configPath)

	configFilePath := filepath.Join(configPath, configFile)
	configFileInfo, _ := os.Stat(configFilePath)

	// if empty, remove to set defaults
	if configFileInfo != nil && configFileInfo.Size() == 0 {
		os.Remove(configFilePath)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatal(err.Error())
		}
	}

	setDefaults()
	viper.WriteConfigAs(configFilePath)
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err.Error())
	}

	return filepath.Join(homeDir, ".lfr")
}

func setDefaults() {
	setDefault(WorkspaceEdition, "portal")
	setDefault(WorkspaceVersion, "7.4")
	setDefault(WorkspaceBuild, "gradle")
	setDefault(WorkspaceInit, false)
	setDefault(LogsFollow, false)
	setDefault(DeployClean, false)
	setDefault(ModulePackage, "org.acme")
	setDefault(DockerMultistage, false)
	setDefault(DockerJDK, 11)
	setDefault(OutputNoColor, false)
	setDefault(OutputAccessible, false)
}

func setDefault(key string, value any) {
	if viper.Get(key) == nil {
		viper.SetDefault(key, value)
	}
}

func createConfigFolderAt(configPath string) {
	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
