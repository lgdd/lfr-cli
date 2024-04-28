package create

import (
	"fmt"
	"slices"

	"github.com/lgdd/lfr-cli/internal/prompt"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/helper"
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
	var sample, name string

	if len(args) == 0 {
		prompt.ForClientExtension(cmd, &sample, &name)
	}

	if len(args) == 1 {
		validateSample(args[0])
		sample = args[0]
		prompt.ForName(&name)
	}

	if len(args) == 2 {
		validateSample(args[0])
		sample = args[0]
		name = args[1]
	}

	if len(args) > 2 && (args[0] == "client-extension" || args[0] == "cx") {
		validateSample(args[1])
		sample = args[1]
		name = args[2]
	}

	scaffold.CreateClientExtension(sample, name)

	logger.PrintlnInfo("\n💡Checkout this tool to help you with client extensions development: https://github.com/bnheise/ce-cli")
}

func validateSample(sampleName string) {
	samples := helper.GetClientExtensionSampleNames()
	if !slices.Contains(samples, sampleName) {
		logger.Error(fmt.Sprintf("'%s' is not a valid sample.", sampleName))
		logger.Print("\n")
		logger.PrintlnInfo("Valid sample names are:\n")
		for _, sample := range samples {
			logger.Println("• " + sample)
		}
		logger.Print("\n")
		logger.Fatal("wrong sample name")
	}
}
