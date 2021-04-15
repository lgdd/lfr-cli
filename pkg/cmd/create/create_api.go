package create

import (
	"github.com/lgdd/deba/pkg/generate/api"
	"github.com/spf13/cobra"
)

var (
	createApiModule = &cobra.Command{
		Use:  "api NAME",
		Args: cobra.ExactArgs(1),
		Run:  generateApiModule,
	}
)

func generateApiModule(cmd *cobra.Command, args []string) {
	name := args[0]
	api.Generate(name)
}
