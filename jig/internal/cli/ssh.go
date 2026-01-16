package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into VM with agent forwarding",
	Long: `Connect to the VM via SSH with agent forwarding enabled.
This allows you to use your local SSH keys for git operations inside the VM.`,
	RunE: runSSH,
}

func init() {
	vmCmd.AddCommand(sshCmd)
}

func runSSH(cmd *cobra.Command, args []string) error {
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
	if cwd != "" {
		mountPoint := client.GetMountPoint(vmName, cwd)
		mounted, _ := client.IsMounted(vmName, cwd)
		if !mounted {
			log.Info("Mounting %s to %s", cwd, mountPoint)
			if err := client.Mount(cwd, vmName+":"+mountPoint); err != nil {
				log.Warn("Failed to mount directory: %v", err)
			}
		}
	}

	// Get mount point for working directory
	mountPoint := client.GetMountPoint(vmName, cwd)

	// Build SSH command with agent forwarding
	sshArgs := []string{
		"-A", // Agent forwarding
		"-t", // Force pseudo-terminal allocation
		"-o", "StrictHostKeyChecking=accept-new",
		fmt.Sprintf("ubuntu@%s", ip),
	}

	// Add custom TERM/COLORTERM and disable focus reporting
	term := getVMTerm()
	colorterm := getColorTerm()

	// Add any additional arguments or default shell
	if len(args) > 0 {
		sshArgs = append(sshArgs, args...)
	} else {
		// Disable focus reporting (\e[?1004l), set TERM/COLORTERM, cd to mount point, and start bash
		sshArgs = append(sshArgs, fmt.Sprintf("printf '\\e[?1004l'; export TERM=%s COLORTERM=%s; cd '%s' && exec bash -l", term, colorterm, mountPoint))
	}

	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh not found: %w", err)
	}

	// Execute SSH
	sshCmd := exec.Command(sshPath, sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}

func autoMountWorkDir(vmName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	mountPoint := filepath.Join("/mnt", sanitizePath(cwd))

	mounted, _ := client.IsMounted(vmName, cwd)
	if !mounted {
		if err := client.Mount(cwd, vmName+":"+mountPoint); err != nil {
			return "", err
		}
	}

	return mountPoint, nil
}

func sanitizePath(path string) string {
	// Convert path to mount-safe name
	// e.g., /home/user/project -> home-user-project
	clean := filepath.Clean(path)
	clean = filepath.ToSlash(clean)
	if len(clean) > 0 && clean[0] == '/' {
		clean = clean[1:]
	}
	return filepath.Join(filepath.SplitList(clean)...)
}
