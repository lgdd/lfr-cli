package initb

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
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
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec initBundle\n\n")
		exec.RunWrapperCmd([]string{"initBundle"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec bundle-support:init\n\n")
		exec.RunWrapperCmd([]string{"bundle-support:init"})
	}
}
