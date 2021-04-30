package create

import (
	"github.com/lgdd/liferay-cli/lfr/pkg/generate/spring"
	"github.com/spf13/cobra"
)

var (
	createSpringPortlet = &cobra.Command{
		Use:     "spring-mvc-portlet NAME",
		Aliases: []string{"spring"},
		Args:    cobra.ExactArgs(1),
		Run:     generateSpringPortlet,
	}
)

func generateSpringPortlet(cmd *cobra.Command, args []string) {
	name := args[0]
	spring.Generate(name)
}
