package workspace

import (
	"os"
	"path/filepath"

	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/fileutil"
)

const (
	Gradle = "gradle"
	Maven  = "maven"
)

func Generate(base, build, version string) error {
	err := os.Mkdir(base, os.ModePerm)
	if err != nil {
		return err
	}

	switch build {
	case Gradle:
		if err := createGradleFiles(base, version); err != nil {
			return err
		}
	case Maven:
		if err := createMavenFiles(base, version); err != nil {
			return err
		}
	}

	return nil
}

func createGradleFiles(base string, version string) error {
	err := fileutil.CreateDirsFromAssets("tmpl/ws/gradle", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tmpl/ws/gradle", base)

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

	err = updateGradleProps(base, version)
	if err != nil {
		return err
	}

	return nil
}

func updateGradleProps(base, version string) error {
	metadata, err := project.NewMetadata(base, version)
	if err != nil {
		return err
	}

	err = fileutil.UpdateWithData(filepath.Join(base, "gradle.properties"), metadata)
	if err != nil {
		return err
	}
	return nil
}

func createMavenFiles(base, version string) error {
	err := fileutil.CreateDirsFromAssets("tmpl/ws/maven", base)

	if err != nil {
		return err
	}

	err = fileutil.CreateFilesFromAssets("tmpl/ws/maven", base)

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

	err = updatePoms(base, version)
	if err != nil {
		return err
	}

	return nil
}

func updatePoms(base, version string) error {
	project, err := project.NewMetadata(base, version)
	if err != nil {
		return err
	}

	poms := []string{
		filepath.Join(base, "pom.xml"),
		filepath.Join(base, "modules", "pom.xml"),
		filepath.Join(base, "themes", "pom.xml"),
		filepath.Join(base, "wars", "pom.xml"),
	}

	for _, pomPath := range poms {
		err = fileutil.UpdateWithData(pomPath, project)
		if err != nil {
			return err
		}
	}

	return nil
}
