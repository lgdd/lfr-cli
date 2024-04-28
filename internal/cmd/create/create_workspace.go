package create

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	config.Init()

	defaultVersion := viper.GetString(config.WorkspaceVersion)
	defaultBuild := viper.GetString(config.WorkspaceBuild)
	defaultEdition := viper.GetString(config.WorkspaceEdition)
	defaultInit := viper.GetBool(config.WorkspaceInit)

	createWorkspace.Flags().StringVarP(&Version, "version", "v", defaultVersion, "Liferay major version (7.x)")
	createWorkspace.Flags().StringVarP(&Build, "build", "b", defaultBuild, "build tool (gradle or maven)")
	createWorkspace.Flags().StringVarP(&Edition, "edition", "e", defaultEdition, "Liferay edition (dxp or portal)")
	createWorkspace.Flags().BoolVarP(&Init, "init", "i", defaultInit, "executes Liferay bundle initialization (i.e. download & unzip in the workspace)")
}

func generateWorkspace(cmd *cobra.Command, args []string) {
	if fileutil.IsInWorkspaceDir() {
		logger.Fatalf("You're already in a Liferay Workspace.")
	}
	name := args[0]
	err := scaffold.CreateWorkspace(name, Build, Version, Edition)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.PrintfSuccess("\nSuccessfully created a Liferay Workspace '%s'\n", name)

	if Init {
		runInit(name, Build)
	} else {
		printInitCmd(name)
	}
}

func runInit(name, build string) {
	if err := os.Chdir(name); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Print("\nInitializing Liferay Bundle using:\n\n")

	switch build {
	case metadata.Gradle:
		logger.PrintlnInfo("lfr exec initBundle\n")
		exec.RunWrapperCmd([]string{"initBundle"})
	case metadata.Maven:
		logger.PrintlnInfo("lfr exec bundle-support:init\n")
		exec.RunWrapperCmd([]string{"bundle-support:init"})
	}
}

func printInitCmd(name string) {
	logger.Println("\nInitialize your Liferay bundle:")
	logger.PrintlnInfo(fmt.Sprintf("cd %s && lfr init\n", name))
}
