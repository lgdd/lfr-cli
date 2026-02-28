// Package logs implements the logs subcommand, which tails the Liferay Tomcat
// catalina.out log file with an optional follow mode.
package logs

import (
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Cmd is the command 'logs' to display Liferay bundle's logs
	Cmd = &cobra.Command{
		Use:     "logs",
		Aliases: []string{"l"},
		Short:   "Display logs from the running Liferay bundle",
		Args:    cobra.NoArgs,
		Run:     run,
	}
	// Follow continuously streams new log lines when set to true.
	Follow bool
)

func init() {
	conf.Init()
	defaultFollow := viper.GetBool(conf.LogsFollow)
	Cmd.Flags().BoolVarP(&Follow, "follow", "f", defaultFollow, "--follow")
}

func run(cmd *cobra.Command, args []string) {
	logFile, err := fileutil.GetCatalinaLogFile()

	if err != nil {
		logger.Fatal(err.Error())
	}

	fileutil.Tail(logFile, Follow)
}
