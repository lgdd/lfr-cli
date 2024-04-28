package create

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	config.Init()
	defaultMultistage := viper.GetBool(config.DockerMultistage)
	defaultJDK := viper.GetInt(config.DockerJDK)
	createDocker.Flags().BoolVarP(&MultiStage, "multi-stage", "m", defaultMultistage, "use multi-stage build")
	createDocker.Flags().IntVarP(&Java, "java", "j", defaultJDK, "Java version (8 or 11)")
}
func generateDocker(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	if Java == 8 || Java == 11 {
		err := scaffold.CreateDockerFiles(liferayWorkspace, MultiStage, Java)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
	} else {
		printutil.Danger(fmt.Sprintf("Java %v is not supported\n", Java))
		os.Exit(1)
	}
}
