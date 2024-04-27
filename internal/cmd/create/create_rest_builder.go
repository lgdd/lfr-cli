package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lgdd/lfr-cli/internal/cmd/exec"
	"github.com/lgdd/lfr-cli/pkg/metadata"
	"github.com/lgdd/lfr-cli/pkg/scaffold"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/printutil"
	"github.com/spf13/cobra"
)

var (
	createRESTBuilder = &cobra.Command{
		Use:     "rest-builder NAME",
		Aliases: []string{"rb"},
		Args:    cobra.ExactArgs(1),
		Run:     generateRESTBuilder,
	}
	// CodeGen holds the option to run the code generation of the rest builder
	CodeGen bool
)

func init() {
	createRESTBuilder.Flags().BoolVarP(&CodeGen, "generate", "g", false, "executes code generation")
}

func generateRESTBuilder(cmd *cobra.Command, args []string) {
	liferayWorkspace, err := fileutil.GetLiferayWorkspacePath()
	if err != nil {
		printutil.Danger(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
	}
	name := args[0]
	name = strcase.ToKebab(strings.ToLower(name))
	scaffold.CreateModuleRESTBuilder(liferayWorkspace, name)

	build := metadata.Maven

	if fileutil.IsGradleWorkspace(liferayWorkspace) {
		build = metadata.Gradle
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
	case metadata.Gradle:
		printutil.Info("lfr exec buildREST\n\n")
		exec.RunWrapperCmd([]string{"buildREST"})
	case metadata.Maven:
		printutil.Info("lfr exec rest-builder:build\n\n")
		exec.RunWrapperCmd([]string{"rest-builder:build"})
	}
}

func printCodeGenSuggestion(workspace, name, build string) {
	moduleImplPath := filepath.Join(workspace, "modules", name, name+"-impl")
	fmt.Println("\nTo execute code generation:")
	switch build {
	case metadata.Gradle:
		printutil.Info(fmt.Sprintf("cd %s && lfr exec buildREST\n", moduleImplPath))
	case metadata.Maven:
		printutil.Info(fmt.Sprintf("cd %s && lfr exec rest-builder:build\n", moduleImplPath))
	}
	fmt.Print("\n")
}
