package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/lgdd/lfr-cli/pkg/generate/sb"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
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
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	name := args[0]
	name = strcase.ToKebab(strings.ToLower(name))
	sb.Generate(liferayWorkspace, name)
}
