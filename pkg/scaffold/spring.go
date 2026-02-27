package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
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
func CreateModuleSpring(name, templateEngine string) error {
	if templateEngine != "thymeleaf" && templateEngine != "jsp" {
		return fmt.Errorf("invalid template engine: use thymeleaf or jsp")
	}

	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return err
	}

	portletPackage, workspacePackage := resolvePackageName(name)
	name = strcase.ToKebab(name)
	destPortletParentPath := filepath.Join(liferayWorkspace, "modules")
	destPortletPath := filepath.Join(destPortletParentPath, name)
	packagePath := strings.ReplaceAll(portletPackage, ".", string(os.PathSeparator))
	packagePath = filepath.Join(destPortletPath, "src", "main", "java", packagePath)
	camelCaseName := strcase.ToCamel(name)
	workspaceName := workspaceBaseName(liferayWorkspace)

	if err = fileutil.CreateDirsFromAssets("tpl/spring", destPortletPath); err != nil {
		return err
	}

	if err = fileutil.CreateFilesFromAssets("tpl/spring", destPortletPath); err != nil {
		return err
	}

	if err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore")); err != nil {
		return err
	}

	fileutil.CreateDirs(packagePath)

	if err = removeUnusedBuildFile(liferayWorkspace, destPortletPath); err != nil {
		return err
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		pomParentPath := filepath.Join(destPortletPath, "../pom.xml")
		if err = fileutil.AppendModuleToPom(pomParentPath, name); err != nil {
			return err
		}
		printModified(pomParentPath)
	}

	version, err := fileutil.GetLiferayWorkspaceProductVersion(liferayWorkspace)
	if err != nil {
		return err
	}

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)
	if err != nil {
		return err
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

	if err = updateSpringFiles(portletData, destPortletPath, packagePath); err != nil {
		return err
	}

	if err = updateFilesWithData(destPortletPath, portletData); err != nil {
		return err
	}

	printCreatedFiles(destPortletPath)
	return nil
}

func updateSpringFiles(portletData *SpringPortletData, modulePath, packagePath string) error {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")

	fileutil.CreateDirs(filepath.Join(packagePath, "controller"))

	if err := os.Rename(filepath.Join(defaultSrcPath, "UserController.java"), filepath.Join(packagePath, "controller", "UserController.java")); err != nil {
		return err
	}

	fileutil.CreateDirs(filepath.Join(packagePath, "dto"))

	if err := os.Rename(filepath.Join(defaultSrcPath, "User.java"), filepath.Join(packagePath, "dto", "User.java")); err != nil {
		return err
	}

	springPortletContextPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "spring-context", "portlet")

	if err := os.Rename(filepath.Join(springPortletContextPath, "Spring.xml"), filepath.Join(springPortletContextPath, portletData.CamelCaseName+".xml")); err != nil {
		return err
	}

	viewsPath := filepath.Join(modulePath, "src", "main", "webapp", "WEB-INF", "views")
	viewsExt := ".jspx"

	if portletData.TemplateEngine == "thymeleaf" {
		viewsExt = ".html"
	}

	if err := os.Rename(filepath.Join(viewsPath, "user.tpl"), filepath.Join(viewsPath, "user"+viewsExt)); err != nil {
		return err
	}

	return os.Rename(filepath.Join(viewsPath, "greeting.tpl"), filepath.Join(viewsPath, "greeting"+viewsExt))
}
