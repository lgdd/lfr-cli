package cmd

import (
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/completion"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/create"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/deploy"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/initb"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/logs"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/shell"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/start"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/status"
	"github.com/lgdd/liferay-cli/lfr/pkg/cmd/stop"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "lfr",
	Short: "Liferay CLI (lfr) is an unofficial tool written in Go that helps you manage your Liferay projects.",
}

func init() {
	root.AddCommand(completion.Cmd)
	root.AddCommand(create.Cmd)
	root.AddCommand(exec.Cmd)
	root.AddCommand(deploy.Cmd)
	root.AddCommand(initb.Cmd)
	root.AddCommand(start.Cmd)
	root.AddCommand(stop.Cmd)
	root.AddCommand(status.Cmd)
	root.AddCommand(logs.Cmd)
	root.AddCommand(shell.Cmd)
	root.PersistentFlags().BoolVar(&printutil.NoColor, "no-color", false, "--no-color (disable color output)")
}

func Execute() error {
	return root.Execute()
}
