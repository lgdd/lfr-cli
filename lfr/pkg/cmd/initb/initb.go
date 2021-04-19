package initb

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

var (
	Cmd = &cobra.Command{
		Use:   "init",
		Short: "Shortcut to initialize your Liferay bundle",
		Args:  cobra.NoArgs,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if fileutil.IsGradleWorkspace() {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec initBundle\n\n")
		exec.RunWrapperCmd([]string{"initBundle"})
	}
	if fileutil.IsMavenWorkspace() {
		fmt.Print("\nRunning ")
		printutil.Info("lfr exec bundle-support:init\n\n")
		exec.RunWrapperCmd([]string{"bundle-support:init"})
	}
}
