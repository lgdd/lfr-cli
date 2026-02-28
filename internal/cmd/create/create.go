// Package create implements the create subcommand and its subcommands for
// scaffolding Liferay workspaces, modules, client extensions, and Docker files.
package create

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/internal/prompt"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

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
				runPrompt(cmd)
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
	conf.Init()
	defaultPackage := viper.GetString(conf.ModulePackage)
	Cmd.PersistentFlags().StringVarP(&metadata.PackageName, "package", "p", defaultPackage, "base package name")
}

func runPrompt(cmd *cobra.Command) {
	var template, sample, name string
	packageName := viper.GetString(conf.ModulePackage)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a template:").
				Options(
					huh.NewOption("Liferay Workspace", "workspace"),
					huh.NewOption("Client Extension", "client-extension"),
					huh.NewOption("Docker (Dockerfile & docker-compose.yml)", "docker"),
					huh.NewOption("OSGi - API Module", "api"),
					huh.NewOption("OSGi - Gogo Shell Command", "command"),
					huh.NewOption("OSGi - MVC Portlet", "mvc-portlet"),
					huh.NewOption("OSGi - Spring MVC Portlet", "spring-mvc-portlet"),
					huh.NewOption("OSGi - REST Builder", "rest-builder"),
					huh.NewOption("OSGi - Service Builder", "service-builder"),
				).
				Value(&template),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err := form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	if template == "workspace" {
		prompt.ForWorkspace(cmd, &name)
		return
	}

	if template == "client-extension" {
		prompt.ForClientExtension(cmd, &sample, &name)
		return
	}

	if template == "docker" {
		prompt.ForDocker(cmd)
		return
	}

	if template == "spring-mvc-portlet" {
		prompt.ForSpring(cmd, &packageName, &name)
		return
	}

	form = huh.NewForm(
		huh.NewGroup(
			prompt.NewInputPackageName(&packageName),
			prompt.NewInputName(&name),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err = form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	os.Args = append(os.Args, template, name, "-p", packageName)
	err = cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}
}
