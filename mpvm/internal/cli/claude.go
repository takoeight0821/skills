package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Run Claude Code in VM",
	Long: `Run Claude Code inside the VM with SSH agent forwarding.
The current directory will be automatically mounted and used as the working directory.`,
	RunE:               runClaude,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(claudeCmd)
}

func runClaude(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	// Get VM IP
	ip, err := client.GetIP(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM IP: %w", err)
	}

	// Auto-mount current directory
	cwd, _ := os.Getwd()
	mountPoint := client.GetMountPoint(vmName, cwd)
	mounted, _ := client.IsMounted(vmName, cwd)
	if !mounted {
		log.Info("Mounting %s to %s", cwd, mountPoint)
		if err := client.Mount(cwd, vmName+":"+mountPoint); err != nil {
			return fmt.Errorf("failed to mount directory: %w", err)
		}
	}

	// Build SSH command
	term := getVMTerm()
	claudeArgs := []string{"claude"}
	claudeArgs = append(claudeArgs, args...)

	sshArgs := []string{
		"-A",
		"-t",
		"-o", "StrictHostKeyChecking=accept-new",
		fmt.Sprintf("ubuntu@%s", ip),
fmt.Sprintf("printf '\\e[?1004l'; cd %s && TERM=%s %s", mountPoint, term, joinArgs(claudeArgs)),
	}

	sshPath, _ := exec.LookPath("ssh")
	sshCmd := exec.Command(sshPath, sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		// Simple quoting for args with spaces
		if containsSpace(arg) {
			result += fmt.Sprintf("'%s'", arg)
		} else {
			result += arg
		}
	}
	return result
}

func containsSpace(s string) bool {
	for _, c := range s {
		if c == ' ' {
			return true
		}
	}
	return false
}
