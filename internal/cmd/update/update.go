// Package update implements the update subcommand, which performs a self-update
// of the lfr binary to the latest GitHub release.
package update

import (
	"github.com/blang/semver"
	"github.com/lgdd/lfr-cli/internal/cmd/version"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'update' which allows the tool to update itself
	Cmd = &cobra.Command{
		Use:   "update",
		Short: "Update Liferay CLI to the latest version",
		Run:   doSelfUpdate,
	}
)

func doSelfUpdate(cmd *cobra.Command, args []string) {
	v := semver.MustParse(version.Number)
	latest, err := selfupdate.UpdateSelf(v, "lgdd/liferay-cli")
	if err != nil {
		logger.Fatal(err.Error())
	}
	if latest.Version.Equals(v) {
		logger.Printf("Current binary (%v) is the latest version (%v)\n", version.Number, latest.Version)
	} else {
		logger.PrintfSuccess("Successfully updated to version %s", latest.Version)
		logger.Printf("Release note:\n%s", latest.ReleaseNotes)
	}
}
