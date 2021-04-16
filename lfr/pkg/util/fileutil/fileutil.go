package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/nxadm/tail"

	"github.com/lgdd/liferay-cli/lfr/pkg/assets"
	"github.com/lgdd/liferay-cli/lfr/pkg/util/printutil"
)

func CreateDirs(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}

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

	printutil.Success("create ")
	fmt.Printf("%s\n", destPath)
}

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

	printutil.Success("create ")
	fmt.Printf("%s\n", destPath)
}

func UpdateWithData(file string, data interface{}) error {
	content, err := ioutil.ReadFile(file)
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

	err = ioutil.WriteFile(file, result.Bytes(), 0664)
	if err != nil {
		return err
	}

	return nil
}

func IsInWorkspaceDir() bool {
	_, err := GetLiferayWorkspacePath()

	if err != nil {
		return false
	}

	return true
}

func IsGradleWorkspace() bool {
	workspace, err := GetLiferayWorkspacePath()

	if err != nil {
		return false
	}

	expectedFiles := []string{
		filepath.Join(workspace, "configs"),
		filepath.Join(workspace, "gradle.properties"),
		filepath.Join(workspace, "settings.gradle"),
		filepath.Join(workspace, "build.gradle"),
		filepath.Join(workspace, "gradlew"),
	}

	return FilesExists(expectedFiles)
}

func IsMavenWorkspace() bool {
	workspace, err := GetLiferayWorkspacePath()

	if err != nil {
		return false
	}

	expectedFiles := []string{
		filepath.Join(workspace, "configs"),
		filepath.Join(workspace, "pom.xml"),
		filepath.Join(workspace, "mvnw"),
	}

	return FilesExists(expectedFiles)
}

func FilesExists(files []string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

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

func FindFileParentInDir(dirPath string, fileName string) (string, error) {
	filePath, err := FindFileInDir(dirPath, fileName)

	if err != nil {
		return "", err
	}

	return filepath.Dir(filePath), nil
}

func FindFileInDir(dirPath string, fileName string) (string, error) {
	targetFilePath := ""

	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == fileName {
				targetFilePath = path
			}
			return nil
		})

	if err != nil {
		return targetFilePath, err
	}

	if targetFilePath == "" {
		return targetFilePath, fmt.Errorf("%s not found in directories under %s", fileName, dirPath)
	}

	return targetFilePath, nil
}

func GetLiferayWorkspacePath() (string, error) {
	workspace, err := FindFileInParent("platform.bndrun")

	if err != nil {
		return "", errors.New("couldn't find Liferay Workspace in the current directory or any parent directory")
	}

	return filepath.Dir(workspace), nil
}

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

func GetTomcatPath() (string, error) {
	catalinaScriptPath, err := GetTomcatScriptPath("catalina")

	if err != nil {
		return "", err
	}
	catalinaScriptPathSplit := strings.Split(catalinaScriptPath, string(os.PathSeparator))
	return strings.Join(catalinaScriptPathSplit[:len(catalinaScriptPathSplit)-2], string(os.PathSeparator)), nil
}

func GetCatalinaLogFile() (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	return FindFileInDir(liferayHome, "catalina.out")
}

func Tail(logFile string, follow bool) {
	t, err := tail.TailFile(logFile, tail.Config{Follow: follow})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
