package create

import (
	"fmt"
	"github.com/lgdd/liferay-cli/lfr/pkg/generate/docker"
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
	if Java == 8 || Java == 11 {
		docker.Generate(MultiStage, Java)
	} else {
		printutil.Danger(fmt.Sprintf("Java %v is not supported\n", Java))
		os.Exit(1)
	}
}
