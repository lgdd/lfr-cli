package create

import (
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	createMvcPortlet = &cobra.Command{
		Use:     "mvc-portlet NAME",
		Aliases: []string{"mvc"},
		Args:    cobra.ExactArgs(1),
		RunE:    generateMvcPortlet,
	}
)

func generateMvcPortlet(cmd *cobra.Command, args []string) error {
	return scaffold.CreateModuleMVC(args[0])
}
