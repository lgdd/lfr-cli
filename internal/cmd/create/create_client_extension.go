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
		RunE:    generateClientExtension,
		Long: `Available Liferay 7.4 U45+/GA45+
Client extensions extend Liferay without using OSGi modules.
Learn more: https://learn.liferay.com/w/dxp/building-applications/client-extensions
Samples available: https://github.com/liferay/liferay-portal/tree/master/workspaces/liferay-sample-workspace/client-extensions
		`,
	}
)

func generateClientExtension(cmd *cobra.Command, args []string) error {
	var sample, name string

	if len(args) == 0 {
		prompt.ForClientExtension(cmd, &sample, &name)
	}

	if len(args) == 1 {
		if err := validateSample(args[0]); err != nil {
			return err
		}
		sample = args[0]
		prompt.ForName(&name)
	}

	if len(args) == 2 {
		if err := validateSample(args[0]); err != nil {
			return err
		}
		sample = args[0]
		name = args[1]
	}

	if len(args) > 2 && (args[0] == "client-extension" || args[0] == "cx") {
		if err := validateSample(args[1]); err != nil {
			return err
		}
		sample = args[1]
		name = args[2]
	}

	if len(args) > 0 {
		if err := scaffold.CreateClientExtension(sample, name); err != nil {
			return err
		}
		logger.PrintlnInfo("\n💡Checkout this tool to help you with client extensions development: https://github.com/bnheise/ce-cli")
	}
	return nil
}

func validateSample(sampleName string) error {
	samples := helper.GetClientExtensionSampleNames()
	if !slices.Contains(samples, sampleName) {
		logger.Error(fmt.Sprintf("'%s' is not a valid sample.", sampleName))
		logger.Print("\n")
		logger.PrintlnInfo("Valid sample names are:\n")
		for _, sample := range samples {
			logger.Println("• " + sample)
		}
		logger.Print("\n")
		return fmt.Errorf("'%s' is not a valid sample", sampleName)
	}
	return nil
}
