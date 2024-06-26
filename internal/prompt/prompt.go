package prompt

import (
	"errors"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/lgdd/lfr-cli/internal/conf"
	"github.com/lgdd/lfr-cli/pkg/util/helper"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/lgdd/lfr-cli/pkg/util/procutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ForDocker(cmd *cobra.Command) {
	var dockerBuildOption string
	javaVersion, _, err := procutil.GetCurrentJavaVersion()

	if err != nil {
		javaVersion = "11"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Java version:").
				Options(
					huh.NewOption("JDK 8", "8"),
					huh.NewOption("JDK 11", "11"),
				).
				Value(&javaVersion),

			huh.NewSelect[string]().
				Title("Choose a build option:").
				Options(
					huh.NewOption("Multi-stage", "-m"),
					huh.NewOption("Single stage", ""),
				).
				Value(&dockerBuildOption),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err = form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	os.Args = append(os.Args, "docker", "-j", javaVersion)

	if dockerBuildOption != "" {
		os.Args = append(os.Args, dockerBuildOption)
	}

	err = cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func ForSpring(cmd *cobra.Command, packageName, name *string) {
	var templateEngine string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a template engine:").
				Options(
					huh.NewOption("Thymeleaf", "thymeleaf"),
					huh.NewOption("JSP", "jsp"),
				).
				Value(&templateEngine),
			NewInputPackageName(packageName),
			NewInputName(name),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err := form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	os.Args = append(os.Args, "spring-mvc-portlet", *name, "-p", *packageName, "-t", templateEngine)
	err = cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func ForWorkspace(cmd *cobra.Command, name *string) {
	ForName(name)

	os.Args = append(os.Args, "workspace", *name)
	err := cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func ForName(name *string) {
	form := huh.NewForm(
		huh.NewGroup(
			NewInputName(name),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err := form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func ForClientExtension(cmd *cobra.Command, sample, name *string) {
	samples := helper.GetClientExtensionSampleNames()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Chosse a sample:").
				Options(huh.NewOptions(samples...)...).
				Value(sample),
			NewInputName(name),
		),
	).WithAccessible(viper.GetBool(conf.OutputAccessible))

	if conf.NoColor {
		form.WithTheme(huh.ThemeBase())
	}

	err := form.Run()

	if err != nil {
		logger.Fatal(err.Error())
	}

	os.Args = append(os.Args, "cx", *sample, *name)
	err = cmd.Execute()

	if err != nil {
		logger.Fatal(err.Error())
	}
}

func NewInputName(name *string) *huh.Input {
	return huh.NewInput().
		Title("Choose a name:").
		Value(name).
		Validate(isNotEmpty)
}

func NewInputPackageName(packageName *string) *huh.Input {
	return huh.NewInput().
		Title("Choose a package name:").
		Value(packageName).
		Validate(isNotEmpty)
}

func isNotEmpty(input string) error {
	if len(input) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}
