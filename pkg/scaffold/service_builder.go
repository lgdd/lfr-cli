package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
)

// ServiceBuilderData contains the data to be injected into the template files
type ServiceBuilderData struct {
	Package                 string
	Name                    string
	CamelCaseName           string
	WorkspaceName           string
	WorkspaceCamelCaseName  string
	WorkspacePackage        string
	MajorVersion            string
	DtdMajorVersion         string
	WorkspaceProductEdition string
}

// Creates the structure for a Service Builder module
func CreateModuleServiceBuilder(liferayWorkspace, name string) error {
	modulePackage, workspacePackage := resolvePackageName(name)

	destModuleParentPath := filepath.Join(liferayWorkspace, "modules")
	destModulePath := filepath.Join(destModuleParentPath, name)
	destModuleAPIPath := filepath.Join(destModuleParentPath, name, name+"-api")
	destModuleServicePath := filepath.Join(destModuleParentPath, name, name+"-service")
	camelCaseName := strcase.ToCamel(name)
	workspaceName := workspaceBaseName(liferayWorkspace)

	if err := fileutil.CreateDirsFromAssets("tpl/service_builder", destModulePath); err != nil {
		return err
	}

	if err := fileutil.CreateFilesFromAssets("tpl/service_builder", destModulePath); err != nil {
		return err
	}

	if err := renameModuleServiceBuilderFiles(destModulePath, destModuleAPIPath, destModuleServicePath); err != nil {
		return err
	}

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		for _, path := range []string{destModulePath, destModuleAPIPath, destModuleServicePath} {
			if err := removeUnusedBuildFile(liferayWorkspace, path); err != nil {
				return err
			}
		}
	}

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		for _, path := range []string{destModuleAPIPath, destModuleServicePath} {
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

	productVersion, err := fileutil.GetLiferayWorkspaceProductVersion(liferayWorkspace)
	if err != nil {
		return err
	}

	var majorVersionBuilder strings.Builder
	majorVersionBuilder.WriteString(productVersion)
	majorVersionBuilder.WriteString(".0")

	workspaceProductEdition, err := fileutil.GetLiferayWorkspaceProductEdition(liferayWorkspace)
	if err != nil {
		return err
	}

	data := &ServiceBuilderData{
		Package:                 modulePackage,
		Name:                    name,
		CamelCaseName:           camelCaseName,
		WorkspaceName:           workspaceName,
		WorkspaceCamelCaseName:  strcase.ToCamel(workspaceName),
		WorkspacePackage:        workspacePackage,
		MajorVersion:            majorVersionBuilder.String(),
		DtdMajorVersion:         strings.ReplaceAll(majorVersionBuilder.String(), ".", "_"),
		WorkspaceProductEdition: workspaceProductEdition,
	}

	if err = updateFilesWithData(destModulePath, data); err != nil {
		return err
	}

	printCreatedFiles(destModulePath)
	return nil
}

func renameModuleServiceBuilderFiles(destModulePath string, destModuleAPIPath string, destModuleServicePath string) error {
	if err := os.Rename(filepath.Join(destModulePath, "api"), destModuleAPIPath); err != nil {
		return err
	}

	if err := os.Rename(filepath.Join(destModulePath, "service"), destModuleServicePath); err != nil {
		return err
	}

	if err := os.Rename(filepath.Join(destModuleAPIPath, "gitignore"), filepath.Join(destModuleAPIPath, ".gitignore")); err != nil {
		return err
	}

	return os.Rename(filepath.Join(destModuleServicePath, "gitignore"), filepath.Join(destModuleServicePath, ".gitignore"))
}
