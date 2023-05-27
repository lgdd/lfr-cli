package deploy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
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
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec deploy\n\n")
		exec.RunWrapperCmd([]string{"deploy"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec package\n\n")
		exec.RunWrapperCmd([]string{"package"})
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec bundle-support:deploy\n\n")
		exec.RunWrapperCmd([]string{"bundle-support:deploy"})
	}
}
