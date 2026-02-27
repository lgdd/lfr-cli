package create

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

var (
	createServiceBuilder = &cobra.Command{
		Use:     "service-builder NAME",
		Aliases: []string{"sb"},
		Args:    cobra.ExactArgs(1),
		RunE:    generateServiceBuilder,
	}
)

func generateServiceBuilder(cmd *cobra.Command, args []string) error {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return err
	}
	name := strcase.ToKebab(strings.ToLower(args[0]))
	return scaffold.CreateModuleServiceBuilder(liferayWorkspace, name)
}
