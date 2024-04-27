package create

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/pkg/generate/workspace"
	"github.com/lgdd/lfr-cli/pkg/project"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createWorkspace = &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Args:    cobra.ExactArgs(1),
		Run:     generateWorkspace,
	}
	// Version is the Liferay version
	Version string
	// Build is Maven or Gradle
	Build string
	// Edition is DXP (Enterprise) or Portal (Community)
	Edition string
	// Init holds the option to initialize the Liferay bundle
	Init bool
)

func init() {
	createWorkspace.Flags().StringVarP(&Version, "version", "v", "7.4", "Liferay major version (7.x)")
	createWorkspace.Flags().StringVarP(&Build, "build", "b", "gradle", "build tool (gradle or maven)")
	createWorkspace.Flags().StringVarP(&Edition, "edition", "e", "portal", "Liferay edition (dxp or portal)")
	createWorkspace.Flags().BoolVarP(&Init, "init", "i", false, "executes Liferay bundle initialization (i.e. download & unzip in the workspace)")
}

func generateWorkspace(cmd *cobra.Command, args []string) {
	if fileutil.IsInWorkspaceDir() {
		printutil.Danger("You're already in a Liferay Workspace and I can't create a new one in it.\n")
		os.Exit(1)
	}
	name := args[0]
	err := workspace.Generate(name, Build, Version, Edition)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	printutil.Success(fmt.Sprintf("\nSuccessfully created a Liferay Workspace '%s'\n", name))

	if Init {
		runInit(name, Build)
	} else {
		printInitCmd(name, Build)
	}
}

func runInit(name, build string) {
	if err := os.Chdir(name); err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fmt.Print("\nInitializing Liferay Bundle using:\n\n")

	switch build {
	case project.Gradle:
		printutil.Info("lfr exec initBundle\n\n")
		exec.RunWrapperCmd([]string{"initBundle"})
	case project.Maven:
		printutil.Info("lfr exec bundle-support:init\n\n")
		exec.RunWrapperCmd([]string{"bundle-support:init"})
	}
}

func printInitCmd(name, build string) {
	fmt.Println("\nInitialize your Liferay bundle:")
	printutil.Info(fmt.Sprintf("cd %s && lfr init\n", name))
	fmt.Print("\n")
}
