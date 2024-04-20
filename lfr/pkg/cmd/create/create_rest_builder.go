package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/lfr/pkg/cmd/exec"
	"github.com/lgdd/lfr-cli/lfr/pkg/generate/rb"
	"github.com/lgdd/lfr-cli/lfr/pkg/project"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createRestBuilder = &cobra.Command{
		Use:     "rest-builder NAME",
		Aliases: []string{"rb"},
		Args:    cobra.ExactArgs(1),
		Run:     generateRestBuilder,
	}
	// CodeGen holds the option to run the code generation of the rest builder
	CodeGen bool
)

func init() {
	createRestBuilder.Flags().BoolVarP(&CodeGen, "generate", "g", false, "executes code generation")
}

func generateRestBuilder(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	name := args[0]
	name = strcase.ToKebab(strings.ToLower(name))
	rb.Generate(liferayWorkspace, name)

	build := project.Maven

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		build = project.Gradle
	}

	if CodeGen {
		runCodeGen(liferayWorkspace, name, build)
	} else {
		printCodeGenSuggestion(liferayWorkspace, name, build)
	}
}

func runCodeGen(workspace, name, build string) {
	moduleImplPath := filepath.Join(workspace, "modules", name, name+"-impl")
	if err := os.Chdir(moduleImplPath); err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}

	fmt.Print("\nExecutes code generation using:\n\n")

	switch build {
	case project.Gradle:
		printutil.Info("lfr exec buildREST\n\n")
		exec.RunWrapperCmd([]string{"buildREST"})
	case project.Maven:
		printutil.Info("lfr exec rest-builder:build\n\n")
		exec.RunWrapperCmd([]string{"rest-builder:build"})
	}
}

func printCodeGenSuggestion(workspace, name, build string) {
	moduleImplPath := filepath.Join(workspace, "modules", name, name+"-impl")
	fmt.Println("\nTo execute code generation:")
	switch build {
	case project.Gradle:
		printutil.Info(fmt.Sprintf("cd %s && lfr exec buildREST\n", moduleImplPath))
	case project.Maven:
		printutil.Info(fmt.Sprintf("cd %s && lfr exec rest-builder:build\n", moduleImplPath))
	}
	fmt.Print("\n")
}
