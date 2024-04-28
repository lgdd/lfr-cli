package build

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

var (
	// Cmd is the command 'build' which builds a Liferay bundle
	Cmd = &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "Shortcut to build your Liferay bundle",
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
		logger.PrintlnInfo("lfr exec build\n")
		exec.RunWrapperCmd([]string{"build"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec clean install\n")
		exec.RunWrapperCmd([]string{"clean", "install"})
	}
}
