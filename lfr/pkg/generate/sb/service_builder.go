package sb

import (
	"encoding/xml"
	"fmt"
	"github.com/magiconair/properties"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/lgdd/liferay-cli/lfr/pkg/project"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

type ServiceBuilderData struct {
	Package                string
	Name                   string
	CamelCaseName          string
	WorkspaceName          string
	WorkspaceCamelCaseName string
	WorkspacePackage       string
	MajorVersion           string
	DtdMajorVersion        string
}

func Generate(name string) {
	sep := string(os.PathSeparator)
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	name = strcase.ToKebab(name)
	destModuleParentPath := filepath.Join(liferayWorkspace, "modules")
	destModulePath := filepath.Join(destModuleParentPath, name)
	destModuleAPIPath := filepath.Join(destModuleParentPath, name, name+"-api")
	destModuleServicePath := filepath.Join(destModuleParentPath, name, name+"-service")
	camelCaseName := strcase.ToCamel(name)
	workspaceSplit := strings.Split(liferayWorkspace, sep)
	workspaceName := workspaceSplit[len(workspaceSplit)-1]
	workspacePackage := strcase.ToDelimited(workspaceName, '.')

	err = fileutil.CreateDirsFromAssets("tpl/sb", destModulePath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = fileutil.CreateFilesFromAssets("tpl/sb", destModulePath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = renameFiles(destModulePath, destModuleAPIPath, destModuleServicePath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		pomPath := filepath.Join(destModulePath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(pomPath)

		pomPath = filepath.Join(destModuleAPIPath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(pomPath)

		pomPath = filepath.Join(destModuleServicePath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(pomPath)
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		buildGradlePath := filepath.Join(destModuleAPIPath, "build.gradle")
		err = os.Remove(buildGradlePath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(buildGradlePath)

		buildGradlePath = filepath.Join(destModuleServicePath, "build.gradle")
		err = os.Remove(buildGradlePath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(buildGradlePath)

		pomParentPath := filepath.Join(destModulePath, "../pom.xml")
		pomParent, err := os.Open(pomParentPath)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
		defer pomParent.Close()

		byteValue, _ := ioutil.ReadAll(pomParent)

		var pom project.Pom
		err = xml.Unmarshal(byteValue, &pom)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		modules := append(pom.Modules.Module, name)
		pom.Modules.Module = modules
		pom.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
		pom.SchemaLocation = "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"

		finalPomBytes, _ := xml.MarshalIndent(pom, "", "  ")

		err = ioutil.WriteFile(pomParentPath, []byte(project.XMLHeader+string(finalPomBytes)), 0644)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Warning("update ")
		fmt.Printf("%s\n", pomParentPath)
	}

	version, err := getLiferayMajorVersion()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	data := &ServiceBuilderData{
		Package:                strcase.ToDelimited(name, '.'),
		Name:                   name,
		CamelCaseName:          camelCaseName,
		WorkspaceName:          workspaceName,
		WorkspaceCamelCaseName: strcase.ToCamel(workspaceName),
		WorkspacePackage:       workspacePackage,
		MajorVersion:           version,
		DtdMajorVersion:        strings.ReplaceAll(version, ".", "_"),
	}

	err = updateModuleWithData(destModulePath, data)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

func getLiferayMajorVersion() (string, error) {
	workspacePath, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return "", err
	}

	if fileutil.IsMavenWorkspace(workspacePath) {
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
		targetVersion := pom.Properties.LiferayBomVersion
		semver := strings.Split(targetVersion, ".")
		version := strings.Join(append(semver[:len(semver)-1], "0"), ".")
		return version, nil
	}

	if fileutil.IsGradleWorkspace(workspacePath) {
		gradlePropsPath := filepath.Join(workspacePath, "gradle.properties")
		gradleProps := properties.MustLoadFile(gradlePropsPath, properties.UTF8)
		product := gradleProps.GetString("liferay.workspace.product", "portal-7.3-ga7")
		version := strings.Split(product, "-")[1]
		version += ".0"
		return version, nil
	}

	return "", nil
}

func renameFiles(destModulePath string, destModuleAPIPath string, destModuleServicePath string) error {
	err := os.Rename(filepath.Join(destModulePath, "sb-api"), destModuleAPIPath)

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(destModulePath, "sb-service"), destModuleServicePath)

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(destModuleAPIPath, "gitignore"), filepath.Join(destModuleAPIPath, ".gitignore"))

	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(destModuleServicePath, "gitignore"), filepath.Join(destModuleServicePath, ".gitignore"))

	if err != nil {
		return err
	}

	return err
}

func updateModuleWithData(destModulePath string, data *ServiceBuilderData) error {
	return filepath.Walk(destModulePath, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = fileutil.UpdateWithData(path, data)
		}

		if err != nil {
			return err
		}

		return nil
	})
}
