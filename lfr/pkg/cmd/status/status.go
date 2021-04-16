package status

import (
	"fmt"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/procutil"
	"github.com/spf13/cobra"
	"os"
)

var (
	Cmd = &cobra.Command{
		Use:   "status",
		Short: "Status (running or stopped) of a Liferay Tomcat bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	isRunning, pid, err := procutil.IsCatalinaRunning()
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if isRunning {
		fmt.Print("Liferay is ")
		printutil.Success("RUNNING")
		fmt.Printf(" (pid=%v)\n", pid)
	} else {
		fmt.Print("Liferay is ")
		printutil.Danger("STOPPED\n")
	}
}
