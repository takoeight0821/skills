package cli

import (
	"github.com/spf13/cobra"
)

// vmCmd represents the vm command group for Multipass VM management
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage Multipass VMs",
	Long: `VM commands for managing Multipass-based coding agent VMs.
Supports launch, start, stop, delete, ssh, claude, gemini, exec, mount, status, and logs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Run parent's PersistentPreRunE first
		if rootCmd.PersistentPreRunE != nil {
			if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
		}
		// Initialize multipass client for vm commands
		return initVMClient()
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)
}
