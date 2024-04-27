package update

import (
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/lgdd/lfr-cli/internal/cmd/version"
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
		log.Println("Binary update failed:", err)
		os.Exit(1)
	}
	if latest.Version.Equals(v) {
		log.Printf("Current binary (%v) is the latest version (%v)\n", version.Number, latest.Version)
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
	}
}
