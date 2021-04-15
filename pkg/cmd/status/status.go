package status

import (
	"fmt"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/lgdd/deba/pkg/util/procutil"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "status",
		Short: "Status of the Liferay bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	isRunning, pid, err := procutil.IsCatalinaRunning()
	if err != nil {
		panic(err)
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
