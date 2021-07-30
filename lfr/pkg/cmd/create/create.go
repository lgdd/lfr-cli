package create

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/liferay-cli/lfr/pkg/project"
)

var (
	Cmd = &cobra.Command{
		Use:   "create TYPE NAME",
		Short: "Create a Liferay project",
	}
)

func init() {
	Cmd.AddCommand(createWorkspace)
	Cmd.AddCommand(createMvcPortlet)
	Cmd.AddCommand(createSpringPortlet)
	Cmd.AddCommand(createApiModule)
	Cmd.AddCommand(createServiceBuilder)
	Cmd.AddCommand(createDocker)
	Cmd.PersistentFlags().StringVarP(&project.PackageName, "package", "p", "org.acme", "base package name")
}
