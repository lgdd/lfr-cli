package logs

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
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
	Follow bool
)

func init() {
	config.Init()
	defaultFollow := viper.GetBool(config.LogsFollow)
	Cmd.Flags().BoolVarP(&Follow, "follow", "f", defaultFollow, "--follow")
}

func run(cmd *cobra.Command, args []string) {
	logFile, err := fileutil.GetCatalinaLogFile()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.Tail(logFile, Follow)
}
