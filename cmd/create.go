package cmd

import (
	"github.com/spf13/cobra"
)

var (
	create = &cobra.Command{
		Use:   "create [type] [name]",
		Short: "Create a Liferay project",
	}

	Version string
	Build   string
)

func init() {
	create.PersistentFlags().StringVarP(&Version, "version", "v", "7.3", "--version 7.3")
	create.PersistentFlags().StringVarP(&Build, "build", "b", "gradle", "--build gradle")
	create.AddCommand(createWorkspace)
}
