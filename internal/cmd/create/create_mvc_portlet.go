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
		Run:     generateMvcPortlet,
	}
)

func generateMvcPortlet(cmd *cobra.Command, args []string) {
	name := args[0]
	scaffold.CreateModuleMVC(name)
}
