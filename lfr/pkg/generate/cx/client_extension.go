package cx

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/lgdd/liferay-cli/lfr/pkg/config"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
	"github.com/manifoldco/promptui"
	cp "github.com/otiai10/copy"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	ClientExtensionSamplePrefix = "liferay-sample-"
)

func Generate(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		panic(err)
	}

	fetchClientExtensionSamples(config.GetConfigPath())

	clientExtensionSamplesPath := filepath.Join(config.GetConfigPath(), "liferay-portal", "workspaces", "liferay-sample-workspace", "client-extensions")
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

func fetchClientExtensionSamples(destination string) {
	liferayPortalPath := filepath.Join(destination, "liferay-portal")

	// Clone & checkout if ~/.lfr/liferay-portal does not exist
	if _, err := os.Stat(filepath.Join(destination, "liferay-portal")); err != nil {
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Fetching samples"),
			progressbar.OptionSpinnerType(11))

		bar.Add(1)

		gitClone := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--no-checkout", "https://github.com/liferay/liferay-portal")
		gitClone.Dir = destination

		if err := gitClone.Run(); err != nil {
			panic(err)
		}

		gitSparseCheckoutInit := exec.Command("git", "sparse-checkout", "init", "--no-cone")
		gitSparseCheckoutInit.Dir = liferayPortalPath

		if err := gitSparseCheckoutInit.Run(); err != nil {
			panic(err)
		}

		gitSparseCheckoutSet := exec.Command("git", "sparse-checkout", "set", "workspaces/liferay-sample-workspace/client-extensions")
		gitSparseCheckoutSet.Dir = liferayPortalPath

		if err := gitSparseCheckoutSet.Run(); err != nil {
			panic(err)
		}

		gitCheckout := exec.Command("git", "checkout")
		gitCheckout.Dir = liferayPortalPath

		if err := gitCheckout.Run(); err != nil {
			panic(err)
		}
		bar.Clear()
	} else {
		// Repo already exists, try to update
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Updating samples"),
			progressbar.OptionSpinnerType(11))

		bar.Add(1)

		gitPull := exec.Command("git", "pull")
		gitPull.Dir = liferayPortalPath

		if err := gitPull.Run(); err != nil {
			panic(err)
		}

		bar.Clear()
	}
}
