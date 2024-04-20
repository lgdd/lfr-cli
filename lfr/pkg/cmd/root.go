package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/build"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/completion"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/create"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/deploy"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/diagnose"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/initb"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/logs"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/shell"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/start"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/status"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/stop"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/update"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/version"
	"github.com/lgdd/lfr-cli/lfr/pkg/config"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
)

var root = &cobra.Command{
	Use:   "lfr",
	Short: "Liferay CLI (lfr) is an unofficial tool written in Go that helps you manage your Liferay projects.",
}

func init() {
	cobra.OnInitialize(config.Init)

	root.AddCommand(completion.Cmd)
	root.AddCommand(create.Cmd)
	root.AddCommand(exec.Cmd)
	root.AddCommand(build.Cmd)
	root.AddCommand(deploy.Cmd)
	root.AddCommand(initb.Cmd)
	root.AddCommand(start.Cmd)
	root.AddCommand(stop.Cmd)
	root.AddCommand(status.Cmd)
	root.AddCommand(logs.Cmd)
	root.AddCommand(shell.Cmd)
	root.AddCommand(version.Cmd)
	root.AddCommand(update.Cmd)
	root.AddCommand(diagnose.Cmd)
	root.PersistentFlags().BoolVar(&printutil.NoColor, "no-color", false, "disable colors for output messages")
}

// Run the the main command
func Execute() error {
	return root.Execute()
}
