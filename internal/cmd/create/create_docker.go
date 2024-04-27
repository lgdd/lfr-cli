package create

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/pkg/generate/docker"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createDocker = &cobra.Command{
		Use:  "docker",
		Args: cobra.NoArgs,
		Run:  generateDocker,
	}
	// MultiStage holds the option to create a Dockerfile with multi-stage build
	MultiStage bool
	// Java is the Java version to use in the Dockerfile
	Java int
)

func init() {
	createDocker.Flags().BoolVarP(&MultiStage, "multi-stage", "m", false, "use multi-stage build")
	createDocker.Flags().IntVarP(&Java, "java", "j", 11, "Java version (8 or 11)")
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
