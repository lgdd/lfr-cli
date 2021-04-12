package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nxadm/tail"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/lgdd/deba/pkg/assets"
	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/printutil"
)

func CopyFromAssets(sourcePath, destPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	source, err := assets.Templates.Open(sourcePath)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	printutil.Success("create ")
	fmt.Printf("%s\n", destPath)
}

func UpdateWithData(pomPath string, metadata *project.Metadata) error {
	pomContent, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return err
	}

	tpl, err := template.New(pomPath).Parse(string(pomContent))
	if err != nil {
		return err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, metadata)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pomPath, result.Bytes(), 0664)
	if err != nil {
		return err
	}

	return nil
}

func VerifyCurrentDirAsWorkspace(build string) bool {
	files := make(map[string]void)
	dir, err := os.Getwd()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files[strings.Split(path, dir)[1]] = void{}
		return nil
	})

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	switch {
	case build == project.Gradle && isGradleWorkspace(files):
		return true
	case build == project.Maven && isMavenWorkspace(files):
		return true
	case build == project.Gradle && isMavenWorkspace(files):
		printutil.Warning(fmt.Sprintln("Oops! It looks like you're trying to do Gradle stuff in a Maven workspace."))
		fmt.Print("Try again with the flag: ")
		printutil.Info("-b maven\n")
		os.Exit(1)
	case build == project.Maven && isGradleWorkspace(files):
		printutil.Warning(fmt.Sprintln("Oops! It looks like you're trying to do Maven stuff in a Gradle workspace."))
		fmt.Print("Try again with the flag: ")
		printutil.Info("-b gradle\n")
		fmt.Print("or without the flag: ")
		printutil.Info("-b maven\n")
		os.Exit(1)
	}
	return false
}

func isGradleWorkspace(files map[string]void) bool {
	sep := string(os.PathSeparator)
	expectedFiles := []string{
		sep + "configs",
		sep + "gradle.properties",
		sep + "settings.gradle",
		sep + "gradle" + sep + "wrapper",
		sep + "build.gradle",
		sep + "gradlew",
		sep + "platform.bndrun",
	}
	for _, expectedFile := range expectedFiles {
		if _, ok := files[expectedFile]; !ok {
			return false
		}
	}
	return true
}

func isMavenWorkspace(files map[string]void) bool {
	sep := string(os.PathSeparator)
	expectedFiles := []string{
		sep + "configs",
		sep + ".mvn" + sep + "wrapper",
		sep + "pom.xml",
		sep + "mvnw",
		sep + "platform.bndrun",
	}
	for _, expectedFile := range expectedFiles {
		if _, ok := files[expectedFile]; !ok {
			return false
		}
	}
	return true
}

func FindFileInParent(fileName string) (string, error) {

	dir, err := os.Getwd()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
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
	filePath := ""

	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == fileName {
				filePath = path
			}
			return nil
		})

	if err != nil {
		return filePath, err
	}

	if filePath == "" {
		return filePath, fmt.Errorf("%s not found in directories under %s", fileName, dirPath)
	}

	return filePath, nil
}

func GetLiferayWorkspacePath() (string, error) {
	workspace, err := FindFileInParent("platform.bndrun")

	if err != nil {
		return "", errors.New("couldn't find Liferay Workspace in the current directory or any parent directory")
	}

	return filepath.Dir(workspace), nil
}

func GetLiferayHomePath() (string, error) {
	workspacePath, err := GetLiferayWorkspacePath()

	if err != nil {
		return "", err
	}

	liferayHome, err := FindFileParentInDir(workspacePath, ".liferay-home")

	if err != nil {
		return "", errors.New("couldn't find Liferay Home in the current directory or any subdirectory")
	}

	return liferayHome, nil
}

func GetTomcatScriptPath(script string) (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n\n", err.Error()))
		fmt.Println("Did you initialize the bundle from the root of your Liferay Workspace?")
		os.Exit(1)
	}

	scriptName := fmt.Sprintf("%s.sh", script)

	if runtime.GOOS == "windows" {
		scriptName = fmt.Sprintf("%s.bat", script)

	}

	scriptParentDir, err := FindFileParentInDir(liferayHome, scriptName)

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	scriptPath := filepath.Join(scriptParentDir, scriptName)

	return scriptPath, nil
}

func GetCatalinaLogFile() (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
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

type void struct{}
