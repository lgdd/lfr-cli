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
	Number string
	Commit string
	Date   string
)

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", Number)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Date: %s\n", Date)
}
