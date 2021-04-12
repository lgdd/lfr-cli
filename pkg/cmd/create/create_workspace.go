package create

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lgdd/deba/pkg/generate/workspace"
	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createWorkspace = &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}
)

func run(cmd *cobra.Command, args []string) {
	fileutil.VerifyCurrentDirAsWorkspace(Build)
	name := args[0]
	err := workspace.Generate(name, Build, Version)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	printutil.Success(fmt.Sprintf("\nSuccessfully created a Liferay Workspace '%s'\n", name))
	printInitCmd(name, Build)
	fmt.Print("\n")
}

func printInitCmd(name, build string) {
	fmt.Println("\nInitialize your Liferay bundle:")
	if runtime.GOOS == "windows" {
		switch build {
		case project.Gradle:
			printutil.Info(fmt.Sprintf("cd %s && gradlew.bat initBundle\n", name))
		case project.Maven:
			printutil.Info(fmt.Sprintf("cd %s && mvnw.cmd bundle-support:init\n", name))
		}
	} else {
		switch build {
		case project.Gradle:
			printutil.Info(fmt.Sprintf("cd %s && ./gradlew initBundle\n", name))
		case project.Maven:
			printutil.Info(fmt.Sprintf("cd %s && ./mvnw bundle-support:init\n", name))
		}
	}
}
