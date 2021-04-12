package logs

import (
	"fmt"
	"os"

	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "logs",
		Short: "Display logs from the running Liferay bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
	Follow bool
)

func init() {
	Cmd.Flags().BoolVarP(&Follow, "follow", "f", false, "--follow")
}

func run(cmd *cobra.Command, args []string) {
	logFile, err := fileutil.GetCatalinaLogFile()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.Tail(logFile, Follow)
}
