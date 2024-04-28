package start

import (
	"os"
	"os/exec"

	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/lgdd/lfr-cli/pkg/util/procutil"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'start' which allows to start the Liferay bundle
	Cmd = &cobra.Command{
		Use:   "start",
		Short: "Start a Liferay Tomcat bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	startupScript, err := fileutil.GetTomcatScriptPath("startup")

	if err != nil {
		logger.Fatal(err.Error())
	}

	tomcatPath, err := fileutil.GetTomcatPath()

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = os.Setenv("CATALINA_HOME", tomcatPath)

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = procutil.SetCatalinaPid()

	if err != nil {
		logger.Fatal(err.Error())
	}

	startupCmd := exec.Command(startupScript)
	startupCmd.Stdout = os.Stdout

	err = startupCmd.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Print("\nFollow the logs:")
	logger.PrintInfo("lfr logs -f\n")
}
