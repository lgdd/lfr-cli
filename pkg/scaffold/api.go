package scaffold

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
)

// ModuleAPIData contains the data to be injected into the template files
type ModuleAPIData struct {
	Package                 string
	Name                    string
	CamelCaseName           string
	WorkspaceName           string
	WorkspaceCamelCaseName  string
	WorkspacePackage        string
	WorkspaceProductEdition string
}

// Creates the structure for an API module
func CreateModuleAPI(name string) {
	sep := string(os.PathSeparator)
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	portletPackage := metadata.PackageName
	workspacePackage, _ := metadata.GetGroupId()

	if portletPackage == "org.acme" && workspacePackage != "org.acme" {
		portletPackage = strings.Join([]string{workspacePackage, strcase.ToDelimited(name, '.')}, ".")
	}

	name = strcase.ToKebab(name)
	destPortletParentPath := filepath.Join(liferayWorkspace, "modules")
	destPortletPath := filepath.Join(destPortletParentPath, name)
	packagePath := strings.ReplaceAll(portletPackage, ".", string(os.PathSeparator))
	packagePath = filepath.Join(destPortletPath, "src", "main", "java", packagePath)
	camelCaseName := strcase.ToCamel(name)
	workspaceSplit := strings.Split(liferayWorkspace, sep)
	workspaceName := workspaceSplit[len(workspaceSplit)-1]

	err = fileutil.CreateDirsFromAssets("tpl/api", destPortletPath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = fileutil.CreateFilesFromAssets("tpl/api", destPortletPath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	if !strings.HasSuffix(packagePath, "api") {
		packagePath = filepath.Join(packagePath, "api")
	}

	fileutil.CreateDirs(packagePath)
	fileutil.CreateDirs(filepath.Join(destPortletPath, "src", "main", "resources", "META-INF", "resources"))
	fileutil.CreateFiles([]string{filepath.Join(destPortletPath, "src", "main", "resources", ".gitkeep")})

	updateModuleAPIJavaFiles(camelCaseName, destPortletPath, packagePath)

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		pomPath := filepath.Join(destPortletPath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		buildGradlePath := filepath.Join(destPortletPath, "build.gradle")
		err = os.Remove(buildGradlePath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		pomParentPath := filepath.Join(destPortletPath, "../pom.xml")
		pomParent, err := os.Open(pomParentPath)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
		defer pomParent.Close()

		byteValue, _ := io.ReadAll(pomParent)

		var pom fileutil.Pom
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

		err = os.WriteFile(pomParentPath, []byte(fileutil.XMLHeader+string(finalPomBytes)), 0644)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Warning("update ")
		fmt.Printf("%s\n", pomParentPath)
	}

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	data := &ModuleAPIData{
		Package:                 portletPackage,
		Name:                    name,
		CamelCaseName:           camelCaseName,
		WorkspaceName:           workspaceName,
		WorkspaceCamelCaseName:  strcase.ToCamel(workspaceName),
		WorkspacePackage:        workspacePackage,
		WorkspaceProductEdition: workspaceProductEdition,
	}

	err = updateModuleAPIWithData(destPortletPath, data)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	_ = filepath.Walk(destPortletPath,
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

}

func updateModuleAPIJavaFiles(camelCaseName, modulePath, packagePath string) {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")
	err := os.Rename(filepath.Join(defaultSrcPath, "Api.java"), filepath.Join(packagePath, camelCaseName+".java"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

func updateModuleAPIWithData(destPortletPath string, data *ModuleAPIData) error {
	return filepath.Walk(destPortletPath, func(path string, info fs.FileInfo, err error) error {

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
