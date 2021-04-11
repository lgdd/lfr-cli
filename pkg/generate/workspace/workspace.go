package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
)

const (
	Gradle string = "gradle"
	Maven         = "maven"
)

func Generate(base, build, version string) error {
	err := os.Mkdir(base, os.ModePerm)
	if err != nil {
		return err
	}

	err = createCommonDirs(base)
	if err != nil {
		return err
	}

	err = createCommonFiles(base)
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

func createCommonDirs(base string) error {
	var wg sync.WaitGroup
	dirs := []string{
		filepath.Join("gradle", "wrapper"),
		filepath.Join("configs", "common"),
		filepath.Join("configs", "dev"),
		filepath.Join("configs", "docker"),
		filepath.Join("configs", "local"),
		filepath.Join("configs", "prod", "osgi", "configs"),
		filepath.Join("configs", "uat", "osgi", "configs"),
		"modules",
		"themes",
		"wars"}

	for _, dir := range dirs {
		wg.Add(1)
		go createDirs(filepath.Join(base, dir), &wg)
	}

	wg.Wait()
	return nil
}

func createCommonFiles(base string) error {
	esConfigFilename := strings.Join([]string{
		"com",
		"liferay",
		"portal",
		"search",
		"elasticsearch",
		"configuration",
		"ElasticsearchConfiguration",
		"config",
	}, ".")
	files := map[string]string{
		"tmpl/ws/gitignore":                           filepath.Join(base, ".gitignore"),
		"tmpl/ws/platform.bndrun":                     filepath.Join(base, "platform.bndrun"),
		"tmpl/ws/configs/dev/portal-ext.properties":   filepath.Join(base, "configs", "dev", "portal-ext.properties"),
		"tmpl/ws/configs/local/portal-ext.properties": filepath.Join(base, "configs", "local", "portal-ext.properties"),
		"tmpl/ws/configs/uat/portal-ext.properties":   filepath.Join(base, "configs", "uat", "portal-ext.properties"),
		"tmpl/ws/configs/prod/portal-ext.properties":  filepath.Join(base, "configs", "prod", "portal-ext.properties"),
		"tmpl/ws/configs/uat/es.config":               filepath.Join(base, "configs", "uat", "osgi", "configs", esConfigFilename),
		"tmpl/ws/configs/prod/es.config":              filepath.Join(base, "configs", "prod", "osgi", "configs", esConfigFilename),
	}

	var wg sync.WaitGroup
	for source, dest := range files {
		wg.Add(1)
		go fileutil.CopyFromAssets(source, dest, &wg)
	}
	wg.Wait()

	emptyFiles := []string{
		filepath.Join(base, "configs", "common", ".touch"),
		filepath.Join(base, "configs", "docker", ".touch"),
	}

	createFiles(emptyFiles)

	return nil
}

func createGradleFiles(base string, version string) error {
	err := os.MkdirAll(filepath.Join(base, filepath.Join("gradle", "wrapper")), os.ModePerm)
	if err != nil {
		return err
	}

	files := map[string]string{
		"tmpl/ws/gradle/gradle-wrapper.properties": filepath.Join(base, "gradle", "wrapper", "gradle-wrapper.properties"),
		"tmpl/ws/gradle/gradle-wrapper.jar":        filepath.Join(base, "gradle", "wrapper", "gradle-wrapper.jar"),
		"tmpl/ws/gradle/gradlew":                   filepath.Join(base, "gradlew"),
		"tmpl/ws/gradle/gradlew.bat":               filepath.Join(base, "gradlew.bat"),
		"tmpl/ws/gradle/settings.gradle":           filepath.Join(base, "settings.gradle"),
		"tmpl/ws/gradle/gradle.properties":         filepath.Join(base, "gradle.properties"),
	}

	emptyFiles := []string{
		filepath.Join(base, "modules", ".touch"),
		filepath.Join(base, "themes", ".touch"),
		filepath.Join(base, "wars", ".touch"),
		filepath.Join(base, "gradle-local.properties"),
		filepath.Join(base, "build.gradle"),
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	for source, dest := range files {
		go fileutil.CopyFromAssets(source, dest, &wg)
	}
	wg.Wait()

	createFiles(emptyFiles)

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
	err := os.MkdirAll(filepath.Join(base, filepath.Join(".mvn", "wrapper")), os.ModePerm)
	if err != nil {
		return err
	}

	files := map[string]string{
		"tmpl/ws/maven/maven-wrapper.properties": filepath.Join(base, ".mvn", "wrapper", "maven-wrapper.properties"),
		"tmpl/ws/maven/maven-wrapper.jar":        filepath.Join(base, ".mvn", "wrapper", "maven-wrapper.jar"),
		"tmpl/ws/maven/mvnw":                     filepath.Join(base, "mvnw"),
		"tmpl/ws/maven/mvnw.cmd":                 filepath.Join(base, "mvnw.cmd"),
		"tmpl/ws/maven/pom.xml":                  filepath.Join(base, "pom.xml"),
		"tmpl/ws/maven/modules/pom.xml":          filepath.Join(base, "modules", "pom.xml"),
		"tmpl/ws/maven/themes/pom.xml":           filepath.Join(base, "themes", "pom.xml"),
		"tmpl/ws/maven/wars/pom.xml":             filepath.Join(base, "wars", "pom.xml"),
	}

	var wg sync.WaitGroup
	for source, dest := range files {
		wg.Add(1)
		go fileutil.CopyFromAssets(source, dest, &wg)
	}
	wg.Wait()

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

func createDirs(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

func createFiles(paths []string) {
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go createFile(path, &wg)
	}
	wg.Wait()
}

func createFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := os.Create(path)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}
