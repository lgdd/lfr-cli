package deploy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

var (
	Cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Shortcut to deploy your modules using Gradle or Maven",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if fileutil.IsGradleWorkspace() {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec deploy\n\n")
		exec.RunWrapperCmd([]string{"deploy"})
	}
	if fileutil.IsMavenWorkspace() {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec package\n\n")
		exec.RunWrapperCmd([]string{"package"})
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec bundle-support:deploy\n\n")
		exec.RunWrapperCmd([]string{"bundle-support:deploy"})
	}
}
