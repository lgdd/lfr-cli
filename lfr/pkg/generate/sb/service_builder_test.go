package sb

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lgdd/liferay-cli/lfr/pkg/generate/workspace"
	"github.com/lgdd/liferay-cli/lfr/pkg/project"
)

func Test_GenerateServiceBuilder_WithGradle_ShouldCreateExpectedFiles(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.3")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedModulePath := filepath.Join(liferayWorkspace, "modules", name)
	if _, err := os.Stat(expectedModulePath); err != nil {
		t.Fatal(err)
	}
	expectedApiPath := filepath.Join(liferayWorkspace, "modules", name, name+"-api")
	if _, err := os.Stat(expectedApiPath); err != nil {
		t.Fatal(err)
	}
	expectedServicePath := filepath.Join(liferayWorkspace, "modules", name, name+"-service")
	if _, err := os.Stat(expectedServicePath); err != nil {
		t.Fatal(err)
	}
	expectedServiceXML := filepath.Join(expectedServicePath, "service.xml")
	if _, err := os.Stat(expectedServiceXML); err != nil {
		t.Fatal(err)
	}

}

func Test_GenerateServiceBuilder_WithMaven_ShouldCreateExpectedFiles(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Maven, "7.3")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedModulePath := filepath.Join(liferayWorkspace, "modules", name)
	if _, err := os.Stat(expectedModulePath); err != nil {
		t.Fatal(err)
	}
	expectedApiPath := filepath.Join(liferayWorkspace, "modules", name, name+"-api")
	if _, err := os.Stat(expectedApiPath); err != nil {
		t.Fatal(err)
	}
	expectedServicePath := filepath.Join(liferayWorkspace, "modules", name, name+"-service")
	if _, err := os.Stat(expectedServicePath); err != nil {
		t.Fatal(err)
	}
	expectedServiceXML := filepath.Join(expectedServicePath, "service.xml")
	if _, err := os.Stat(expectedServiceXML); err != nil {
		t.Fatal(err)
	}
	expectedPomXMLs := []string{
		filepath.Join(expectedModulePath, "pom.xml"),
		filepath.Join(expectedApiPath, "pom.xml"),
		filepath.Join(expectedServicePath, "pom.xml"),
	}
	for _, pomXML := range expectedPomXMLs {
		if _, err := os.Stat(pomXML); err != nil {
			t.Fatal(err)
		}
	}
	modulesPomXML, err := os.Open(filepath.Join(liferayWorkspace, "modules", "pom.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer modulesPomXML.Close()
	byteValue, _ := ioutil.ReadAll(modulesPomXML)
	var pom project.Pom
	err = xml.Unmarshal(byteValue, &pom)
	if err != nil {
		t.Fatal(err)
	}
	foundExpectedModule := false
	for _, module := range pom.Modules.Module {
		if name == module {
			foundExpectedModule = true
		}
	}
	if !foundExpectedModule {
		t.Fatal("the module was not found in the pom.xml of the modules folder")
	}
}

func Test_GenerateServiceBuilder_With73_ShouldContainCorrespondingDoctype(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.3")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedDoctype := "<!DOCTYPE service-builder PUBLIC \"-//Liferay//DTD Service Builder 7.3.0//EN\" \"http://www.liferay.com/dtd/liferay-service-builder_7_3_0.dtd\">"
	serviceXML, err := os.Open(filepath.Join(liferayWorkspace, "modules", name, name+"-service", "service.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer serviceXML.Close()
	serviceXMLBytes, _ := ioutil.ReadAll(serviceXML)
	serviceXMLContent := string(serviceXMLBytes)
	if !strings.Contains(serviceXMLContent, expectedDoctype) {
		t.Fatal("valid doctype for 7.3 wasn't found in service.xml")
	}
}

func Test_GenerateServiceBuilder_With72_ShouldContainCorrespondingDoctype(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.2")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedDoctype := "<!DOCTYPE service-builder PUBLIC \"-//Liferay//DTD Service Builder 7.2.0//EN\" \"http://www.liferay.com/dtd/liferay-service-builder_7_2_0.dtd\">"
	serviceXML, err := os.Open(filepath.Join(liferayWorkspace, "modules", name, name+"-service", "service.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer serviceXML.Close()
	serviceXMLBytes, _ := ioutil.ReadAll(serviceXML)
	serviceXMLContent := string(serviceXMLBytes)
	if !strings.Contains(serviceXMLContent, expectedDoctype) {
		t.Fatal("valid doctype for 7.2 wasn't found in service.xml")
	}
}

func Test_GenerateServiceBuilder_With71_ShouldContainCorrespondingDoctype(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.1")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedDoctype := "<!DOCTYPE service-builder PUBLIC \"-//Liferay//DTD Service Builder 7.1.0//EN\" \"http://www.liferay.com/dtd/liferay-service-builder_7_1_0.dtd\">"
	serviceXML, err := os.Open(filepath.Join(liferayWorkspace, "modules", name, name+"-service", "service.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer serviceXML.Close()
	serviceXMLBytes, _ := ioutil.ReadAll(serviceXML)
	serviceXMLContent := string(serviceXMLBytes)
	if !strings.Contains(serviceXMLContent, expectedDoctype) {
		t.Fatal("valid doctype for 7.1 wasn't found in service.xml")
	}
}

func Test_GenerateServiceBuilder_With70_ShouldContainCorrespondingDoctype(t *testing.T) {
	liferayWorkspace := filepath.Join(t.TempDir(), "liferay-workspace")
	err := workspace.Generate(liferayWorkspace, project.Gradle, "7.0")
	if err != nil {
		t.Fatal(err)
	}
	name := "example-app"
	Generate(liferayWorkspace, name)
	expectedDoctype := "<!DOCTYPE service-builder PUBLIC \"-//Liferay//DTD Service Builder 7.0.0//EN\" \"http://www.liferay.com/dtd/liferay-service-builder_7_0_0.dtd\">"
	serviceXML, err := os.Open(filepath.Join(liferayWorkspace, "modules", name, name+"-service", "service.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer serviceXML.Close()
	serviceXMLBytes, _ := ioutil.ReadAll(serviceXML)
	serviceXMLContent := string(serviceXMLBytes)
	if !strings.Contains(serviceXMLContent, expectedDoctype) {
		t.Fatal("valid doctype for 7.1 wasn't found in service.xml")
	}
}
