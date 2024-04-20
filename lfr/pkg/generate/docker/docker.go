package docker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lgdd/lfr-cli/lfr/pkg/project"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
	"github.com/magiconair/properties"
)

// DockerData contains the data to be injected into the template files
type DockerData struct {
	Image       string
	JavaVersion string
}

// Generate opinionated Docker and Docker Compose files
// for a give Java version with an option for Docker multi-stage build
func Generate(liferayWorkspace string, multistage bool, java int) error {
	projectType := project.Gradle

	if !project.IsSupportedJavaVersion(java) {
		return errors.New("java " + strconv.Itoa(java) + " is not supported\n")
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		projectType = project.Maven
	}

	tplDockerfile := "tpl/docker/Dockerfile"
	tplDockerCompose := "tpl/docker/docker-compose.yml"
	destDockerfile := filepath.Join(liferayWorkspace, "Dockerfile")
	destDockerCompose := filepath.Join(liferayWorkspace, "docker-compose.yml")

	if multistage {
		tplDockerfile = "tpl/docker/multistage/Dockerfile." + projectType
	}

	dockerImage, err := getLiferayDockerImage(liferayWorkspace)

	if err != nil {
		return err
	}

	fileutil.CopyFromTemplates(tplDockerfile, destDockerfile)
	fileutil.CopyFromTemplates(tplDockerCompose, destDockerCompose)
	err = fileutil.UpdateWithData(destDockerfile, &DockerData{
		Image:       dockerImage,
		JavaVersion: strconv.Itoa(java),
	})

	if err != nil {
		return err
	}

	fileutil.CreateDirs("build/docker/deploy")
	fileutil.CreateDirs("build/docker/configs/local")
	fileutil.CreateDirs("build/docker/configs/dev")
	fileutil.CreateDirs("build/docker/configs/uat")
	fileutil.CreateDirs("build/docker/configs/prod")

	_ = filepath.Walk("build/docker",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			printutil.Success("created ")
			fmt.Printf("%s\n", path)
			return nil
		})

	printutil.Success("created ")
	fmt.Println("Dockerfile")
	printutil.Success("created ")
	fmt.Println("docker-compose.yml")

	fmt.Print("\nTo get started:\n\n")
	fmt.Println("Deploy your modules with:")
	printutil.Info("lfr deploy\n\n")
	fmt.Println("Start your docker containers with:")
	printutil.Info("docker-compose up -d\n\n")
	fmt.Println("And follow the logs with:")
	printutil.Info("docker-compose logs -f\n\n")

	return nil
}

func getLiferayDockerImage(workspacePath string) (string, error) {
	if fileutil.IsMavenWorkspace(workspacePath) {
		pomWorkspacePath := filepath.Join(workspacePath, "pom.xml")
		pomWorkspace, err := os.Open(pomWorkspacePath)
		if err != nil {
			return "", err
		}
		defer pomWorkspace.Close()
		byteValue, _ := io.ReadAll(pomWorkspace)

		var pom project.WorkspacePom
		err = xml.Unmarshal(byteValue, &pom)
		if err != nil {
			return "", err
		}
		return pom.Properties.LiferayDockerImage, nil
	}

	if fileutil.IsGradleWorkspace(workspacePath) {
		gradlePropsPath := filepath.Join(workspacePath, "gradle.properties")
		gradleProps := properties.MustLoadFile(gradlePropsPath, properties.UTF8)
		dockerImage := gradleProps.GetString("liferay.workspace.docker.image.liferay", "liferay/portal:7.4.3.30-ga30")
		return dockerImage, nil
	}

	return "", nil
}
