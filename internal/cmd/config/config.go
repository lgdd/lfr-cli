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
	Cmd = &cobra.Command{
		Use:   "config",
		Short: "Get and set your configuration for lfr",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
				return err
			}
			if len(args) > 0 && !slices.Contains(viper.AllKeys(), args[0]) {
				return fmt.Errorf("invalid config key")
			}
			return nil
		},
		Run: run,
	}
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
