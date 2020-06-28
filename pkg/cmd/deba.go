package cmd

import (
	"github.com/lgdd/deba/pkg/cmd/create"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var deba = &cobra.Command{
	Use:   "deba",
	Short: "Deba is yet another command-line tool to work on Liferay projects.",
}

func init() {
	deba.AddCommand(create.Cmd)
	deba.PersistentFlags().BoolVar(&printutil.NoColor, "no-color", false, "--no-color (disable color output)")
}

func Execute() error {
	return deba.Execute()
}
