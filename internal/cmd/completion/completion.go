// Package completion implements the completion subcommand, which generates
// shell completion scripts for Bash, Zsh, Fish, and PowerShell.
package completion

import (
	"os"

	"github.com/lgdd/lfr-cli/pkg/util/logger"
	"github.com/spf13/cobra"
)

// Cmd is the command 'completion' which generate a completion script
var Cmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(lfr completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ lfr completion bash > /etc/bash_completion.d/lfr
  # macOS:
  $ lfr completion bash > /usr/local/etc/bash_completion.d/lfr

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ lfr completion zsh > "${fpath[1]}/_lfr"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ lfr completion fish | source

  # To load completions for each session, execute once:
  $ lfr completion fish > ~/.config/fish/completions/lfr.fish

PowerShell:

  PS> lfr completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> lfr completion powershell > lfr.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		if err != nil {
			logger.Fatal(err.Error())
		}
	},
}
