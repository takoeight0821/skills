package cli

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start stopped VM",
	Long:  `Start a stopped Multipass VM.`,
	RunE:  runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	vmName := getVMName()

	if err := checkVMExists(); err != nil {
		return err
	}

	running, err := client.VMRunning(vmName)
	if err != nil {
		return err
	}

	if running {
		log.Info("VM '%s' is already running.", vmName)
		return nil
	}

	log.Info("Starting VM '%s'...", vmName)
	if err := client.Start(vmName); err != nil {
		return err
	}

	log.Success("VM '%s' started.", vmName)
	return nil
}
