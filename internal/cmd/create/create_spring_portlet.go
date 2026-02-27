package create

import (
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	createSpringPortlet = &cobra.Command{
		Use:     "spring-mvc-portlet NAME",
		Aliases: []string{"spring"},
		Args:    cobra.ExactArgs(1),
		RunE:    generateSpringPortlet,
	}
	// TemplateEngine holds the option for the Spring template engine to use
	TemplateEngine string
)

func init() {
	createSpringPortlet.Flags().StringVarP(&TemplateEngine, "template", "t", "thymeleaf", "template engine (thymeleaf or jsp)")
}

func generateSpringPortlet(cmd *cobra.Command, args []string) error {
	return scaffold.CreateModuleSpring(args[0], TemplateEngine)
}
