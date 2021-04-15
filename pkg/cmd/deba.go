package cmd

import (
	"github.com/lgdd/deba/pkg/cmd/create"
	"github.com/lgdd/deba/pkg/cmd/exec"
	"github.com/lgdd/deba/pkg/cmd/logs"
	"github.com/lgdd/deba/pkg/cmd/start"
	"github.com/lgdd/deba/pkg/cmd/status"
	"github.com/lgdd/deba/pkg/cmd/stop"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var deba = &cobra.Command{
	Use:   "deba",
	Short: "Deba is yet another command-line tool to work on Liferay projects.",
}

func init() {
	deba.AddCommand(create.Cmd)
	deba.AddCommand(exec.Cmd)
	deba.AddCommand(start.Cmd)
	deba.AddCommand(stop.Cmd)
	deba.AddCommand(status.Cmd)
	deba.AddCommand(logs.Cmd)
	deba.PersistentFlags().BoolVar(&printutil.NoColor, "no-color", false, "--no-color (disable color output)")
}

func Execute() error {
	return deba.Execute()
}
