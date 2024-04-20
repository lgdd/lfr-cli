package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/spf13/viper"
)

func Init() {
	configPath := GetConfigPath()
	configName := "config"
	configType := "toml"
	configFullPath := filepath.Join(configPath, configName+"."+configType)

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	if _, err := os.Stat(configFullPath); err != nil {
		fileutil.CreateDirs(configPath)
		fileutil.CreateFiles([]string{configFullPath})
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Read config error: %s", err))
	}
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".lfr")
}
