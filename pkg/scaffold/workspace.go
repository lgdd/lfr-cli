package scaffold

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/lgdd/lfr-cli/pkg/util/procutil"
)

// Build options (i.e. Maven or Gradle)
const (
	Gradle = "gradle"
	Maven  = "maven"
)

// Create the structure of a Liferay workspace
func CreateWorkspace(base, build, version, edition string) error {
	workspaceData, err := metadata.NewWorkspaceData(base, version, edition)
	if err != nil {
		return err
	}

	if build == Maven {
		err = os.Mkdir(base, os.ModePerm)
		if err != nil {
			return err
		}
		if err := createMavenFiles(base, workspaceData); err != nil {
			return err
		}
		createCommonEmptyDirs(base)
	} else if build == Gradle {
		err = os.Mkdir(base, os.ModePerm)
		if err != nil {
			return err
		}
		if err := createGradleFiles(base, version, edition, workspaceData); err != nil {
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
				logger.PrintSuccess("created ")
				fmt.Printf("%s\n", path)
			}
			return nil
		})

	return nil
}

func createGradleFiles(base, version, edition string, workspaceData *metadata.WorkspaceData) error {
	err := fileutil.CreateDirsFromAssets("tpl/workspace/gradle", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tpl/workspace/gradle", base)

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

	err = updateGradleWrapper(base, workspaceData)
	if err != nil {
		return err
	}

	err = updateGradleProps(base, workspaceData)
	if err != nil {
		return err
	}

	err = updateGradleSettings(base, version, edition)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleProps(base string, workspaceData *metadata.WorkspaceData) error {
	err := fileutil.UpdateWithData(filepath.Join(base, "gradle.properties"), workspaceData)
	if err != nil {
		return err
	}
	err = fileutil.UpdateWithData(filepath.Join(base, "build.gradle"), workspaceData)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleWrapper(base string, workspaceData *metadata.WorkspaceData) error {
	err := fileutil.UpdateWithData(filepath.Join(base, "gradle", "wrapper", "gradle-wrapper.properties"), workspaceData)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleSettings(base, version, edition string) error {
	workspaceGradlePluginVersion := "10.1.9"

	if version != "7.4" || edition == metadata.Portal {
		err := fileutil.UpdateWithData(filepath.Join(base, "settings.gradle"), struct {
			WorkspaceGradlePluginVersion string
		}{WorkspaceGradlePluginVersion: workspaceGradlePluginVersion})
		if err != nil {
			return err
		}
		return nil
	}

	resp, err := http.Get("https://raw.githubusercontent.com/lgdd/liferay-product-info/main/com.liferay.gradle.plugins.workspace")

	if err != nil {
		logger.PrintfWarn("couldn't get latest version for com.liferay.gradle.plugins.workspace (default to %s)\n", workspaceGradlePluginVersion)
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

func createMavenFiles(base string, workspaceData *metadata.WorkspaceData) error {
	err := fileutil.CreateDirsFromAssets("tpl/workspace/maven", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tpl/workspace/maven", base)

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

	err = updatePoms(base, workspaceData)
	if err != nil {
		return err
	}

	return nil
}

func updatePoms(base string, workspaceData *metadata.WorkspaceData) error {
	poms := []string{
		filepath.Join(base, "pom.xml"),
		filepath.Join(base, "modules", "pom.xml"),
		filepath.Join(base, "themes", "pom.xml"),
		filepath.Join(base, "wars", "pom.xml"),
	}

	for _, pomPath := range poms {
		err := fileutil.UpdateWithData(pomPath, workspaceData)
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

	err := fileutil.CreateFilesFromAssets("tpl/github", githubWorkflowsDir)

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
