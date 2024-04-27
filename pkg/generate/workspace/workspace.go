package workspace

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lgdd/lfr-cli/pkg/project"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/lgdd/lfr-cli/pkg/util/procutil"
)

// Build options (i.e. Maven or Gradle)
const (
	Gradle = "gradle"
	Maven  = "maven"
)

// Generate the structure of a Liferay workspace
func Generate(base, build, version, edition string) error {
	metadata, err := project.NewMetadata(base, version, edition)
	if err != nil {
		return err
	}

	if build == Maven {
		err = os.Mkdir(base, os.ModePerm)
		if err != nil {
			return err
		}
		if err := createMavenFiles(base, metadata); err != nil {
			return err
		}
		createCommonEmptyDirs(base)
	} else if build == Gradle {
		err = os.Mkdir(base, os.ModePerm)
		if err != nil {
			return err
		}
		if err := createGradleFiles(base, metadata); err != nil {
			return err
		}
		createCommonEmptyDirs(base)
	} else {
		return errors.New("only Gradle and Maven are supported")
	}

	createGithubWorkflows(base)
	procutil.Exec("git", "init", base)

	_ = filepath.Walk(base,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				printutil.Success("created ")
				fmt.Printf("%s\n", path)
			}
			return nil
		})

	return nil
}

func createGradleFiles(base string, metadata *project.Metadata) error {
	err := fileutil.CreateDirsFromAssets("tpl/ws/gradle", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tpl/ws/gradle", base)

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(base, "gitignore"), filepath.Join(base, ".gitignore"))

	if err != nil {
		return err
	}

	err = os.Chmod(filepath.Join(base, "gradlew"), 0774)

	if err != nil {
		return err
	}

	err = updateGradleProps(base, metadata)
	if err != nil {
		return err
	}

	err = updateGradleSettings(base)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleProps(base string, metadata *project.Metadata) error {
	err := fileutil.UpdateWithData(filepath.Join(base, "gradle.properties"), metadata)
	if err != nil {
		return err
	}
	err = fileutil.UpdateWithData(filepath.Join(base, "build.gradle"), metadata)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleSettings(base string) error {
	workspaceGradlePluginVersion := "10.0.6"

	resp, err := http.Get("https://raw.githubusercontent.com/lgdd/liferay-product-info/main/com.liferay.gradle.plugins.workspace")

	if err != nil {
		warningMessage := fmt.Sprintf("couldn't get latest version for com.liferay.gradle.plugins.workspace (default to %s)\n", workspaceGradlePluginVersion)
		printutil.Warning(warningMessage)
		err := fileutil.UpdateWithData(filepath.Join(base, "settings.gradle"), struct {
			WorkspaceGradlePluginVersion string
		}{WorkspaceGradlePluginVersion: workspaceGradlePluginVersion})
		if err != nil {
			return err
		}
		return nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	workspaceGradlePluginVersion = string(body)
	err = fileutil.UpdateWithData(filepath.Join(base, "settings.gradle"), struct {
		WorkspaceGradlePluginVersion string
	}{WorkspaceGradlePluginVersion: workspaceGradlePluginVersion})
	if err != nil {
		return err
	}
	return nil
}

func createMavenFiles(base string, metadata *project.Metadata) error {
	err := fileutil.CreateDirsFromAssets("tpl/ws/maven", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tpl/ws/maven", base)

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(base, "gitignore"), filepath.Join(base, ".gitignore"))

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(base, "mvn"), filepath.Join(base, ".mvn"))

	if err != nil {
		return err
	}

	err = os.Chmod(filepath.Join(base, "mvnw"), 0774)

	if err != nil {
		return err
	}

	err = updatePoms(base, metadata)
	if err != nil {
		return err
	}

	return nil
}

func updatePoms(base string, metadata *project.Metadata) error {
	poms := []string{
		filepath.Join(base, "pom.xml"),
		filepath.Join(base, "modules", "pom.xml"),
		filepath.Join(base, "themes", "pom.xml"),
		filepath.Join(base, "wars", "pom.xml"),
	}

	for _, pomPath := range poms {
		err := fileutil.UpdateWithData(pomPath, metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func createCommonEmptyDirs(base string) {
	configCommonDir := filepath.Join(base, "configs", "common")
	configDockerDir := filepath.Join(base, "configs", "docker")
	fileutil.CreateDirs(configCommonDir)
	fileutil.CreateDirs(configDockerDir)
	fileutil.CreateFiles([]string{filepath.Join(configCommonDir, ".gitkeep")})
	fileutil.CreateFiles([]string{filepath.Join(configDockerDir, ".gitkeep")})
}

func createGithubWorkflows(base string) error {
	javaVersion := "11"
	githubWorkflowsDir := filepath.Join(base, ".github", "workflows")
	fileutil.CreateDirs(filepath.Join(base, ".github", "workflows"))

	err := fileutil.CreateFilesFromAssets("tpl/gh", githubWorkflowsDir)

	if err != nil {
		return err
	}

	major, _, err := procutil.GetCurrentJavaVersion()

	if err == nil && (major == "8" || major == "11") {
		javaVersion = major
	}

	err = fileutil.UpdateWithData(filepath.Join(githubWorkflowsDir, "liferay-upgrade.yml"), struct {
		JavaVersion string
	}{JavaVersion: javaVersion})

	return err
}
