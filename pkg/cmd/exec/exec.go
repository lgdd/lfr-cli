package exec

import (
	"fmt"
	"github.com/lgdd/deba/pkg/util/printutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "exec TASK... -- [TASK_FLAG]...",
		Short: "Execute Gradle or Maven task(s)",
		Args:  cobra.MinimumNArgs(1),
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	RunWrapperCmd(args)
}

func RunWrapperCmd(args []string) {
	wrapper, err := getWrapper()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	wrapperCmd := exec.Command(wrapper, args...)
	wrapperCmd.Stdout = os.Stdout
	wrapperCmd.Stderr = os.Stderr

	err = wrapperCmd.Run()

	if err != nil {
		printutil.Danger(err.Error())
		os.Exit(1)
	}
}

func getWrapper() (string, error) {
	if _, err := os.Stat("pom.xml"); !os.IsNotExist(err) {
		return getMavenWrapper()
	}
	return getGradleWrapper()
}

func getGradleWrapper() (string, error) {
	scriptName := "gradlew"

	if runtime.GOOS == "windows" {
		scriptName = "gradlew.bat"
	}

	return fileutil.FindFileInParent(scriptName)
}

func getMavenWrapper() (string, error) {
	scriptName := "mvnw"

	if runtime.GOOS == "windows" {
		scriptName = "mvnw.cmd"
	}

	return fileutil.FindFileInParent(scriptName)
}
