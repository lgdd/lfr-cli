// Package config implements the config subcommand, which lets users read and
// write entries in the ~/.lfr/config.toml configuration file.
package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Cmd is the command 'config' which gets and sets lfr configuration values.
	Cmd = &cobra.Command{
		Use:   "config",
		Short: "Get and set your configuration for lfr",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
				return err
			}
			if len(args) == 1 && strings.Contains(args[0], "=") {
				key := strings.Split(args[0], "=")[0]
				if len(key) > 0 && !slices.Contains(viper.AllKeys(), key) {
					return fmt.Errorf("invalid config key")
				}
			}
			if len(args) > 1 && !slices.Contains(viper.AllKeys(), args[0]) {
				return fmt.Errorf("invalid config key")
			}
			return nil
		},
		Run: run,
	}
	// List enables printing all configuration key/value pairs when set to true.
	List bool
)

func init() {
	conf.Init()
	Cmd.Flags().BoolVarP(&List, "list", "l", false, "list the key/values for your current configuration")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		printConfKeyValues()
		cmd.Help()
		os.Exit(0)
	}

	if len(args) == 1 {
		if strings.Contains(args[0], "=") {
			setKeyValue(args[0])
		}
		logger.Println(viper.GetString(args[0]))
		os.Exit(0)
	}

	if len(args) == 2 {
		viper.Set(args[0], args[1])
		viper.WriteConfig()
	}
}

func printConfKeyValues() {
	if List {
		for _, key := range viper.AllKeys() {
			var keyValueBuilder strings.Builder
			keyValueBuilder.WriteString(key)
			keyValueBuilder.WriteString("=")
			keyValueBuilder.WriteString(viper.GetString(key))
			logger.Println(keyValueBuilder.String())
		}
		os.Exit(0)
	}
}

func setKeyValue(arg string) {
	keyValue := strings.Split(arg, "=")
	key := keyValue[0]
	value := keyValue[1]
	if len(key) > 0 && len(value) > 0 {
		viper.Set(key, value)
		viper.WriteConfig()
	}
}
