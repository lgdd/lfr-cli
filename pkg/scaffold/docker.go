package scaffold

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/helper"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// DockerData contains the data to be injected into the template files
type DockerData struct {
	Image       string
	JavaVersion string
}

// Create opinionated Docker and Docker Compose files
// for a give Java version with an option for Docker multi-stage build
func CreateDockerFiles(liferayWorkspace string, multistage bool, java int) error {
	projectType := metadata.Gradle

	if !helper.IsSupportedJavaVersion(java) {
		return errors.New("java " + strconv.Itoa(java) + " is not supported\n")
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		projectType = metadata.Maven
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

	fileutil.CreateDirs(filepath.Join(liferayWorkspace, "build/docker/deploy"))
	fileutil.CreateDirs(filepath.Join(liferayWorkspace, "build/docker/configs/local"))
	fileutil.CreateDirs(filepath.Join(liferayWorkspace, "build/docker/configs/dev"))
	fileutil.CreateDirs(filepath.Join(liferayWorkspace, "build/docker/configs/uat"))
	fileutil.CreateDirs(filepath.Join(liferayWorkspace, "build/docker/configs/prod"))

	_ = filepath.Walk("build/docker",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			logger.PrintSuccess("created ")
			fmt.Printf("%s\n", path)
			return nil
		})

	logger.PrintSuccess("created ")
	logger.Println("Dockerfile")
	logger.PrintSuccess("created ")
	logger.Println("docker-compose.yml")

	logger.Print("\nTo get started:\n\n")
	logger.Println("Deploy your modules with:")
	logger.PrintlnInfo("lfr deploy\n")
	logger.Println("Start your docker containers with:")
	logger.PrintlnInfo("docker-compose up -d\n")
	logger.Println("And follow the logs with:")
	logger.PrintlnInfo("docker-compose logs -f\n")

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

		var pom fileutil.WorkspacePom
		err = xml.Unmarshal(byteValue, &pom)
		if err != nil {
			return "", err
		}
		return pom.Properties.LiferayDockerImage, nil
	}

	if fileutil.IsGradleWorkspace(workspacePath) {
		dockerImageRegex := regexp.MustCompile(`(liferay\/)(portal|dxp):.+`)
		gradlePropsPath := filepath.Join(workspacePath, "gradle.properties")
		gradlePropsBytes, err := os.ReadFile(gradlePropsPath)

		if err != nil {
			return "", err
		}
		return dockerImageRegex.FindString(string(gradlePropsBytes)), nil
	}

	return "", nil
}
