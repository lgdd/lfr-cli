package status

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/lgdd/lfr-cli/pkg/util/procutil"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'status' which tells if Liferay is running or not
	Cmd = &cobra.Command{
		Use:   "status",
		Short: "Status (running or stopped) of a Liferay Tomcat bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		logger.PrintInfo("not available for Windows since the Tomcat process run in another command window")
		os.Exit(0)
	}
	isRunning, pid, err := procutil.IsCatalinaRunning()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if isRunning {
		fmt.Print("Liferay is ")
		logger.PrintSuccess("RUNNING")
		fmt.Printf(" (pid=%v)\n", pid)
	} else {
		fmt.Print("Liferay is ")
		logger.PrintError("STOPPED\n")
	}
}
