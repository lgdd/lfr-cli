// Package initb implements the init subcommand, which runs the Gradle or Maven
// wrapper task to download and initialize the Liferay bundle.
package initb

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

var (
	// Cmd is the command 'init' which initialize a Liferay bundle
	Cmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Shortcut to initialize your Liferay bundle",
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
		logger.PrintlnInfo("lfr exec initBundle\n")
		exec.RunWrapperCmd([]string{"initBundle"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		logger.Print("\nRunning ")
		logger.PrintlnInfo("lfr exec bundle-support:init\n")
		exec.RunWrapperCmd([]string{"bundle-support:init"})
	}
}
