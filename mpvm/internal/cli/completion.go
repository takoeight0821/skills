package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for skills vm.

To load completions:

Bash:
  $ source <(skills vm completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ skills vm completion bash > /etc/bash_completion.d/skills vm
  # macOS:
  $ skills vm completion bash > $(brew --prefix)/etc/bash_completion.d/skills vm

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ skills vm completion zsh > "${fpath[1]}/_skills vm"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ skills vm completion fish | source
  # To load completions for each session, execute once:
  $ skills vm completion fish > ~/.config/fish/completions/skills vm.fish

PowerShell:
  PS> skills vm completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> skills vm completion powershell > skills vm.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
