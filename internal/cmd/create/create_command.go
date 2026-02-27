package create

import (
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	createCmdModule = &cobra.Command{
		Use:     "command NAME",
		Aliases: []string{"cmd"},
		Args:    cobra.ExactArgs(1),
		RunE:    generateCmdModule,
	}
)

func generateCmdModule(cmd *cobra.Command, args []string) error {
	return scaffold.CreateModuleGogoCommand(args[0])
}
