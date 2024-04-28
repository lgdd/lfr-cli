package create

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
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
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		logger.Fatal(err.Error())
	}
	name := args[0]
	name = strcase.ToKebab(strings.ToLower(name))
	scaffold.CreateModuleServiceBuilder(liferayWorkspace, name)
}
