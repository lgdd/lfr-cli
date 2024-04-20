package create

import (
	"github.com/lgdd/lfr-cli/lfr/pkg/generate/cx"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createClientExtension = &cobra.Command{
		Use:     "client-extension NAME",
		Aliases: []string{"cx"},
		Run:     generateClientExtension,
		Long: `Available Liferay 7.4 U45+/GA45+
Client extensions extend Liferay without using OSGi modules.
Learn more: https://learn.liferay.com/w/dxp/building-applications/client-extensions
Samples available: https://github.com/liferay/liferay-portal/tree/master/workspaces/liferay-sample-workspace/client-extensions
		`,
	}
)

func generateClientExtension(cmd *cobra.Command, args []string) {
	cx.Generate(cmd, args)
	printutil.Info("\n💡Checkout this tool to help you with client extensions development: https://github.com/bnheise/ce-cli\n")
}
