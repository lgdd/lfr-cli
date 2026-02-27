package create

import (
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/spf13/cobra"
)

var (
	createApiModule = &cobra.Command{
		Use:  "api NAME",
		Args: cobra.ExactArgs(1),
		RunE: generateApiModule,
	}
)

func generateApiModule(cmd *cobra.Command, args []string) error {
	return scaffold.CreateModuleAPI(args[0])
}
