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
	TemplateEngine string
)

func init() {
	createSpringPortlet.Flags().StringVarP(&TemplateEngine, "template", "t", "thymeleaf", "template engine (thymeleaf or jsp)")
}

func generateSpringPortlet(cmd *cobra.Command, args []string) {
	name := args[0]
	spring.Generate(name, TemplateEngine)
}
