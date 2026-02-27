package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/lgdd/lfr-cli/pkg/util/procutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createWorkspace = &cobra.Command{
		Use:     "workspace NAME",
		Aliases: []string{"ws"},
		Args:    cobra.ExactArgs(1),
		RunE:    generateWorkspace,
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
	conf.Init()
	defaultVersion := viper.GetString(conf.WorkspaceVersion)
	defaultBuild := viper.GetString(conf.WorkspaceBuild)
	defaultEdition := viper.GetString(conf.WorkspaceEdition)
	defaultInit := viper.GetBool(conf.WorkspaceInit)

	createWorkspace.Flags().StringVarP(&Version, "version", "v", defaultVersion, "Liferay major version (7.x)")
	createWorkspace.Flags().StringVarP(&Build, "build", "b", defaultBuild, "build tool (gradle or maven)")
	createWorkspace.Flags().StringVarP(&Edition, "edition", "e", defaultEdition, "Liferay edition (dxp or portal)")
	createWorkspace.Flags().BoolVarP(&Init, "init", "i", defaultInit, "executes Liferay bundle initialization (i.e. download & unzip in the workspace)")
}

func generateWorkspace(cmd *cobra.Command, args []string) error {
	if fileutil.IsInWorkspaceDir() {
		return fmt.Errorf("you're already in a Liferay Workspace")
	}
	name := args[0]
	if err := scaffold.CreateWorkspace(name, Build, Version, Edition); err != nil {
		return err
	}
	logger.PrintfSuccess("\nSuccessfully created a Liferay Workspace '%s'\n", name)

	if Init {
		return runInit(name, Build)
	}
	printInitCmd(name)
	return nil
}

func runInit(name, build string) error {
	if err := os.Chdir(name); err != nil {
		return err
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
	return nil
}

func printInitCmd(name string) {
	logger.Println("\nInitialize your Liferay bundle:")
	logger.PrintlnInfo(fmt.Sprintf("cd %s && lfr init\n", name))
	major, version, err := procutil.GetCurrentJavaVersion()

	if err != nil {
		logger.PrintError("[✗] ")
		logger.Print("Liferay requires Java 8, 11, 17 or 21.\n")
	}

	if major == "8" {
		logger.PrintWarn("[!] ")
		logger.Printf("Java %s is used\n", strings.Split(version, "\n")[0])
		logger.PrintWarn("Liferay DXP 2024.Q1 and Liferay Portal 7.4 GA112 will be the last version to support Java 8.")
	} else if major == "11" {
		logger.PrintWarn("[!] ")
		logger.Printf("Java %s is used\n", strings.Split(version, "\n")[0])
		logger.PrintWarn("Liferay DXP DXP 2024.Q2 and Liferay Portal 7.4 GA120 will be the last version to support Java 11.")
	} else if major == "17" || major == "21" {
		logger.PrintSuccess("[✓] ")
		logger.Printf("Java %s is used\n", strings.Split(version, "\n")[0])
		logger.PrintlnBold("Liferay DXP 2024.Q2+ and Liferay Portal CE 7.4 GA120+ are fully certified to run on both Java JDK 17 and 21.")
		logger.PrintfBold("Building with Java %s requires Gradle 8.5 and Liferay Gradle Workspace Plugin v10.1+.", major)
	} else {
		logger.PrintWarn("[!] ")
		logger.Printf("Java %s is used\n", strings.Split(version, "\n")[0])
		logger.PrintlnBold("Liferay only supports Java 8, 11, 17 and 21. See https://www.liferay.com/compatibility-matrix.")
	}
}
