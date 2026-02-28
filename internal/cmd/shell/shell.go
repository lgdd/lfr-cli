// Package shell implements the shell subcommand, which opens a Gogo Shell
// session connected to the running Liferay bundle.
package shell

import (
	"fmt"

	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'shell' which returns the Gogo Shell from a running Liferay bundle
	Cmd = &cobra.Command{
		Use:     "shell",
		Short:   "Connect and get Liferay Gogo Shell",
		Aliases: []string{"sh"},
		Args:    cobra.NoArgs,
		RunE:    run,
	}
	// Host is the hostname or IP address of the Gogo Shell server.
	Host string
	// Port is the TCP port of the Gogo Shell server.
	Port int
)

func init() {
	Cmd.Flags().StringVar(&Host, "host", "localhost", "--host localhost")
	Cmd.Flags().IntVarP(&Port, "port", "p", 11311, "--port 11311")
}

func run(cmd *cobra.Command, args []string) error {
	destination := fmt.Sprintf("%s:%v", Host, Port)
	fmt.Printf("Connecting to %v...\n", destination)
	err := connectGogoShell(Host, Port)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return nil
}
