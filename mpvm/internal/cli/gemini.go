package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var geminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Run Gemini CLI in VM",
	Long: `Run Gemini CLI inside the VM with SSH agent forwarding.
The current directory will be automatically mounted and used as the working directory.`,
	RunE:               runGemini,
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(geminiCmd)
}

func runGemini(cmd *cobra.Command, args []string) error {
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
	geminiArgs := []string{"gemini"}
	geminiArgs = append(geminiArgs, args...)

	sshArgs := []string{
		"-A",
		"-t",
		"-o", "StrictHostKeyChecking=accept-new",
		fmt.Sprintf("ubuntu@%s", ip),
		fmt.Sprintf("printf '\\e[?1004l'; cd %s && TERM=%s %s", mountPoint, term, joinArgs(geminiArgs)),
	}

	sshPath, _ := exec.LookPath("ssh")
	sshCmd := exec.Command(sshPath, sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}
