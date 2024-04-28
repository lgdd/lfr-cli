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
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// SpringPortletData contains the data to be injected into the template files
type SpringPortletData struct {
	Package                 string
	Name                    string
	CamelCaseName           string
	WorkspaceName           string
	WorkspaceCamelCaseName  string
	WorkspacePackage        string
	PortletIDKey            string
	PortletIDValue          string
	MajorVersion            string
	DtdMajorVersion         string
	TemplateEngine          string
	WorkspaceProductEdition string
}

// Creates the structure for a Spring portlet module (i.e. PortletMVC4Spring)
func CreateModuleSpring(name, templateEngine string) {

	if templateEngine != "thymeleaf" && templateEngine != "jsp" {
		logger.Fatal("invalid template engine: use thymeleaf or jsp")
	}

	sep := string(os.PathSeparator)
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		logger.Fatal(err.Error())
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

	err = fileutil.CreateDirsFromAssets("tpl/spring", destPortletPath)

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = fileutil.CreateFilesFromAssets("tpl/spring", destPortletPath)

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore"))

	if err != nil {
		logger.Fatal(err.Error())
	}

	fileutil.CreateDirs(packagePath)

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		pomPath := filepath.Join(destPortletPath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		buildGradlePath := filepath.Join(destPortletPath, "build.gradle")
		err = os.Remove(buildGradlePath)

		if err != nil {
			logger.Fatal(err.Error())
		}

		pomParentPath := filepath.Join(destPortletPath, "../pom.xml")
		pomParent, err := os.Open(pomParentPath)
		if err != nil {
			logger.Fatal(err.Error())
		}
		defer pomParent.Close()

		byteValue, _ := io.ReadAll(pomParent)

		var pom fileutil.Pom
		err = xml.Unmarshal(byteValue, &pom)

		if err != nil {
			logger.Fatal(err.Error())
		}

		modules := append(pom.Modules.Module, name)
		pom.Modules.Module = modules
		pom.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
		pom.SchemaLocation = "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"

		finalPomBytes, _ := xml.MarshalIndent(pom, "", "  ")

		err = os.WriteFile(pomParentPath, []byte(fileutil.XMLHeader+string(finalPomBytes)), 0644)

		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.PrintWarn("modified ")
		fmt.Printf("%s\n", pomParentPath)
	}

	version, err := fileutil.GetLiferayWorkspaceProductVersion(liferayWorkspace)

	if err != nil {
		logger.Fatal(err.Error())
	}

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)

	if err != nil {
		logger.Fatal(err.Error())
	}

	portletIDKey := strcase.ToScreamingDelimited(name, '_', "", true)
	portletIDKey = strings.ToUpper(portletIDKey)
	portletIDValue := strings.ToLower(portletIDKey) + "_" + camelCaseName

	portletData := &SpringPortletData{
		Package:                 portletPackage,
		Name:                    name,
		CamelCaseName:           camelCaseName,
		PortletIDKey:            portletIDKey,
		PortletIDValue:          portletIDValue,
		WorkspaceName:           workspaceName,
		WorkspaceCamelCaseName:  strcase.ToCamel(workspaceName),
		WorkspacePackage:        workspacePackage,
		MajorVersion:            version,
		DtdMajorVersion:         strings.ReplaceAll(version, ".", "_"),
		TemplateEngine:          templateEngine,
		WorkspaceProductEdition: workspaceProductEdition,
	}

	updateFiles(portletData, destPortletPath, packagePath)

	err = updateMvcPortletWithData(destPortletPath, portletData)

	if err != nil {
		logger.Fatal(err.Error())
	}

	_ = filepath.Walk(destPortletPath,
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
}

func updateFiles(portletData *SpringPortletData, modulePath, packagePath string) {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")

	fileutil.CreateDirs(filepath.Join(packagePath, "controller"))

	err := os.Rename(filepath.Join(defaultSrcPath, "UserController.java"), filepath.Join(packagePath, "controller", "UserController.java"))

	if err != nil {
		logger.Fatal(err.Error())
	}

	fileutil.CreateDirs(filepath.Join(packagePath, "dto"))

	err = os.Rename(filepath.Join(defaultSrcPath, "User.java"), filepath.Join(packagePath, "dto", "User.java"))

	if err != nil {
		logger.Fatal(err.Error())
	}

	springPortletContextPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "spring-context", "portlet")

	err = os.Rename(filepath.Join(springPortletContextPath, "Spring.xml"), filepath.Join(springPortletContextPath, portletData.CamelCaseName+".xml"))

	if err != nil {
		logger.Fatal(err.Error())
	}

	viewsPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "views")
	viewsExt := ".jspx"

	if portletData.TemplateEngine == "thymeleaf" {
		viewsExt = ".html"
	}

	err = os.Rename(filepath.Join(viewsPath, "user.tpl"), filepath.Join(viewsPath, "user"+viewsExt))

	if err != nil {
		logger.Fatal(err.Error())
	}

	err = os.Rename(filepath.Join(viewsPath, "greeting.tpl"), filepath.Join(viewsPath, "greeting"+viewsExt))

	if err != nil {
		logger.Fatal(err.Error())
	}

}

func updateMvcPortletWithData(destPortletPath string, portletData *SpringPortletData) error {
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
