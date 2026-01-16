package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [command...]",
	Short: "Execute command in VM",
	Long: `Execute an arbitrary command inside the VM.
Use -- to separate skills vm flags from the command to execute.

Example:
  skills vm exec -- ls -la
  skills vm exec -- bash -c "echo hello"`,
	RunE:               runExec,
	DisableFlagParsing: true,
}

func init() {
	vmCmd.AddCommand(execCmd)
}

func runExec(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	// Remove leading -- if present
	if args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	vmName := getVMName()

	// Get VM IP for SSH (to support agent forwarding)
	ip, err := client.GetIP(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM IP: %w", err)
	}

	// Build SSH command with agent forwarding
	sshArgs := []string{
		"-A",
		"-o", "StrictHostKeyChecking=accept-new",
		fmt.Sprintf("ubuntu@%s", ip),
	}
	sshArgs = append(sshArgs, args...)

	sshPath, _ := exec.LookPath("ssh")
	sshCmd := exec.Command(sshPath, sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}
