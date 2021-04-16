package create

import (
	"github.com/lgdd/liferay-cli/lfr/pkg/generate/mvc"
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
	mvc.Generate(name)
}
