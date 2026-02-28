// Package deploy implements the deploy subcommand, which runs the Gradle or
// Maven wrapper to deploy modules to the running Liferay bundle.
package deploy

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

var (
	// Cmd is the command 'deploy' which runs modules deployment to the Liferay bundle
	Cmd = &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d"},
		Short:   "Shortcut to deploy your modules using Gradle or Maven",
		Args:    cobra.NoArgs,
		Run:     run,
	}

	// DeployClean runs a clean task before deploying when set to true.
	DeployClean bool
)

func init() {
	conf.Init()
	defaultDeployClean := viper.GetBool(conf.DeployClean)
	Cmd.Flags().BoolVarP(&DeployClean, "clean", "c", defaultDeployClean, "run a clean deploy command")
}

func run(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		cmdArgs := []string{"deploy"}
		if DeployClean {
			cmdArgs = []string{"clean", "deploy"}
		}
		logger.Print("\nRunning ")
		logger.PrintfInfo("lfr exec %s\n", strings.Join(cmdArgs, " "))
		exec.RunWrapperCmd(cmdArgs)
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		cmdArgs := []string{"package"}
		if DeployClean {
			cmdArgs = []string{"clean", "package"}
		}
		logger.Print("\nRunning ")
		logger.PrintfInfo("lfr exec %s\n", strings.Join(cmdArgs, " "))
		exec.RunWrapperCmd(cmdArgs)
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec bundle-support:deploy\n")
		exec.RunWrapperCmd([]string{"bundle-support:deploy"})
	}
}
