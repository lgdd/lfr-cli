package deploy

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
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
)

func run(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec deploy\n")
		exec.RunWrapperCmd([]string{"deploy"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec package\n")
		exec.RunWrapperCmd([]string{"package"})
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec bundle-support:deploy\n")
		exec.RunWrapperCmd([]string{"bundle-support:deploy"})
	}
}
