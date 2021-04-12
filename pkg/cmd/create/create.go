package create

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "create TYPE NAME",
		Short: "Create a Liferay project",
	}

	Version string
	Build   string
)

func init() {
	Cmd.PersistentFlags().StringVarP(&Version, "version", "v", "7.3", "--version 7.3")
	Cmd.PersistentFlags().StringVarP(&Build, "build", "b", "gradle", "--build gradle")
	Cmd.AddCommand(createWorkspace)
}
