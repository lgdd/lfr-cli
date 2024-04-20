package exec

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"

	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'exec' which executes a task using Gradle or Maven wrapper
	Cmd = &cobra.Command{
		Use:     "exec TASK... -- [TASK_FLAG]...",
		Aliases: []string{"x"},
		Short:   "Execute Gradle or Maven task(s)",
		Args:    cobra.MinimumNArgs(1),
		Run:     run,
	}
)

func run(cmd *cobra.Command, args []string) {
	RunWrapperCmd(args)
}

// Get the Maven or Gradle wrapper to execute tasks
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
