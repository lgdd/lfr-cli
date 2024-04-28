package create

import (
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

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
	scaffold.CreateClientExtension(cmd, args)
	logger.PrintInfo("\n💡Checkout this tool to help you with client extensions development: https://github.com/bnheise/ce-cli")
}
