package spring

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/magiconair/properties"

	"github.com/lgdd/liferay-cli/lfr/pkg/project"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

type SpringPortletData struct {
	Package                string
	Name                   string
	CamelCaseName          string
	WorkspaceName          string
	WorkspaceCamelCaseName string
	WorkspacePackage       string
	PortletIDKey           string
	PortletIDValue         string
	MajorVersion           string
	DtdMajorVersion        string
	TemplateEngine         string
}

func Generate(name, templateEngine string) {

	if templateEngine != "thymeleaf" && templateEngine != "jsp" {
		printutil.Danger("invalid template engine\n")
		fmt.Println("Please use thymeleaf or jsp")
		os.Exit(1)
	}

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

	err = fileutil.CreateDirsFromAssets("tpl/spring", destPortletPath)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = fileutil.CreateFilesFromAssets("tpl/spring", destPortletPath)

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

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		pomPath := filepath.Join(destPortletPath, "pom.xml")
		err = os.Remove(pomPath)

		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Danger("remove ")
		fmt.Println(pomPath)
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
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

	version, err := getLiferayMajorVersion(liferayWorkspace)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	portletIDKey := strcase.ToScreamingDelimited(name, '_', 0, true)
	portletIDKey = strings.ToUpper(portletIDKey)
	portletIDValue := strings.ToLower(portletIDKey) + "_" + camelCaseName

	portletData := &SpringPortletData{
		Package:                strcase.ToDelimited(name, '.'),
		Name:                   name,
		CamelCaseName:          camelCaseName,
		PortletIDKey:           portletIDKey,
		PortletIDValue:         portletIDValue,
		WorkspaceName:          workspaceName,
		WorkspaceCamelCaseName: strcase.ToCamel(workspaceName),
		WorkspacePackage:       workspacePackage,
		MajorVersion:           version,
		DtdMajorVersion:        strings.ReplaceAll(version, ".", "_"),
		TemplateEngine:         templateEngine,
	}

	updateFiles(portletData, destPortletPath, packagePath)

	err = updateMvcPortletWithData(destPortletPath, portletData)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

func updateFiles(portletData *SpringPortletData, modulePath, packagePath string) {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")

	fileutil.CreateDirs(filepath.Join(packagePath, "controller"))

	err := os.Rename(filepath.Join(defaultSrcPath, "UserController.java"), filepath.Join(packagePath, "controller", "UserController.java"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fileutil.CreateDirs(filepath.Join(packagePath, "dto"))

	err = os.Rename(filepath.Join(defaultSrcPath, "User.java"), filepath.Join(packagePath, "dto", "User.java"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	springPortletContextPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "spring-context", "portlet")

	err = os.Rename(filepath.Join(springPortletContextPath, "Spring.xml"), filepath.Join(springPortletContextPath, portletData.CamelCaseName+".xml"))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	viewsPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "views")
	viewsExt := ".jspx"

	if portletData.TemplateEngine == "thymeleaf" {
		viewsExt = ".html"
	}

	err = os.Rename(filepath.Join(viewsPath, "user.tpl"), filepath.Join(viewsPath, "user"+viewsExt))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = os.Rename(filepath.Join(viewsPath, "greeting.tpl"), filepath.Join(viewsPath, "greeting"+viewsExt))

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
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

func getLiferayMajorVersion(workspacePath string) (string, error) {
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