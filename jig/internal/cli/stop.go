package cli

import (
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop running VM",
	Long:  `Stop a running Multipass VM.`,
	RunE:  runStop,
}

func init() {
	vmCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	vmName := getVMName()

	if err := checkVMExists(); err != nil {
		return err
	}

	running, err := client.VMRunning(vmName)
	if err != nil {
		return err
	}

	if !running {
		log.Info("VM '%s' is already stopped.", vmName)
		return nil
	}

	log.Info("Stopping VM '%s'...", vmName)
	if err := client.Stop(vmName); err != nil {
		return err
	}

	log.Success("VM '%s' stopped.", vmName)
	return nil
}
