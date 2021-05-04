package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Liferay CLI",
		Run:   getVersion,
	}
	number string
	commit string
	date   string
)

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", number)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
}
