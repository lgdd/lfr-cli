package build

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
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
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec build\n\n")
		exec.RunWrapperCmd([]string{"build"})
	}
	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec clean install\n\n")
		exec.RunWrapperCmd([]string{"clean", "install"})
	}
}
