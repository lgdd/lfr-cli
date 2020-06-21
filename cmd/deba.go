package cmd

import (
	"github.com/spf13/cobra"
)

var deba = &cobra.Command{
	Use:   "deba",
	Short: "Deba is yet another command-line tool to work on Liferay projects.",
}

func init() {
	deba.AddCommand(create)
}

func Execute() error {
	return deba.Execute()
}
