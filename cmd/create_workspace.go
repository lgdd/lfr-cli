package cmd

import (
	"fmt"
	"github.com/lgdd/deba/gen"
	"github.com/lgdd/deba/util"
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
	err := gen.CreateWorkspace(name, Build, Version)
	if err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
	util.PrintSuccess(fmt.Sprintf("Successfully created a Liferay Workspace '%s' ", name))
	printInitCmd(name, Build)
}

func printInitCmd(name, build string) {
	fmt.Println("\nInitialize your Liferay bundle:")
	if runtime.GOOS == "windows" {
		switch build {
		case gen.Gradle:
			util.PrintInfo(fmt.Sprintf("cd %s && gradlew.bat initBundle", name))
		case gen.Maven:
			util.PrintInfo(fmt.Sprintf("cd %s && mvnw.cmd bundle-support:init", name))
		}
	} else {
		switch build {
		case gen.Gradle:
			util.PrintInfo(fmt.Sprintf("cd %s && ./gradlew initBundle", name))
		case gen.Maven:
			util.PrintInfo(fmt.Sprintf("cd %s && ./mvnw bundle-support:init", name))
		}
	}
}
