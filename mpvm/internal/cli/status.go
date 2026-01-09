package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show VM status",
	Long:  `Display detailed information about the VM.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	vmName := getVMName()

	exists, err := client.VMExists(vmName)
	if err != nil {
		return err
	}

	if !exists {
		log.Info("VM '%s' does not exist.", vmName)
		return nil
	}

	info, err := client.Info(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM info: %w", err)
	}

	fmt.Println(info)
	return nil
}
