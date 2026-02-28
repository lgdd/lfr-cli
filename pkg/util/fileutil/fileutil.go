// Package fileutil provides helpers for creating and manipulating files and
// directories, processing Go templates, working with Maven pom.xml files, and
// locating Liferay workspace and bundle paths on the filesystem.
package fileutil

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/charmbracelet/huh/spinner"
	"github.com/magiconair/properties"
	"github.com/nxadm/tail"
	"github.com/spf13/viper"

	"github.com/lgdd/lfr-cli/internal/assets"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

// XMLHeader is the first line to be written in an XML file
const (
	XMLHeader = `<?xml version="1.0"?>` + "\n"
)

// Pom represents the common structure of a Maven pom file
type Pom struct {
	XMLName        xml.Name `xml:"project"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xmlns:schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	Parent         struct {
		GroupId      string `xml:"groupId"`
		ArtifactId   string `xml:"artifactId"`
		Version      string `xml:"version"`
		RelativePath string `xml:"relativePath"`
	} `xml:"parent"`
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Name       string `xml:"name"`
	Packaging  string `xml:"packaging"`
	Modules    struct {
		Module []string `xml:"module"`
	} `xml:"modules"`
}

// WorkspacePom represents the common structure of a parent Maven pom file in a Liferay Workspace
type WorkspacePom struct {
	XMLName        xml.Name `xml:"project"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xmlns:schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	Parent         struct {
		GroupId      string `xml:"groupId"`
		ArtifactId   string `xml:"artifactId"`
		Version      string `xml:"version"`
		RelativePath string `xml:"relativePath"`
	} `xml:"parent"`
	ArtifactId string `xml:"artifactId"`
	Name       string `xml:"name"`
	Packaging  string `xml:"packaging"`
	Modules    struct {
		Module []string `xml:"module"`
	} `xml:"modules"`
	Properties struct {
		LiferayBomVersion          string `xml:"liferay.bom.version"`
		LiferayDockerImage         string `xml:"liferay.docker.image"`
		LiferayWorkspaceBundleURL  string `xml:"liferay.workspace.bundle.url"`
		LiferayRepositoryURL       string `xml:"liferay.repository.url"`
		ProjectBuildSourceEncoding string `xml:"project.build.sourceEncoding"`
	} `xml:"properties"`
}

// CreateDirs creates all directories in path, including any necessary parents.
func CreateDirs(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// CreateFiles creates empty files for each path in the given list, concurrently.
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
		logger.Fatal(err.Error())
	}
}

// CreateDirsFromAssets walks the embedded template assets under assetsRoot and
// recreates the directory structure under baseDest.
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

// CreateFilesFromAssets walks the embedded template assets under assetsRoot and
// copies each file to the corresponding path under baseDest.
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

// CopyFromTemplates copies an embedded template asset at sourcePath to destPath on disk.
func CopyFromTemplates(sourcePath, destPath string) {
	source, err := assets.Templates.Open(sourcePath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		logger.Fatal(err.Error())
	}

}

