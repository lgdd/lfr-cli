package gen

import (
	"github.com/iancoleman/strcase"
	"github.com/lgdd/deba/util"
	"os"
	"path/filepath"
)

const (
	Gradle string = "gradle"
	Maven         = "maven"
)

func CreateWorkspace(base, build, version string) error {
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
		err := os.MkdirAll(filepath.Join(base, dir), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func createCommonFiles(base string) error {
	files := map[string]string{
		"/tpl/ws/gitignore":                           filepath.Join(base, ".gitignore"),
		"/tpl/ws/configs/dev/portal-ext.properties":   filepath.Join(base, "configs", "dev", "portal-ext.properties"),
		"/tpl/ws/configs/local/portal-ext.properties": filepath.Join(base, "configs", "local", "portal-ext.properties"),
		"/tpl/ws/configs/uat/portal-ext.properties":   filepath.Join(base, "configs", "uat", "portal-ext.properties"),
		"/tpl/ws/configs/prod/portal-ext.properties":  filepath.Join(base, "configs", "prod", "portal-ext.properties"),
		"/tpl/ws/configs/uat/es.config":               filepath.Join(base, "configs", "uat", "osgi", "configs", "com.liferay.portal.search.elasticsearch.configuration.ElasticsearchConfiguration.config"),
		"/tpl/ws/configs/prod/es.config":              filepath.Join(base, "configs", "prod", "osgi", "configs", "com.liferay.portal.search.elasticsearch.configuration.ElasticsearchConfiguration.config"),
	}

	for source, dest := range files {
		err := util.CopyPkgedFile(source, dest)
		if err != nil {
			return err
		}
	}

	emptyFiles := []string{
		filepath.Join(base, "configs", "common", ".touch"),
		filepath.Join(base, "configs", "docker", ".touch"),
	}

	for _, emptyFile := range emptyFiles {
		_, err := os.Create(emptyFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func createGradleFiles(base string, version string) error {
	err := os.MkdirAll(filepath.Join(base, filepath.Join("gradle", "wrapper")), os.ModePerm)
	if err != nil {
		return err
	}

	files := map[string]string{
		"/tpl/ws/gradle/gradle-wrapper.properties": filepath.Join(base, "gradle", "wrapper", "gradle-wrapper.properties"),
		"/tpl/ws/gradle/gradle-wrapper.jar":        filepath.Join(base, "gradle", "wrapper", "gradle-wrapper.jar"),
		"/tpl/ws/gradle/gradlew":                   filepath.Join(base, "gradlew"),
		"/tpl/ws/gradle/gradlew.bat":               filepath.Join(base, "gradlew.bat"),
		"/tpl/ws/gradle/settings.gradle":           filepath.Join(base, "settings.gradle"),
		"/tpl/ws/gradle/gradle.properties":         filepath.Join(base, "gradle.properties"),
	}

	emptyFiles := []string{
		filepath.Join(base, "modules", ".touch"),
		filepath.Join(base, "themes", ".touch"),
		filepath.Join(base, "wars", ".touch"),
		filepath.Join(base, "gradle-local.properties"),
		filepath.Join(base, "build.gradle"),
	}

	for source, dest := range files {
		err := util.CopyPkgedFile(source, dest)
		if err != nil {
			return err
		}
	}

	for _, emptyFile := range emptyFiles {
		_, err := os.Create(emptyFile)
		if err != nil {
			return err
		}
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
	project, err := util.GetProject(version)
	if err != nil {
		return err
	}

	err = util.UpdateFile(filepath.Join(base, "gradle.properties"), project)
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
		"/tpl/ws/maven/maven-wrapper.properties": filepath.Join(base, ".mvn", "wrapper", "maven-wrapper.properties"),
		"/tpl/ws/maven/maven-wrapper.jar":        filepath.Join(base, ".mvn", "wrapper", "maven-wrapper.jar"),
		"/tpl/ws/maven/mvnw":                     filepath.Join(base, "mvnw"),
		"/tpl/ws/maven/mvnw.cmd":                 filepath.Join(base, "mvnw.cmd"),
		"/tpl/ws/maven/parent-pom.xml":           filepath.Join(base, "pom.xml"),
		"/tpl/ws/maven/modules-pom.xml":          filepath.Join(base, "modules", "pom.xml"),
		"/tpl/ws/maven/themes-pom.xml":           filepath.Join(base, "themes", "pom.xml"),
		"/tpl/ws/maven/wars-pom.xml":             filepath.Join(base, "wars", "pom.xml"),
	}

	for source, dest := range files {
		err := util.CopyPkgedFile(source, dest)
		if err != nil {
			return err
		}
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
	project, err := util.GetProject(version)
	if err != nil {
		return err
	}
	project.GroupId = strcase.ToDelimited(base, '.')
	project.ArtifactId = base
	project.Name = strcase.ToCamel(base)

	poms := []string{
		filepath.Join(base, "pom.xml"),
		filepath.Join(base, "modules", "pom.xml"),
		filepath.Join(base, "themes", "pom.xml"),
		filepath.Join(base, "wars", "pom.xml"),
	}

	for _, pomPath := range poms {
		err = util.UpdateFile(pomPath, project)
		if err != nil {
			return err
		}
	}

	return nil
}
