package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

// PortletData contains the data to be injected into the template files
type PortletData struct {
	Package                 string
	Name                    string
	CamelCaseName           string
	WorkspaceName           string
	WorkspaceCamelCaseName  string
	WorkspacePackage        string
	PortletIDKey            string
	PortletIDValue          string
	WorkspaceProductEdition string
}

// Creates the structure for a portlet module
func CreateModuleMVC(name string) error {
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

	if err = fileutil.CreateDirsFromAssets("tpl/mvc_portlet", destPortletPath); err != nil {
		return err
	}

	if err = fileutil.CreateFilesFromAssets("tpl/mvc_portlet", destPortletPath); err != nil {
		return err
	}

	if err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore")); err != nil {
		return err
	}

	fileutil.CreateDirs(packagePath)

	if err = updateModuleMVCJavaFiles(camelCaseName, destPortletPath, packagePath); err != nil {
		return err
	}

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

	portletIDKey := strcase.ToScreamingDelimited(name, '_', "", true)
	portletIDKey = strings.ToUpper(portletIDKey)
	portletIDValue := strings.ToLower(portletIDKey) + "_" + camelCaseName

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)
	if err != nil {
		return err
	}

	portletData := &PortletData{
		Package:                 portletPackage,
		Name:                    name,
		CamelCaseName:           camelCaseName,
		PortletIDKey:            portletIDKey,
		PortletIDValue:          portletIDValue,
		WorkspaceName:           workspaceName,
		WorkspaceCamelCaseName:  strcase.ToCamel(workspaceName),
		WorkspacePackage:        workspacePackage,
		WorkspaceProductEdition: workspaceProductEdition,
	}

	if err = updateFilesWithData(destPortletPath, portletData); err != nil {
		return err
	}

	printCreatedFiles(destPortletPath)
	return nil
}

func updateModuleMVCJavaFiles(camelCaseName, modulePath, packagePath string) error {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")
	if err := os.Rename(filepath.Join(defaultSrcPath, "Portlet.java"), filepath.Join(packagePath, camelCaseName+".java")); err != nil {
		return err
	}

	fileutil.CreateDirs(filepath.Join(packagePath, "constants"))

	return os.Rename(filepath.Join(defaultSrcPath, "PortletKeys.java"), filepath.Join(packagePath, "constants", camelCaseName+"Keys.java"))
}
