package create

import (
	"github.com/spf13/cobra"

	"github.com/lgdd/liferay-cli/lfr/pkg/generate/sb"
)

var (
	createServiceBuilder = &cobra.Command{
		Use:     "service-builder NAME",
		Aliases: []string{"sb"},
		Args:    cobra.ExactArgs(1),
		Run:     generateServiceBuilder,
	}
)

func generateServiceBuilder(cmd *cobra.Command, args []string) {
	name := args[0]
	sb.Generate(name)
}
