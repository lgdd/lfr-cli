package create

import (
	"github.com/lgdd/liferay-cli/lfr/pkg/generate/gogocmd"
	"github.com/spf13/cobra"
)

var (
	createCmdModule = &cobra.Command{
		Use:     "command NAME",
		Aliases: []string{"cmd"},
		Args:    cobra.ExactArgs(1),
		Run:     generateCmdModule,
	}
)

func generateCmdModule(cmd *cobra.Command, args []string) {
	name := args[0]
	gogocmd.Generate(name)
}