// CopyFromAssets copies an embedded asset at sourcePath to destPath on disk.
// It calls wg.Done() when the copy completes and is designed for concurrent use.
func CopyFromAssets(sourcePath, destPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	source, err := assets.Templates.Open(sourcePath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// AppendModuleToPom adds a module entry to an existing Maven pom.xml file
func AppendModuleToPom(pomPath, moduleName string) error {
	pomFile, err := os.Open(pomPath)
	if err != nil {
		return err
	}
	defer pomFile.Close()

	byteValue, err := io.ReadAll(pomFile)
	if err != nil {
		return err
	}

	var pom Pom
	if err = xml.Unmarshal(byteValue, &pom); err != nil {
		return err
	}

	pom.Modules.Module = append(pom.Modules.Module, moduleName)
	pom.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	pom.SchemaLocation = "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"

	finalPomBytes, err := xml.MarshalIndent(pom, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pomPath, []byte(XMLHeader+string(finalPomBytes)), 0644)
}

// UpdateWithData renders file as a Go template with the provided data and
// writes the result back to the same file.
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

// IsInWorkspaceDir reports whether the current working directory is inside a Liferay workspace.
func IsInWorkspaceDir() bool {
	_, err := GetLiferayWorkspacePath()

	return err == nil
}

// IsGradleWorkspace reports whether the workspace at path is a Gradle workspace.
func IsGradleWorkspace(path string) bool {
	expectedFiles := []string{
		filepath.Join(path, "configs"),
		filepath.Join(path, "gradle.properties"),
		filepath.Join(path, "settings.gradle"),
		filepath.Join(path, "gradlew"),
	}

	return FilesExist(expectedFiles)
}

// IsMavenWorkspace reports whether the workspace at path is a Maven workspace.
func IsMavenWorkspace(path string) bool {
	expectedFiles := []string{
		filepath.Join(path, "configs"),
		filepath.Join(path, "pom.xml"),
		filepath.Join(path, "mvnw"),
	}

	return FilesExist(expectedFiles)
}

// FilesExist reports whether every path in files exists on disk.
func FilesExist(files []string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// FindFileInParent searches the working directory and each of its parents for a
// file named fileName, returning its full path when found.
func FindFileInParent(fileName string) (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		logger.Fatal(err.Error())
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

// FindFileParentInDir returns the parent directory of the first file named
// fileName found anywhere under dirPath.
func FindFileParentInDir(dirPath string, fileName string) (string, error) {
	filePath, err := FindFileInDir(dirPath, fileName)

	if err != nil {
		return "", err
	}

	return filepath.Dir(filePath), nil
}

// FindFileInDir recursively searches dirPath for a file named fileName and
// returns its full path. A spinner is shown while scanning.
func FindFileInDir(dirPath string, fileName string) (string, error) {
	targetFilePath := ""

	scan := func() {
		filepath.Walk(dirPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if err != nil {
					logger.Fatal(err.Error())
				}
				if info.Name() == fileName {
					targetFilePath = path
				}
				return nil
			})
	}

	_ = spinner.New().
		Title(fmt.Sprintf("Scanning files under %s", dirPath)).
		Action(scan).
		Accessible(viper.GetBool(conf.OutputAccessible)).
		Run()

	if targetFilePath == "" {
		return targetFilePath, fmt.Errorf("%s not found in directories under %s", fileName, dirPath)
	}

	return targetFilePath, nil
}

// GetLiferayWorkspacePath returns the root path of the Liferay workspace by
// locating the platform.bndrun file in the current directory or any parent.
func GetLiferayWorkspacePath() (string, error) {
	workspace, err := FindFileInParent("platform.bndrun")

	if err != nil {
		return "", errors.New("couldn't find Liferay Workspace in the current directory or any parent directory")
	}

	return filepath.Dir(workspace), nil
}

// GetLiferayHomePath returns the root path of the Liferay bundle by locating
// the .liferay-home marker file within or near the workspace.
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

// GetTomcatScriptPath returns the path to the named Tomcat script (e.g. "catalina")
// inside the Liferay bundle, using the appropriate extension for the current OS.
func GetTomcatScriptPath(script string) (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		logger.Errorf("%s\n", err.Error())
		logger.Warn("Did you initialize the Tomcat bundle from the root of your Liferay Workspace?")
		os.Exit(1)
	}

	scriptName := fmt.Sprintf("%s.sh", script)

	if runtime.GOOS == "windows" {
		scriptName = fmt.Sprintf("%s.bat", script)
	}

	scriptParentDir, err := FindFileParentInDir(liferayHome, scriptName)

	if err != nil {
		logger.Fatal(err.Error())
	}

	scriptPath := filepath.Join(scriptParentDir, scriptName)

	return scriptPath, nil
}

// GetTomcatPath returns the Tomcat home directory path within the Liferay bundle.
func GetTomcatPath() (string, error) {
	catalinaScriptPath, err := GetTomcatScriptPath("catalina")

	if err != nil {
		return "", err
	}
	catalinaScriptPathSplit := strings.Split(catalinaScriptPath, string(os.PathSeparator))
	return strings.Join(catalinaScriptPathSplit[:len(catalinaScriptPathSplit)-2], string(os.PathSeparator)), nil
}

// GetCatalinaLogFile returns the path to the main Tomcat log file (catalina.out).
func GetCatalinaLogFile() (string, error) {
	liferayHome, err := GetLiferayHomePath()

	if err != nil {
		logger.Fatal(err.Error())
	}

	return FindFileInDir(liferayHome, "catalina.out")
}

// Tail streams the contents of logFile to stdout. When follow is true, new
// lines are printed as they are appended to the file.
func Tail(logFile string, follow bool) {
	t, err := tail.TailFile(logFile, tail.Config{Follow: follow})
	if err != nil {
		logger.Fatal(err.Error())
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

// DirSize returns the total size in bytes of all files under path.
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

// GetLiferayWorkspaceProduct returns the Liferay workspace product string
// (e.g. "portal-7.4.3.120" or "dxp-2024.q4.0") for the workspace at workspacePath.
func GetLiferayWorkspaceProduct(workspacePath string) (string, error) {
	if IsMavenWorkspace(workspacePath) {
		pomWorkspacePath := filepath.Join(workspacePath, "pom.xml")
		pomWorkspace, err := os.Open(pomWorkspacePath)

		if err != nil {
			return "", err
		}

		defer pomWorkspace.Close()
		byteValue, _ := io.ReadAll(pomWorkspace)

		var pom WorkspacePom
		err = xml.Unmarshal(byteValue, &pom)

		if err != nil {
			return "", err
		}

		var productBuilder strings.Builder
		editionRegex := regexp.MustCompile(`(portal|dxp)`)
		edition := editionRegex.FindString(pom.Properties.LiferayWorkspaceBundleURL)

		productBuilder.WriteString(edition)
		productBuilder.WriteString("-")
		productBuilder.WriteString(pom.Properties.LiferayBomVersion)

		return productBuilder.String(), nil
	}

	if IsGradleWorkspace(workspacePath) {
		gradlePropsPath := filepath.Join(workspacePath, "gradle.properties")
		gradleProps := properties.MustLoadFile(gradlePropsPath, properties.UTF8)
		return gradleProps.GetString("liferay.workspace.product", ""), nil
	}

	return "", nil
}

// GetLiferayWorkspaceProductVersion returns the major Liferay product version
// (e.g. "7.4") for the workspace at workspacePath.
func GetLiferayWorkspaceProductVersion(workspacePath string) (string, error) {
	product, err := GetLiferayWorkspaceProduct(workspacePath)

	if err != nil {
		return "", err
	}

	majorVersionRegex := regexp.MustCompile(`7\.\d`)
	quarterlyVersionRegex := regexp.MustCompile(`-\d+\.q\d`)

	majorVersion := majorVersionRegex.FindString(product)
	quarterlyVersion := quarterlyVersionRegex.FindString(product)

	if len(majorVersion) == 0 && len(quarterlyVersion) == 0 {
		return "", fmt.Errorf("liferay workspace product version not found")
	}

	if len(quarterlyVersion) > 0 {
		return "7.4", nil
	}

	return majorVersion, nil
}

// GetLiferayWorkspaceProductEdition returns the edition ("dxp" or "portal") of
// the Liferay product for the workspace at workspacePath.
func GetLiferayWorkspaceProductEdition(workspacePath string) (string, error) {
	product, err := GetLiferayWorkspaceProduct(workspacePath)

	if err != nil {
		return "", err
	}

	if strings.Contains(product, "dxp") || strings.Contains(product, ".q") {
		return "dxp", nil
	}

	if strings.Contains(product, "portal") {
		return "portal", nil
	}

	return "", fmt.Errorf("liferay workspace product edition not found")
}
