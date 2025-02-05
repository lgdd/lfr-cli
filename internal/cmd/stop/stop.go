package stop

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/lgdd/lfr-cli/pkg/util/procutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Cmd is the command 'stop' which shutdowns the Liferay bundle
	Cmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a Liferay Tomcat bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		logger.PrintInfo("not available for Windows since the Tomcat process run in another command window")
		os.Exit(0)
	}
	shutdownScript, err := fileutil.GetTomcatScriptPath("shutdown")

	if err != nil {
		logger.Fatal(err.Error())
	}

	_, err = procutil.GetCatalinaPid()

	if err != nil {
		logger.Fatal(err.Error())
	}

	shutdownCmd := exec.Command(shutdownScript)
	shutdownCmd.Stdout = os.Stdout

	err = shutdownCmd.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	isRunning, pid, err := procutil.IsCatalinaRunning()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if !isRunning {
		logger.Print("Liferay is already stopped.")
		os.Exit(0)
	}

	_ = spinner.New().
		Title(fmt.Sprintf("Stopping Liferay (pid=%v)", pid)).
		Action(checkLiferayStatus).
		Accessible(viper.GetBool(conf.OutputAccessible)).
		Run()

	isRunning, pid, err = procutil.IsCatalinaRunning()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if isRunning {
		logger.PrintlnWarn("Liferay is still running.")
		logger.PrintfWarn("Be patient, but you might need to kill the process (pid=%v).", pid)
	} else {
		logger.PrintSuccess("Liferay stopped successfully.")
	}
}

func checkLiferayStatus() {
	isRunning, _, err := procutil.IsCatalinaRunning()

	if err != nil {
		logger.Fatal(err.Error())
	}

	maxIterations := 60
	iteration := 0
	for isRunning && iteration < maxIterations {
		iteration = iteration + 1
		isRunning, _, err = procutil.IsCatalinaRunning()
		if err != nil {
			logger.Fatal(err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
