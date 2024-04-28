package create

import (
	"errors"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/internal/config"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Cmd is the command 'create' which create the structure of a given module type
	Cmd = &cobra.Command{
		Use:     "create TYPE NAME",
		Aliases: []string{"c"},
		Short:   "Create a Liferay project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				promptCreateChoices(cmd, args)
			}
		},
	}
)

func init() {
	Cmd.AddCommand(createWorkspace)
	Cmd.AddCommand(createClientExtension)
	Cmd.AddCommand(createMvcPortlet)
	Cmd.AddCommand(createSpringPortlet)
	Cmd.AddCommand(createApiModule)
	Cmd.AddCommand(createCmdModule)
	Cmd.AddCommand(createServiceBuilder)
	Cmd.AddCommand(createRESTBuilder)
	Cmd.AddCommand(createDocker)
	config.Init()
	defaultPackage := viper.GetString(config.ModulePackage)
	Cmd.PersistentFlags().StringVarP(&metadata.PackageName, "package", "p", defaultPackage, "base package name")
}

func promptCreateChoices(cmd *cobra.Command, args []string) {
	promptTemplate := promptui.Select{
		Label: "Choose a template",
		Items: []string{"client-extension", "api", "command", "docker", "mvc-portlet", "rest-builder", "service-builder", "spring-mvc-portlet", "workspace"},
	}

	_, template, err := promptTemplate.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if template == "docker" {
		promptJavaVersion := promptui.Select{
			Label: "Choose a Java version",
			Items: []string{"8", "11"},
		}
		_, javaVersion, err := promptJavaVersion.Run()

		if err != nil {
			logger.Fatal(err.Error())
		}

		promptBuildOption := promptui.Select{
			Label: "Choose a build option",
			Items: []string{"Multi-stage", "Single stage"},
		}

		_, buildOption, err := promptBuildOption.Run()

		if err != nil {
			logger.Fatal(err.Error())
		}

		if buildOption == "Multi-stage" {
			os.Args = append(os.Args, template, "-j", javaVersion, "-m")
		} else {
			os.Args = append(os.Args, template, "-j", javaVersion)
		}

		err = cmd.Execute()

		if err != nil {
			logger.Fatal(err.Error())
		}

		return
	}

	promptName := promptui.Prompt{
		Label: "Choose a name",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("the name cannot be empty")
			}
			return nil
		},
	}

	name, err := promptName.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	workspacePackage, _ := metadata.GetGroupId()
	defaultPackageName := strings.Join([]string{workspacePackage, strcase.ToDelimited(name, '.')}, ".")

	if template == "workspace" {
		defaultPackageName = "org.acme"
	}

	promptPackageName := promptui.Prompt{
		Label:   "Choose a package name",
		Default: defaultPackageName,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("the name cannot be empty")
			}
			return nil
		},
	}

	packageName, err := promptPackageName.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if template == "spring-mvc-portlet" {
		promptTemplateEngine := promptui.Select{
			Label: "Choose a template engine",
			Items: []string{"thymeleaf", "jsp"},
		}

		_, templateEngine, err := promptTemplateEngine.Run()

		if err != nil {
			logger.Fatal(err.Error())
		}

		os.Args = append(os.Args, template, name, "-p", packageName, "-t", templateEngine)
		err = cmd.Execute()

		if err != nil {
			logger.Fatal(err.Error())
		}

		return
	}

	os.Args = append(os.Args, template, name, "-p", packageName)
	err = cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}

}
