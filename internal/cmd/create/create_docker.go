package create

import (
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

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
	conf.Init()
	defaultMultistage := viper.GetBool(conf.DockerMultistage)
	defaultJDK := viper.GetInt(conf.DockerJDK)
	createDocker.Flags().BoolVarP(&MultiStage, "multi-stage", "m", defaultMultistage, "use multi-stage build")
	createDocker.Flags().IntVarP(&Java, "java", "j", defaultJDK, "Java version (8, 11, 17 or 21)")
}
func generateDocker(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if Java == 8 || Java == 11 {
		err := scaffold.CreateDockerFiles(liferayWorkspace, MultiStage, Java)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		logger.Fatalf("Java %v is not supported", Java)
	}
}
