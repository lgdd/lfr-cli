package create

import (
	"github.com/spf13/cobra"
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
	Cmd.AddCommand(createApiModule)
	Cmd.AddCommand(createServiceBuilder)
	Cmd.AddCommand(createDocker)
}
