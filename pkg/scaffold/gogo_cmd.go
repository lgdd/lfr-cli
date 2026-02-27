package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

// CmdData contains the data to be injected into the template files
type CmdData struct {
	Package                 string
	Name                    string
	CamelCaseName           string
	WorkspaceName           string
	WorkspaceCamelCaseName  string
	WorkspacePackage        string
	WorkspaceProductEdition string
}

// Creates the structure for a Gogo shell command module
func CreateModuleGogoCommand(name string) error {
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

	if err = fileutil.CreateDirsFromAssets("tpl/gogo_cmd", destPortletPath); err != nil {
		return err
	}

	if err = fileutil.CreateFilesFromAssets("tpl/gogo_cmd", destPortletPath); err != nil {
		return err
	}

	if err = os.Rename(filepath.Join(destPortletPath, "gitignore"), filepath.Join(destPortletPath, ".gitignore")); err != nil {
		return err
	}

	fileutil.CreateDirs(packagePath)
	fileutil.CreateDirs(filepath.Join(destPortletPath, "src", "main", "resources", "META-INF", "resources"))
	fileutil.CreateFiles([]string{filepath.Join(destPortletPath, "src", "main", "resources", ".gitkeep")})

	if err = updateModuleGogoCommandJavaFiles(camelCaseName, destPortletPath, packagePath); err != nil {
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

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)
	if err != nil {
		return err
	}

	data := &CmdData{
		Package:                 portletPackage,
		Name:                    name,
		CamelCaseName:           camelCaseName,
		WorkspaceName:           workspaceName,
		WorkspaceCamelCaseName:  strcase.ToCamel(workspaceName),
		WorkspacePackage:        workspacePackage,
		WorkspaceProductEdition: workspaceProductEdition,
	}

	if err = updateFilesWithData(destPortletPath, data); err != nil {
		return err
	}

	printCreatedFiles(destPortletPath)
	return nil
}

func updateModuleGogoCommandJavaFiles(camelCaseName, modulePath, packagePath string) error {
	defaultSrcPath := filepath.Join(modulePath, "src", "main", "java")
	return os.Rename(filepath.Join(defaultSrcPath, "Cmd.java"), filepath.Join(packagePath, camelCaseName+".java"))
}
