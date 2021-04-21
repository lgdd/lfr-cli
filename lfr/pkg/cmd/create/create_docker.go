package create

import (
	"fmt"
	"github.com/lgdd/liferay-cli/lfr/pkg/generate/docker"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/spf13/cobra"
	"os"
)

var (
	createDocker = &cobra.Command{
		Use:  "docker",
		Args: cobra.NoArgs,
		Run:  generateDocker,
	}
	MultiStage bool
	Java       int
)

func init() {
	createDocker.Flags().BoolVarP(&MultiStage, "multi-stage", "m", false, "--multi-stage")
	createDocker.Flags().IntVarP(&Java, "java", "j", 8, "--java")
}
func generateDocker(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if Java == 8 || Java == 11 {
		err := docker.Generate(liferayWorkspace, MultiStage, Java)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
	} else {
		printutil.Danger(fmt.Sprintf("Java %v is not supported\n", Java))
		os.Exit(1)
	}
}
