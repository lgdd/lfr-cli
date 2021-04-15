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
		printutil.Danger(err.Error())
		os.Exit(1)
	}

	err = procutil.SetCatalinaPid()

	if err != nil {
		printutil.Danger(err.Error())
		os.Exit(1)
	}

	shutdownCmd := exec.Command(shutdownScript)
	shutdownCmd.Stdout = os.Stdout

	err = shutdownCmd.Run()

	if err != nil {
		printutil.Danger(err.Error())
		os.Exit(1)
	}

}
