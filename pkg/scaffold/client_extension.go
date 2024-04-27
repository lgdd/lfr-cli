package scaffold

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/manifoldco/promptui"
	cp "github.com/otiai10/copy"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	ClientExtensionSamplePrefix      = "liferay-sample-"
	ClientExtensionSampleProjectName = "liferay-client-extensions-samples"
)

func CreateClientExtension(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		panic(err)
	}

	if err := FetchClientExtensionSamples(config.GetConfigPath()); err != nil {
		HandleClientExtensionsOffline(config.GetConfigPath())
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
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
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
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}

	name = strcase.ToKebab(strings.ToLower(name))
	clientExtensionDir := filepath.Join(clientExtensionsWorkspaceDir, name)

	fileutil.CreateDirs(clientExtensionDir)

	if err := cp.Copy(template, clientExtensionDir); err != nil {
		panic(err)
	}

	_ = filepath.Walk(clientExtensionDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				printutil.Success("created ")
				fmt.Printf("%s\n", path)
			}
			return nil
		})

	if fileutil.IsMavenWorkspace(liferayWorkspace) {
		printutil.Warning("\nClient Extensions are not supported with Maven")
	}
}

func getTemplateNames(clientExtensionSamplesPath string) []string {
	sampleDirs, err := os.ReadDir(clientExtensionSamplesPath)

	if err != nil {
		panic(err)
	}

	var samples []string

	for _, sampleDir := range sampleDirs {
		if sampleDir.IsDir() && strings.Contains(sampleDir.Name(), ClientExtensionSamplePrefix) {
			samples = append(samples, strings.Split(sampleDir.Name(), ClientExtensionSamplePrefix)[1])
		}
	}

	return samples
}

func FetchClientExtensionSamples(destination string) error {
	clientExtensionsSamplesPath := filepath.Join(destination, ClientExtensionSampleProjectName)

	// Clone & checkout if ~/.lfr/liferay-portal does not exist
	if _, err := os.Stat(filepath.Join(destination, ClientExtensionSampleProjectName)); err != nil {
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Fetching samples"),
			progressbar.OptionSpinnerType(11))

		bar.Add(1)
		var gitProject strings.Builder
		gitProject.WriteString("https://github.com/lgdd/")
		gitProject.WriteString(ClientExtensionSampleProjectName)

		gitClone := exec.Command("git", "clone", "--depth", "1", gitProject.String())
		gitClone.Dir = destination

		if err := gitClone.Run(); err != nil {
			bar.Clear()
			return err
		}
		bar.Clear()
	} else {
		// Repo already exists, try to update
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Updating samples"),
			progressbar.OptionSpinnerType(11))

		bar.Add(1)

		gitPull := exec.Command("git", "pull")
		gitPull.Dir = clientExtensionsSamplesPath

		if err := gitPull.Run(); err != nil {
			bar.Clear()
			return err
		}

		bar.Clear()
	}
	return nil
}

func HandleClientExtensionsOffline(configPath string) {
	if _, err := os.Stat(filepath.Join(configPath, ClientExtensionSampleProjectName)); err != nil {
		printutil.Warning("Couldn't fetch client extensions samples from GitHub.\n")
		fmt.Println("Copying embedded versions from the CLI instead.")
		err = fileutil.CreateDirsFromAssets("tpl/client_extension", configPath)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		err = fileutil.CreateFilesFromAssets("tpl/client_extension", configPath)
		if err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		oldGitDirectory := filepath.Join(configPath, ClientExtensionSampleProjectName, "git")
		newGitDirectory := filepath.Join(configPath, ClientExtensionSampleProjectName, ".git")
		if err := os.Rename(oldGitDirectory, newGitDirectory); err != nil {
			printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
	} else {
		printutil.Warning("Couldn't update client extensions samples from GitHub.\n")
		fmt.Println("Using latest versions fetched.")
	}
}
