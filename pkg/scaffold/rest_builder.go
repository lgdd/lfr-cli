package scaffold

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

// RESTBuilderData contains the data to be injected into the template files
type RESTBuilderData struct {
	Package                string
	Name                   string
	CamelCaseName          string
	WorkspaceName          string
	WorkspaceCamelCaseName string
	WorkspacePackage       string
	MajorVersion           string
	DtdMajorVersion        string
	User                   string
}

// Creates the structure of a REST Builder module
func CreateModuleRESTBuilder(liferayWorkspace, name string) error {
	modulePackage, workspacePackage := resolvePackageName(name)

	destModuleParentPath := filepath.Join(liferayWorkspace, "modules")
	destModulePath := filepath.Join(destModuleParentPath, name)
	destModuleAPIPath := filepath.Join(destModuleParentPath, name, name+"-api")
	destModuleImplPath := filepath.Join(destModuleParentPath, name, name+"-impl")
	camelCaseName := strcase.ToCamel(name)
	workspaceName := workspaceBaseName(liferayWorkspace)

	if err := fileutil.CreateDirsFromAssets("tpl/rest_builder", destModulePath); err != nil {
		return err
	}

	if err := fileutil.CreateFilesFromAssets("tpl/rest_builder", destModulePath); err != nil {
		return err
	}

	if err := renameModuleRESTBuilderFiles(destModulePath, destModuleAPIPath, destModuleImplPath); err != nil {
		return err
	}

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		for _, path := range []string{destModulePath, destModuleAPIPath, destModuleImplPath} {
			if err := removeUnusedBuildFile(liferayWorkspace, path); err != nil {
				return err
			}
		}
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		for _, path := range []string{destModuleAPIPath, destModuleImplPath} {
			if err := removeUnusedBuildFile(liferayWorkspace, path); err != nil {
				return err
			}
		}

		pomParentPath := filepath.Join(destModulePath, "../pom.xml")
		if err := fileutil.AppendModuleToPom(pomParentPath, name); err != nil {
			return err
		}
		printModified(pomParentPath)
	}

	version, err := fileutil.GetLiferayWorkspaceProductVersion(liferayWorkspace)
	if err != nil {
		return err
	}

	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	data := &RESTBuilderData{
		Package:                modulePackage,
		Name:                   name,
		CamelCaseName:          camelCaseName,
		WorkspaceName:          workspaceName,
		WorkspaceCamelCaseName: strcase.ToCamel(workspaceName),
		WorkspacePackage:       workspacePackage,
		MajorVersion:           version,
		DtdMajorVersion:        strings.ReplaceAll(version, ".", "_"),
		User:                   currentUser.Username,
	}

	if err = updateFilesWithData(destModulePath, data); err != nil {
		return err
	}

	printCreatedFiles(destModulePath)
	return nil
}

func renameModuleRESTBuilderFiles(destModulePath string, destModuleAPIPath string, destModuleImplPath string) error {
	if err := os.Rename(filepath.Join(destModulePath, "api"), destModuleAPIPath); err != nil {
		return err
	}

	return os.Rename(filepath.Join(destModulePath, "impl"), destModuleImplPath)
}
