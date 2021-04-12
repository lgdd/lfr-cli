package stop

import (
	"os"
	"os/exec"

	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/lgdd/deba/pkg/util/procutil"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop Liferay bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {

	shutdownScript, err := fileutil.GetTomcatScriptPath("shutdown")

	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}

	procutil.SetCatalinaPid()

	shutdownCmd := exec.Command(shutdownScript)
	shutdownCmd.Stdout = os.Stdout

	err = shutdownCmd.Run()

	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}
}
