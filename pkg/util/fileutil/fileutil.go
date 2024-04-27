package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/nxadm/tail"
	progressbar "github.com/schollz/progressbar/v3"

	"github.com/lgdd/lfr-cli/internal/assets"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
)

// Create all the directories of a given path
func CreateDirs(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

// Create all the files from a given list of paths
func CreateFiles(paths []string) {
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go createFile(path, &wg)
	}
	wg.Wait()
}

func createFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := os.Create(path)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

// Walk through the template assets and create all the directories contained in it
func CreateDirsFromAssets(assetsRoot, baseDest string) error {
	return fs.WalkDir(assets.Templates, assetsRoot, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if file.IsDir() {
			relativePath := strings.Split(path, assetsRoot)[1]
			if len(relativePath) > 0 {
				destPath := baseDest + relativePath
				CreateDirs(destPath)
			}
		}
		return nil
	})
}

// Walk through the template assets and create all the files contained in it
func CreateFilesFromAssets(assetsRoot, baseDest string) error {
	return fs.WalkDir(assets.Templates, assetsRoot, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			relativePath := strings.Split(path, assetsRoot)[1]
			if len(relativePath) > 0 {
				destPath := baseDest + relativePath
				CopyFromTemplates(path, destPath)
			}
		}
		return nil
	})
}

// Copy templates to a given destination
func CopyFromTemplates(sourcePath, destPath string) {
	source, err := assets.Templates.Open(sourcePath)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

}

// Copy assets to a given destination
func CopyFromAssets(sourcePath, destPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	source, err := assets.Templates.Open(sourcePath)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

// Update template files with given data
func UpdateWithData(file string, data interface{}) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tpl, err := template.New(file).Parse(string(content))
	if err != nil {
		return err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, data)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, result.Bytes(), 0664)
	if err != nil {
		return err
	}

	return nil
}

// Check if the current directory is under a Liferay workspace
func IsInWorkspaceDir() bool {
	_, err := GetLiferayWorkspacePath()

	return err == nil
}

// Check if the Liferay workspace is using Gradle
func IsGradleWorkspace(path string) bool {
	expectedFiles := []string{
		filepath.Join(path, "configs"),
		filepath.Join(path, "gradle.properties"),
		filepath.Join(path, "settings.gradle"),
		filepath.Join(path, "gradlew"),
	}

	return FilesExist(expectedFiles)
}

// Check if the Liferay workspace is using Maven
func IsMavenWorkspace(path string) bool {
	expectedFiles := []string{
		filepath.Join(path, "configs"),
		filepath.Join(path, "pom.xml"),
		filepath.Join(path, "mvnw"),
	}

	return FilesExist(expectedFiles)
}

// Check if a given list of files exist
func FilesExist(files []string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Find a file with a given name in parent directories
func FindFileInParent(fileName string) (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	var filePath string
	sep := string(os.PathSeparator)

	slice := strings.Split(dir, sep)
	slice = slice[1:]

	for len(slice) > 0 {
		elemsPath := append([]string{sep}, slice...)
		elemsPath = append(elemsPath, fileName)
		filePath = filepath.Join(elemsPath...)

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return filePath, nil
		}
		slice = slice[:len(slice)-1]
	}

	return "", fmt.Errorf("%s not found", fileName)
}

// Find a file with a given name in a given directory and walk through its parents
func FindFileParentInDir(dirPath string, fileName string) (string, error) {
	filePath, err := FindFileInDir(dirPath, fileName)

	if err != nil {
		return "", err
	}

	return filepath.Dir(filePath), nil
}

// Find a file with a given name under a give directory
func FindFileInDir(dirPath string, fileName string) (string, error) {
	targetFilePath := ""

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(fmt.Sprintf("Scanning files under %s", dirPath)),
		progressbar.OptionSpinnerType(11),
		progressbar.OptionThrottle(65*time.Millisecond))

	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			err = bar.Add(1)
			if err != nil {
				printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
				os.Exit(1)
			}
			if info.Name() == fileName {
				targetFilePath = path
			}
			return nil
		})

	if err != nil {
		return targetFilePath, err
	}

	err = bar.Clear()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	if targetFilePath == "" {
		return targetFilePath, fmt.Errorf("%s not found in directories under %s", fileName, dirPath)
	}

	return targetFilePath, nil
}

// Get Liferay workspace path
func GetLiferayWorkspacePath() (string, error) {
	workspace, err := FindFileInParent("platform.bndrun")

	if err != nil {
		return "", errors.New("couldn't find Liferay Workspace in the current directory or any parent directory")
	}

	return filepath.Dir(workspace), nil
}

// Get Liferay home path (i.e. root path of the Liferay bundle)
func GetLiferayHomePath() (string, error) {
	workingPath, err := GetLiferayWorkspacePath()

	if err != nil {
		workingPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	liferayHome, err := FindFileParentInDir(workingPath, ".liferay-home")

	if err != nil {
		liferayHome, err = FindFileInParent(".liferay-home")
		if err != nil {
			return "", errors.New("couldn't find a Liferay Tomcat bundle")
		}
	}

	if strings.HasSuffix(liferayHome, ".liferay-home") {
		liferayHomeSplit := strings.Split(liferayHome, string(os.PathSeparator))
		liferayHomeSplit = liferayHomeSplit[:len(liferayHomeSplit)-1]
		liferayHome = strings.Join(liferayHomeSplit, string(os.PathSeparator))
	}

	return liferayHome, nil
}

// Get the bin directory path under Tomcat home within the Liferay bundle
func GetTomcatScriptPath(script string) (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n\n", err.Error()))
		fmt.Println("Did you initialize the Tomcat bundle from the root of your Liferay Workspace?")
		os.Exit(1)
	}

	scriptName := fmt.Sprintf("%s.sh", script)

	if runtime.GOOS == "windows" {
		scriptName = fmt.Sprintf("%s.bat", script)
	}

	scriptParentDir, err := FindFileParentInDir(liferayHome, scriptName)

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	scriptPath := filepath.Join(scriptParentDir, scriptName)

	return scriptPath, nil
}

// Get the Tomcat home directory within the Liferay bundle
func GetTomcatPath() (string, error) {
	catalinaScriptPath, err := GetTomcatScriptPath("catalina")

	if err != nil {
		return "", err
	}
	catalinaScriptPathSplit := strings.Split(catalinaScriptPath, string(os.PathSeparator))
	return strings.Join(catalinaScriptPathSplit[:len(catalinaScriptPathSplit)-2], string(os.PathSeparator)), nil
}

// Get the main Tomcat log file
func GetCatalinaLogFile() (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	return FindFileInDir(liferayHome, "catalina.out")
}

// Tail a given log file with the option to follow updates
func Tail(logFile string, follow bool) {
	t, err := tail.TailFile(logFile, tail.Config{Follow: follow})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
