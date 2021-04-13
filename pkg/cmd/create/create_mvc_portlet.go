package create

import (
	"github.com/lgdd/deba/pkg/generate/mvc_portlet"
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
	mvc_portlet.Generate(name)
}
