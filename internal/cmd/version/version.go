// Package version implements the version subcommand and exposes the build-time
// version variables injected by GoReleaser via ldflags.
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'version' displaying version, commit and date infos.
	Cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Liferay CLI",
		Run:   getVersion,
	}
	// Number is the semantic version string, injected at build time.
	Number string
	// Commit is the Git commit hash, injected at build time.
	Commit string
	// Date is the build date, injected at build time.
	Date string
)

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", Number)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Date: %s\n", Date)
}
