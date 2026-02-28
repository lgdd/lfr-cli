// Package trial implements the trial subcommand, which downloads a Liferay DXP
// trial license key into the specified directory.
package trial

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the command 'trial' which downloads a DXP trial license key.
	Cmd = &cobra.Command{
		Use:     "trial",
		Aliases: []string{"t"},
		Short:   "Get a DXP trial key to be placed in the current directory",
		Args:    cobra.NoArgs,
		Run:     run,
	}
	// Directory is the destination path where the trial license key is saved.
	Directory string
)

func init() {
	conf.Init()
	Cmd.Flags().StringVarP(&Directory, "directory", "d", ".", "--directory")
}

func run(cmd *cobra.Command, args []string) {
	destination := path.Join(Directory, "trial.xml")
	resp, error := http.Get("https://raw.githubusercontent.com/lgdd/liferay-product-info/refs/heads/main/dxp-trial/trial.xml")

	if error != nil {
		logger.Fatal(error.Error())
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if _, error = os.Stat(destination); error == nil {
		logger.PrintfWarn("%s already exists", destination)
	} else {
		os.WriteFile(destination, body, 0666)
		logger.PrintSuccess("created ")
		logger.Print(destination)
	}
}
