package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// removeUnusedBuildFile removes pom.xml for Gradle workspaces or build.gradle
// for Maven workspaces from the given module directory.
func removeUnusedBuildFile(workspacePath, modulePath string) error {
	if fileutil.IsGradleWorkspace(workspacePath) {
		return os.Remove(filepath.Join(modulePath, "pom.xml"))
	}
	if fileutil.IsMavenWorkspace(workspacePath) {
		return os.Remove(filepath.Join(modulePath, "build.gradle"))
	}
	return nil
}

// resolvePackageName returns the module package and workspace group ID.
// If the module package is still the default "org.acme" but the workspace has
// a real group ID, the workspace group ID is used as a prefix for the module.
func resolvePackageName(name string) (pkg, workspacePkg string) {
	pkg = metadata.PackageName
	workspacePkg, _ = metadata.GetGroupId()
	if pkg == "org.acme" && workspacePkg != "org.acme" {
		pkg = strings.Join([]string{workspacePkg, strcase.ToDelimited(name, '.')}, ".")
	}
	return
}

// workspaceBaseName returns the last path component of a workspace path,
// used as the project name in template data.
func workspaceBaseName(workspacePath string) string {
	parts := strings.Split(workspacePath, string(os.PathSeparator))
	return parts[len(parts)-1]
}

// updateFilesWithData applies template data to all files under destPath.
func updateFilesWithData(destPath string, data any) error {
	return filepath.Walk(destPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = fileutil.UpdateWithData(path, data)
		}
		return err
	})
}

// printModified logs the given path as "modified".
func printModified(path string) {
	logger.PrintWarn("modified ")
	logger.Printf("%s\n", path)
}

// printCreatedFiles walks destPath and logs each file as "created".
func printCreatedFiles(destPath string) {
	_ = filepath.Walk(destPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			logger.PrintSuccess("created ")
			logger.Printf("%s\n", path)
		}
		return nil
	})
}
