package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	cp "github.com/otiai10/copy"
)

// CreateClientExtension creates a client extension directory named name inside
// the current workspace by copying the files from the given sample project.
func CreateClientExtension(sample, name string) error {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		return err
	}

	clientExtensionSamplesPath := filepath.Join(conf.GetConfigPath(), conf.ClientExtensionSampleProjectName)
	clientExtensionsWorkspaceDir := filepath.Join(liferayWorkspace, "client-extensions")

	samplePath := filepath.Join(clientExtensionSamplesPath, conf.ClientExtensionSamplePrefix+sample)

	name = strcase.ToKebab(strings.ToLower(name))
	clientExtensionDir := filepath.Join(clientExtensionsWorkspaceDir, name)

	fileutil.CreateDirs(clientExtensionDir)

	if err := cp.Copy(samplePath, clientExtensionDir); err != nil {
		clientExtensionExtraSamplesPath := filepath.Join(conf.GetConfigPath(), conf.ClientExtensionExtraSampleProjectName)
		extraSamplePath := filepath.Join(clientExtensionExtraSamplesPath, conf.ClientExtensionExtraSamplePrefix+sample)
		if err := cp.Copy(extraSamplePath, clientExtensionDir); err != nil {
			return err
		}
	}

	printCreatedFiles(clientExtensionDir)

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		logger.PrintlnWarn("\nClient Extensions are not supported with Maven")
	}

	return nil
}
