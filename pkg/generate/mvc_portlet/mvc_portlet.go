package mvc_portlet

import (
	"encoding/xml"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type MvcPortletData struct {
	Package                string
	Name                   string
	CamelCaseName          string
	WorkspaceName          string
	WorkspaceCamelCaseName string
	WorkspacePackage       string
	PortletIdKey           string
	PortletIdValue         string
}

func Generate(name string) {
	sep := string(os.PathSeparator)
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	name = strcase.ToKebab(name)
	destPortletParentPath := filepath.Join(liferayWorkspace, "modules")
	destPortletPath := filepath.Join(destPortletParentPath, name)
	packagePath := strings.ReplaceAll(name, "-", string(os.PathSeparator))
	packagePath = filepath.Join(destPortletPath, "src", "main", "java", packagePath)
	camelCaseName := strcase.ToCamel(name)
	workspaceSplit := strings.Split(liferayWorkspace, sep)
	workspaceName := workspaceSplit[len(workspaceSplit)-1]
	workspacePackage := strcase.ToDelimited(workspaceName, '.')

	err = fileutil.CreateDirsFromAssets("tpl/mvc-portlet", destPortletPath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = fileutil.CreateFilesFromAssets("tpl/mvc-portlet", destPortletPath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.CreateDirs(packagePath)

	updateJavaFiles(camelCaseName, destPortletPath, packagePath)

	if fileutil.IsGradleWorkspace() {
		pomPath := filepath.Join(destPortletPath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(pomPath)
	}

	if fileutil.IsMavenWorkspace() {
		buildGradlePath := filepath.Join(destPortletPath, "build.gradle")
		err = os.Remove(buildGradlePath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(buildGradlePath)

		pomParentPath := filepath.Join(destPortletPath, "../pom.xml")
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

	portletIdKey := strcase.ToScreamingDelimited(name, '_', 0, true)
	portletIdKey = strings.ToUpper(portletIdKey)
	portletIdValue := strings.ToLower(portletIdKey) + "_" + camelCaseName

	portletData := &MvcPortletData{
		Package:                strcase.ToDelimited(name, '.'),
		Name:                   name,
		CamelCaseName:          camelCaseName,
		PortletIdKey:           portletIdKey,
		PortletIdValue:         portletIdValue,
		WorkspaceName:          workspaceName,
		WorkspaceCamelCaseName: strcase.ToCamel(workspaceName),
		WorkspacePackage:       workspacePackage,
	}

	err = updateMvcPortletWithData(destPortletPath, portletData)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

func updateJavaFiles(camelCaseName, modulePath, packagePath string) {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")
	err := os.Rename(filepath.Join(defaultSrcPath, "Portlet.java"), filepath.Join(packagePath, camelCaseName+".java"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.CreateDirs(filepath.Join(packagePath, "constants"))

	err = os.Rename(filepath.Join(defaultSrcPath, "PortletKeys.java"), filepath.Join(packagePath, "constants", camelCaseName+"Keys.java"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

}

func updateMvcPortletWithData(destPortletPath string, portletData *MvcPortletData) error {
	return filepath.Walk(destPortletPath, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = fileutil.UpdateWithData(path, portletData)
		}

		if err != nil {
			return err
		}

		return nil
	})
}
