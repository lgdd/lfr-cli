package helper

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

var supportedJavaVersions = [2]int{8, 11}

// Checks if the Java version is supported by Liferay
func IsSupportedJavaVersion(javaVersion int) bool {
	for _, version := range supportedJavaVersions {
		if javaVersion == version {
			return true
		}
	}
	return false
}

func FetchClientExtensionSamples(destination string) error {
	clientExtensionsSamplesPath := filepath.Join(destination, conf.ClientExtensionSampleProjectName)

	// Clone & checkout if ~/.lfr/liferay-portal does not exist
	if _, err := os.Stat(filepath.Join(destination, conf.ClientExtensionSampleProjectName)); err != nil {
		var gitProject strings.Builder
		gitProject.WriteString("https://github.com/lgdd/")
		gitProject.WriteString(conf.ClientExtensionSampleProjectName)

		gitClone := exec.Command("git", "clone", "--depth", "1", gitProject.String())
		gitClone.Dir = destination

		if err := gitClone.Run(); err != nil {
			return err
		}
	} else {
		// Repo already exists, try to update
		go updateClientExtensionSamples(clientExtensionsSamplesPath)
	}
	return nil
}

func updateClientExtensionSamples(path string) {
	gitPull := exec.Command("git", "pull")
	gitPull.Dir = path

	if err := gitPull.Run(); err != nil {
		defer logger.Error(err.Error())
	}
}

func HandleClientExtensionsOffline(configPath string) {
	if _, err := os.Stat(filepath.Join(configPath, conf.ClientExtensionSampleProjectName)); err != nil {
		logger.PrintWarn("Couldn't fetch client extensions samples from GitHub.\n")
		logger.Println("Copying embedded versions from the CLI instead.")
		err = fileutil.CreateDirsFromAssets("tpl/client_extension", configPath)
		if err != nil {
			logger.Fatal(err.Error())
		}

		err = fileutil.CreateFilesFromAssets("tpl/client_extension", configPath)
		if err != nil {
			logger.Fatal(err.Error())
		}

		oldGitDirectory := filepath.Join(configPath, conf.ClientExtensionSampleProjectName, "git")
		newGitDirectory := filepath.Join(configPath, conf.ClientExtensionSampleProjectName, ".git")
		if err := os.Rename(oldGitDirectory, newGitDirectory); err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		logger.PrintWarn("Couldn't update client extensions samples from GitHub.\n")
		logger.Print("Using latest versions fetched.")
	}
}

func GetClientExtensionSampleNames() []string {
	if err := FetchClientExtensionSamples(conf.GetConfigPath()); err != nil {
		HandleClientExtensionsOffline(conf.GetConfigPath())
	}

	clientExtensionSamplesPath := filepath.Join(conf.GetConfigPath(), conf.ClientExtensionSampleProjectName)
	sampleDirs, err := os.ReadDir(clientExtensionSamplesPath)

	if err != nil {
		logger.Fatal(err.Error())
	}

	var samples []string

	for _, sampleDir := range sampleDirs {
		if sampleDir.IsDir() && strings.Contains(sampleDir.Name(), conf.ClientExtensionSamplePrefix) {
			samples = append(samples, strings.Split(sampleDir.Name(), conf.ClientExtensionSamplePrefix)[1])
		}
	}

	return samples
}
