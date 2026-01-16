package cli

import (
	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/jig/internal/git"
)

var configureGitCmd = &cobra.Command{
	Use:   "configure-git",
	Short: "Configure git in VM",
	Long: `Re-configure git settings in the VM.
This sets up git user name, email, and SSH signing key.`,
	RunE: runConfigureGit,
}

func init() {
	vmCmd.AddCommand(configureGitCmd)
}

func runConfigureGit(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	log.Info("Configuring git in VM '%s'...", vmName)

	if err := git.ConfigureInVM(client, cfg, vmName); err != nil {
		return err
	}

	log.Success("Git configured successfully.")
	return nil
}
