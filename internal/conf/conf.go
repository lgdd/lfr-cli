// Package conf manages the CLI configuration stored in ~/.lfr/config.toml via
// Viper, and exposes the Viper key constants and default values used throughout
// the application.
package conf

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Client extension sample repository constants.
const (
	// ClientExtensionSamplePrefix is the directory name prefix for official client extension samples.
	ClientExtensionSamplePrefix = "liferay-sample-"
	// ClientExtensionExtraSamplePrefix is the directory name prefix for extra client extension samples.
	ClientExtensionExtraSamplePrefix = "sample-"
	// ClientExtensionSampleProjectName is the GitHub project name for official client extension samples.
	ClientExtensionSampleProjectName = "liferay-client-extensions-samples"
	// ClientExtensionExtraSampleProjectName is the GitHub project name for extra client extension samples.
	ClientExtensionExtraSampleProjectName = "liferay-client-extensions-extra-samples"
)

// Viper config file keys used throughout the application.
const (
	// WorkspaceEdition is the config key for the Liferay edition (portal or dxp).
	WorkspaceEdition = "workspace.edition"
	// WorkspaceVersion is the config key for the target Liferay version.
	WorkspaceVersion = "workspace.version"
	// WorkspaceBuild is the config key for the build tool (gradle or maven).
	WorkspaceBuild = "workspace.build"
	// WorkspaceInit is the config key for whether to initialize the bundle on workspace creation.
	WorkspaceInit = "workspace.init"
	// LogsFollow is the config key for whether to follow log output continuously.
	LogsFollow = "logs.follow"
	// DeployClean is the config key for whether to run a clean before deploying.
	DeployClean = "deploy.clean"
	// ModulePackage is the config key for the default Java package name.
	ModulePackage = "module.package"
	// DockerMultistage is the config key for whether to use a multi-stage Docker build.
	DockerMultistage = "docker.multistage"
	// DockerJDK is the config key for the JDK version used in Docker images.
	DockerJDK = "docker.jdk"
	// OutputNoColor is the config key for disabling colored output.
	OutputNoColor = "output.no_color"
	// OutputAccessible is the config key for enabling accessible (non-animated) output.
	OutputAccessible = "output.accessible"
)

// NoColor disables colors for printed messages when set to true.
var NoColor bool

// Init loads the configuration from the default config path (~/.lfr/config.toml),
// creating the file with defaults if it does not exist.
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

	if _, statErr := os.Stat(configFilePath); statErr == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err.Error())
		}
	}

	setDefaults()
	viper.WriteConfigAs(configFilePath)
}

// GetConfigPath returns the path to the CLI configuration directory (~/.lfr).
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
