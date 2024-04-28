package scaffold

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/helper"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/manifoldco/promptui"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

const (
	ClientExtensionSamplePrefix      = "liferay-sample-"
	ClientExtensionSampleProjectName = "liferay-client-extensions-samples"
)

func CreateClientExtension(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := helper.FetchClientExtensionSamples(config.GetConfigPath()); err != nil {
		helper.HandleClientExtensionsOffline(config.GetConfigPath())
	}

	clientExtensionSamplesPath := filepath.Join(config.GetConfigPath(), ClientExtensionSampleProjectName)
	templates := getTemplateNames(clientExtensionSamplesPath)
	clientExtensionsWorkspaceDir := filepath.Join(liferayWorkspace, "client-extensions")

	promptTemplate := promptui.Select{
		Label: "Choose a template",
		Items: templates,
	}

	_, template, err := promptTemplate.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	template = filepath.Join(clientExtensionSamplesPath, ClientExtensionSamplePrefix+template)

	var name string
	if len(args) >= 1 && len(args[0]) > 0 {
		name = args[0]
	} else {
		promptName := promptui.Prompt{
			Label: "Choose a name",
			Validate: func(input string) error {
				if len(input) == 0 {
					return errors.New("the name cannot be empty")
				}
				return nil
			},
		}

		name, err = promptName.Run()

		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	name = strcase.ToKebab(strings.ToLower(name))
	clientExtensionDir := filepath.Join(clientExtensionsWorkspaceDir, name)

	fileutil.CreateDirs(clientExtensionDir)

	if err := cp.Copy(template, clientExtensionDir); err != nil {
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
		logger.PrintWarn("\nClient Extensions are not supported with Maven")
	}
}

func getTemplateNames(clientExtensionSamplesPath string) []string {
	sampleDirs, err := os.ReadDir(clientExtensionSamplesPath)

	if err != nil {
		logger.Fatal(err.Error())
	}

	var samples []string

	for _, sampleDir := range sampleDirs {
		if sampleDir.IsDir() && strings.Contains(sampleDir.Name(), ClientExtensionSamplePrefix) {
			samples = append(samples, strings.Split(sampleDir.Name(), ClientExtensionSamplePrefix)[1])
		}
	}

	return samples
}
