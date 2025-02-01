package trial

import (
	"io"
	"net/http"
	"os"

	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "trial",
		Aliases: []string{"t"},
		Short:   "Get a DXP trial key to be placed in the current directory",
		Args:    cobra.NoArgs,
		Run:     run,
	}
	Follow bool
)

func init() {
	conf.Init()
}

func run(cmd *cobra.Command, args []string) {
	resp, error := http.Get("https://raw.githubusercontent.com/lgdd/liferay-product-info/refs/heads/main/dxp-trial/trial.xml")

	if error != nil {
		logger.Fatal(error.Error())
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if _, error = os.Stat("trial.xml"); error == nil {
		logger.PrintWarn("trial.xml already exists")
	} else {
		os.WriteFile("trial.xml", body, 0666)
		logger.PrintSuccess("created ")
		logger.Print("trial.xml")
	}
}
