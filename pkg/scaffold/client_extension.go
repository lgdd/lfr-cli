package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	cp "github.com/otiai10/copy"
)

func CreateClientExtension(sample, name string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		logger.Fatal(err.Error())
	}

	clientExtensionSamplesPath := filepath.Join(config.GetConfigPath(), config.ClientExtensionSampleProjectName)
	clientExtensionsWorkspaceDir := filepath.Join(liferayWorkspace, "client-extensions")

	samplePath := filepath.Join(clientExtensionSamplesPath, config.ClientExtensionSamplePrefix+sample)

	name = strcase.ToKebab(strings.ToLower(name))
	clientExtensionDir := filepath.Join(clientExtensionsWorkspaceDir, name)

	fileutil.CreateDirs(clientExtensionDir)

	if err := cp.Copy(samplePath, clientExtensionDir); err != nil {
		logger.Fatal(err.Error())
	}

	_ = filepath.Walk(clientExtensionDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				logger.PrintSuccess("created ")
				logger.Printf("%s\n", path)
			}
			return nil
		})

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		logger.PrintlnWarn("\nClient Extensions are not supported with Maven")
	}
}
