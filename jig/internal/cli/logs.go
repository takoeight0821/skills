package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	logsLines int
	logsFollow bool
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show cloud-init logs",
	Long:  `Display cloud-init logs from the VM.`,
	RunE:  runLogs,
}

func init() {
	logsCmd.Flags().IntVarP(&logsLines, "lines", "l", 100, "Number of lines to show")
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	vmCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	if logsFollow {
		// Use interactive exec for streaming output
		tailCmd := fmt.Sprintf("tail -f /var/log/cloud-init-output.log 2>/dev/null || tail -f /var/log/cloud-init.log")
		return client.ExecInteractive(vmName, "bash", "-c", tailCmd)
	}

	// Non-follow mode: get output and print
	tailCmd := fmt.Sprintf("tail -n %d /var/log/cloud-init-output.log", logsLines)
	output, err := client.ExecOutput(vmName, "bash", "-c", tailCmd)
	if err != nil {
		// Try alternative log location
		output, err = client.ExecOutput(vmName, "bash", "-c",
			fmt.Sprintf("tail -n %d /var/log/cloud-init.log 2>/dev/null || echo 'No cloud-init logs found'", logsLines))
		if err != nil {
			return fmt.Errorf("failed to get logs: %w", err)
		}
	}

	fmt.Print(output)
	return nil
}
