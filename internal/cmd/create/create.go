package create

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/pkg/project"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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
	Cmd.AddCommand(createRestBuilder)
	Cmd.AddCommand(createDocker)
	Cmd.PersistentFlags().StringVarP(&project.PackageName, "package", "p", "org.acme", "base package name")
}

func promptCreateChoices(cmd *cobra.Command, args []string) {
	promptTemplate := promptui.Select{
		Label: "Choose a template",
		Items: []string{"client-extension", "api", "command", "docker", "mvc-portlet", "rest-builder", "service-builder", "spring-mvc-portlet", "workspace"},
	}

	_, template, err := promptTemplate.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if template == "docker" {
		promptJavaVersion := promptui.Select{
			Label: "Choose a Java version",
			Items: []string{"8", "11"},
		}
		_, javaVersion, err := promptJavaVersion.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		promptBuildOption := promptui.Select{
			Label: "Choose a build option",
			Items: []string{"Multi-stage", "Single stage"},
		}

		_, buildOption, err := promptBuildOption.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		if buildOption == "Multi-stage" {
			os.Args = append(os.Args, template, "-j", javaVersion, "-m")
		} else {
			os.Args = append(os.Args, template, "-j", javaVersion)
		}

		err = cmd.Execute()

		if err != nil {
			printutil.Danger(err.Error())
			os.Exit(1)
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
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	workspacePackage, _ := project.GetGroupId()
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
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if template == "spring-mvc-portlet" {
		promptTemplateEngine := promptui.Select{
			Label: "Choose a template engine",
			Items: []string{"thymeleaf", "jsp"},
		}

		_, templateEngine, err := promptTemplateEngine.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		os.Args = append(os.Args, template, name, "-p", packageName, "-t", templateEngine)
		err = cmd.Execute()

		if err != nil {
			printutil.Danger(err.Error())
			os.Exit(1)
		}

		return
	}

	os.Args = append(os.Args, template, name, "-p", packageName)
	err = cmd.Execute()

	if err != nil {
		printutil.Danger(err.Error())
		os.Exit(1)
	}

}
