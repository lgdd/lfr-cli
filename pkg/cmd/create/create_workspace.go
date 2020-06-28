package create

import (
	"fmt"
	"github.com/lgdd/deba/pkg/generate/workspace"
	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	createWorkspace = &cobra.Command{
		Use:     "workspace [name]",
		Aliases: []string{"ws"},
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}
)

func run(cmd *cobra.Command, args []string) {
	name := args[0]
	err := workspace.Generate(name, Build, Version)
	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}
	printutil.Success(fmt.Sprintf("Successfully created a Liferay Workspace '%s' ", name))
	printInitCmd(name, Build)
}

func printInitCmd(name, build string) {
	fmt.Println("\nInitialize your Liferay bundle:")
	if runtime.GOOS == "windows" {
		switch build {
		case workspace.Gradle:
			printutil.Info(fmt.Sprintf("cd %s && gradlew.bat initBundle", name))
		case workspace.Maven:
			printutil.Info(fmt.Sprintf("cd %s && mvnw.cmd bundle-support:init", name))
		}
	} else {
		switch build {
		case workspace.Gradle:
			printutil.Info(fmt.Sprintf("cd %s && ./gradlew initBundle", name))
		case workspace.Maven:
			printutil.Info(fmt.Sprintf("cd %s && ./mvnw bundle-support:init", name))
		}
	}
}
