package start

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
		Use:   "start",
		Short: "Start Liferay bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	startupScript, err := fileutil.GetTomcatScriptPath("startup")

	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}

	procutil.SetCatalinaPid()

	startupCmd := exec.Command(startupScript)
	startupCmd.Stdout = os.Stdout

	err = startupCmd.Run()

	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}
}
