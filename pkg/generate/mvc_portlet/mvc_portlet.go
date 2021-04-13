package mvc_portlet

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/fileutil"
	"github.com/lgdd/deba/pkg/util/printutil"
)

type MvcPortletData struct {
	Package          string
	Name             string
	CamelCaseName    string
	WorkspaceName    string
	WorkspacePackage string
	PortletIdKey     string
	PortletIdValue   string
}

func Generate(name string) {
	sep := string(os.PathSeparator)
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	base := filepath.Join(liferayWorkspace, "modules")
	portletBase := filepath.Join(base, name)
	packagePath := strings.ReplaceAll(name, "-", string(os.PathSeparator))
	camelCaseName := strings.ReplaceAll(name, "-", " ")
	camelCaseName = strings.Title(camelCaseName)
	camelCaseName = strings.ReplaceAll(camelCaseName, " ", "")
	workspaceSplit := strings.Split(liferayWorkspace, sep)
	workspaceName := workspaceSplit[len(workspaceSplit)-1]
	workspacePackage := strings.ReplaceAll(workspaceName, "-", ".")

	dirs := []string{
		name,
		filepath.Join(name, "src", "main", "java", packagePath),
		filepath.Join(name, "src", "main", "java", packagePath, "constants"),
		filepath.Join(name, "src", "main", "resources", "content"),
		filepath.Join(name, "src", "main", "resources", "META-INF", "resources", "css"),
	}

	files := map[string]string{
		"tmpl/mvc-portlet/gitignore": filepath.Join(portletBase, ".gitignore"),
		"tmpl/mvc-portlet/bnd.bnd":   filepath.Join(portletBase, "bnd.bnd"),
		"tmpl/mvc-portlet/src/java/Portlet.java": filepath.Join(portletBase, "src", "main", "java",
			packagePath, fmt.Sprintf("%s.java", camelCaseName)),
		"tmpl/mvc-portlet/src/java/PortletKeys.java": filepath.Join(portletBase, "src", "main", "java",
			packagePath, "constants", fmt.Sprintf("%sKeys.java", camelCaseName)),
		"tmpl/mvc-portlet/src/resources/Language.properties": filepath.Join(portletBase, "src", "main",
			"resources", "content", "Language.properties"),
		"tmpl/mvc-portlet/src/resources/init.jsp": filepath.Join(portletBase, "src", "main",
			"resources", "META-INF", "resources", "init.jsp"),
		"tmpl/mvc-portlet/src/resources/view.jsp": filepath.Join(portletBase, "src", "main",
			"resources", "META-INF", "resources", "view.jsp"),
		"tmpl/mvc-portlet/src/resources/main.scss": filepath.Join(portletBase, "src", "main",
			"resources", "META-INF", "resources", "css", "main.scss"),
	}

	if fileutil.IsGradleWorkspace() {
		files["tmpl/mvc-portlet/build.gradle"] = filepath.Join(portletBase, "build.gradle")
	}

	if fileutil.IsMavenWorkspace() {
		files["tmpl/mvc-portlet/pom.xml"] = filepath.Join(portletBase, "pom.xml")

		pomParentPath := filepath.Join(portletBase, "../pom.xml")
		pomParent, err := os.Open(pomParentPath)
		if err != nil {
			printutil.Error(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
		defer pomParent.Close()

		byteValue, _ := ioutil.ReadAll(pomParent)

		var pom project.Pom
		xml.Unmarshal(byteValue, &pom)

		modules := append(pom.Modules.Module, name)
		pom.Modules.Module = modules
		pom.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
		pom.SchemaLocation = "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"

		finalPomBytes, _ := xml.MarshalIndent(pom, "", "  ")
		finalPom := fileutil.NormalizeXmlString(string(finalPomBytes))

		err = ioutil.WriteFile(pomParentPath, []byte(project.XMLHeader+finalPom), 0644)

		if err != nil {
			printutil.Error(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}

		printutil.Warning("update ")
		fmt.Printf("%s\n", pomParentPath)
	}

	var wg sync.WaitGroup

	wg.Add(len(dirs))
	for _, dir := range dirs {
		go fileutil.CreateDirs(filepath.Join(base, dir), &wg)
	}
	wg.Wait()

	wg.Add(len(files))
	for source, dest := range files {
		go fileutil.CopyFromAssets(source, dest, &wg)
	}
	wg.Wait()

	portletIdKey := strings.ReplaceAll(name, "-", "_")
	portletIdKey = strings.ToUpper(portletIdKey)
	portletIdValue := strings.ToLower(portletIdKey) + "_" + camelCaseName

	wg.Add(len(files))
	for _, dest := range files {
		go updateMvcPortletWithData(&wg, dest, &MvcPortletData{
			Package:          strings.ReplaceAll(name, "-", "."),
			Name:             name,
			CamelCaseName:    camelCaseName,
			PortletIdKey:     portletIdKey,
			PortletIdValue:   portletIdValue,
			WorkspaceName:    workspaceName,
			WorkspacePackage: workspacePackage,
		})
	}
	wg.Wait()
}

func updateMvcPortletWithData(wg *sync.WaitGroup, file string, data *MvcPortletData) {
	defer wg.Done()

	content, err := ioutil.ReadFile(file)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	tpl, err := template.New(file).Parse(string(content))
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, data)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	err = ioutil.WriteFile(file, result.Bytes(), 0664)
	if err != nil {
		printutil.Error(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
}
