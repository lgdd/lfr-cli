package docker

import (
	"encoding/xml"
	"fmt"
	"github.com/lgdd/liferay-cli/lfr/pkg/project"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/magiconair/properties"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type DockerData struct {
	Image       string
	JavaVersion string
}

func Generate(multistage bool, java int) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	projectType := project.Gradle

	if fileutil.IsMavenWorkspace() {
		projectType = project.Maven
	}

	tplDockerfile := "tpl/docker/Dockerfile"
	tplDockerCompose := "tpl/docker/docker-compose.yml"
	destDockerfile := filepath.Join(liferayWorkspace, "Dockerfile")
	destDockerCompose := filepath.Join(liferayWorkspace, "docker-compose.yml")

	if multistage {
		tplDockerfile = "tpl/docker/multistage/Dockerfile." + projectType
	}

	dockerImage, err := getLiferayDockerImage()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.CopyFromTemplates(tplDockerfile, destDockerfile)
	fileutil.CopyFromTemplates(tplDockerCompose, destDockerCompose)
	err = fileutil.UpdateWithData(destDockerfile, &DockerData{
		Image:       dockerImage,
		JavaVersion: strconv.Itoa(java),
	})

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.CreateDirs("build/docker/deploy")
	fileutil.CreateDirs("build/docker/configs/local")
	fileutil.CreateDirs("build/docker/configs/dev")
	fileutil.CreateDirs("build/docker/configs/uat")
	fileutil.CreateDirs("build/docker/configs/prod")

	fmt.Print("\nTo get started:\n\n")
	fmt.Println("Deploy your modules with:")
	printutil.Info("lfr deploy\n\n")
	fmt.Println("Start your docker containers with:")
	printutil.Info("docker-compose up -d\n\n")
	fmt.Println("And follow the logs with:")
	printutil.Info("docker-compose logs -f\n\n")
}

func getLiferayDockerImage() (string, error) {
	workspacePath, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return "", err
	}

	if fileutil.IsMavenWorkspace() {
		pomWorkspacePath := filepath.Join(workspacePath, "pom.xml")
		pomWorkspace, err := os.Open(pomWorkspacePath)
		if err != nil {
			return "", err
		}
		defer pomWorkspace.Close()
		byteValue, _ := ioutil.ReadAll(pomWorkspace)

		var pom project.WorkspacePom
		err = xml.Unmarshal(byteValue, &pom)
		if err != nil {
			return "", err
		}
		return pom.Properties.LiferayDockerImage, nil
	}

	if fileutil.IsGradleWorkspace() {
		gradlePropsPath := filepath.Join(workspacePath, "gradle.properties")
		gradleProps := properties.MustLoadFile(gradlePropsPath, properties.UTF8)
		dockerImage := gradleProps.GetString("liferay.workspace.docker.image.liferay", "liferay/portal:7.3.6-ga7")
		return dockerImage, nil
	}

	return "", nil
}
