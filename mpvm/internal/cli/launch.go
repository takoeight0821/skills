package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/cloudinit"
	"github.com/takoeight0821/skills/skills-cli/internal/git"
	"github.com/takoeight0821/skills/skills-cli/internal/multipass"
)

var launchCmd = &cobra.Command{
	Use:     "launch",
	Aliases: []string{"create"},
	Short:   "Create and start VM with cloud-init",
	Long: `Create a new Multipass VM configured for coding agents.
If the VM already exists, it will be started instead.`,
	RunE: runLaunch,
}

func init() {
	vmCmd.AddCommand(launchCmd)
}

func runLaunch(cmd *cobra.Command, args []string) error {
	vmName := getVMName()

	// Check if VM exists
	exists, err := client.VMExists(vmName)
	if err != nil {
		return err
	}

	if exists {
		log.Warn("VM '%s' already exists.", vmName)
		running, _ := client.VMRunning(vmName)
		if !running {
			log.Info("Starting existing VM...")
			if err := client.Start(vmName); err != nil {
				return err
			}
		} else {
			log.Info("VM is already running.")
		}
	} else {
		log.Info("Creating VM '%s'...", vmName)
		log.Info("  CPUs: %d", cfg.VM.CPUs)
		log.Info("  Memory: %s", cfg.VM.Memory)
		log.Info("  Disk: %s", cfg.VM.Disk)

		// Write cloud-init to temp file
		cloudInitPath, cleanup, err := cloudinit.WriteCloudInitTempFile()
		if err != nil {
			return fmt.Errorf("failed to create cloud-init file: %w", err)
		}
		defer cleanup()

		// Launch VM
		if err := client.Launch(multipass.LaunchOptions{
			Name:      vmName,
			CPUs:      cfg.VM.CPUs,
			Memory:    cfg.VM.Memory,
			Disk:      cfg.VM.Disk,
			CloudInit: cloudInitPath,
		}); err != nil {
			return fmt.Errorf("failed to launch VM: %w", err)
		}

		// Wait for cloud-init
		log.Info("Waiting for cloud-init to complete (this may take a few minutes)...")
		if err := client.WaitForCloudInit(vmName, 600, log.Info); err != nil {
			log.Warn("Cloud-init may not have completed: %v", err)
		} else {
			log.Success("Cloud-init completed!")
		}
	}

	// Copy SSH public key for authentication
	if err := copySSHPublicKey(vmName); err != nil {
		log.Warn("SSH key setup failed: %v", err)
	}

	// Configure git
	if cfg.Git.UserName != "" || cfg.Git.UserEmail != "" {
		log.Info("Configuring git...")
		if err := git.ConfigureInVM(client, cfg, vmName); err != nil {
			log.Warn("Git configuration failed: %v", err)
		}
	}

	// Setup Claude
	if err := setupClaude(vmName); err != nil {
		log.Warn("Claude setup failed: %v", err)
	}

	// Get IP and show success message
	ip, _ := client.GetIP(vmName)

	fmt.Println()
	log.Success("VM '%s' is ready!", vmName)
	if ip != "" {
		log.Success("IP Address: %s", ip)
	}

	fmt.Println()
	fmt.Println("Connect with:")
	if ip != "" {
		fmt.Printf("  ssh -A ubuntu@%s\n", ip)
	}
	fmt.Println()
	fmt.Println("Or use:")
	fmt.Println("  skills vm ssh      # Interactive shell")
	fmt.Println("  skills vm claude   # Run Claude Code")
	fmt.Println("  skills vm gemini   # Run Gemini CLI")

	return nil
}

func setupClaude(vmName string) error {
	// Create settings file in VM
	settingsContent, err := cloudinit.GetVMSettings()
	if err != nil {
		return err
	}

	// Write settings to VM (non-interactive)
	settingsPath := "/home/ubuntu/.claude/settings.json"
	cmd := fmt.Sprintf("mkdir -p /home/ubuntu/.claude && cat > %s << 'EOF'\n%s\nEOF", settingsPath, string(settingsContent))
	if err := client.Exec(vmName, "bash", "-c", cmd); err != nil {
		return fmt.Errorf("failed to write Claude settings: %w", err)
	}

	log.Info("Claude settings configured. Run 'skills vm claude' to authenticate.")
	return nil
}

func copySSHPublicKey(vmName string) error {
	// Auto-detect user's SSH public key
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keyPaths := []string{
		filepath.Join(home, ".ssh", "id_ed25519.pub"),
		filepath.Join(home, ".ssh", "id_rsa.pub"),
	}

	var pubKey []byte
	var usedPath string
	for _, path := range keyPaths {
		if data, err := os.ReadFile(path); err == nil {
			pubKey = data
			usedPath = path
			break
		}
	}

	if pubKey == nil {
		return fmt.Errorf("no SSH public key found in ~/.ssh/")
	}

	log.Info("Copying SSH key: %s", usedPath)

	// Add to VM's authorized_keys
	cmd := fmt.Sprintf(`mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys`,
		strings.TrimSpace(string(pubKey)))
	if err := client.Exec(vmName, "bash", "-c", cmd); err != nil {
		return fmt.Errorf("failed to copy SSH key: %w", err)
	}

	return nil
}

func runMultipassInteractive(vmName string, command ...string) error {
	args := append([]string{"exec", vmName, "--"}, command...)

	mpCmd := exec.Command("multipass", args...)
	mpCmd.Stdin = os.Stdin
	mpCmd.Stdout = os.Stdout
	mpCmd.Stderr = os.Stderr

	return mpCmd.Run()
}
