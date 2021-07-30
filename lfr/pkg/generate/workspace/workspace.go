package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lgdd/liferay-cli/lfr/pkg/project"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

const (
	Gradle = "gradle"
	Maven  = "maven"
)

func Generate(base, build, version string) error {
	metadata, err := project.NewMetadata(base, version)
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
