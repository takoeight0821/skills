package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/takoeight0821/skills/skills-cli/internal/config"
	"github.com/takoeight0821/skills/skills-cli/internal/multipass"
)

func ConfigureInVM(client multipass.Client, cfg *config.Config, vmName string) error {
	// Set git user name
	if cfg.Git.UserName != "" {
		if err := client.Exec(vmName, "git", "config", "--global", "user.name", cfg.Git.UserName); err != nil {
			return fmt.Errorf("failed to set git user.name: %w", err)
		}
	}

	// Set git user email
	if cfg.Git.UserEmail != "" {
		if err := client.Exec(vmName, "git", "config", "--global", "user.email", cfg.Git.UserEmail); err != nil {
			return fmt.Errorf("failed to set git user.email: %w", err)
		}
	}

	// Setup SSH signing key
	if err := setupSSHSigning(client, cfg, vmName); err != nil {
		return err
	}

	return nil
}

func setupSSHSigning(client multipass.Client, cfg *config.Config, vmName string) error {
	signingKeyPath := config.ExpandPath(cfg.SSH.SigningKey)
	if signingKeyPath == "" {
		return nil
	}

	// Check if key exists locally
	if _, err := os.Stat(signingKeyPath); err != nil {
		return fmt.Errorf("SSH signing key not found: %s", signingKeyPath)
	}

	// Read the public key
	pubKey, err := os.ReadFile(signingKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH signing key: %w", err)
	}

	// Copy public key to VM
	tmpFile, err := os.CreateTemp("", "signing-key-*.pub")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pubKey); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	vmKeyPath := "/home/ubuntu/.ssh/signing_key.pub"
	if err := client.Transfer(tmpFile.Name(), vmName+":"+vmKeyPath); err != nil {
		return fmt.Errorf("failed to transfer signing key: %w", err)
	}

	// Set git signing key
	if err := client.Exec(vmName, "git", "config", "--global", "user.signingkey", vmKeyPath); err != nil {
		return fmt.Errorf("failed to set git signing key: %w", err)
	}

	// Setup allowed_signers
	allowedSignersPath := "/home/ubuntu/.ssh/allowed_signers"
	allowedSignersContent := fmt.Sprintf("%s %s", cfg.Git.UserEmail, string(pubKey))

	// Create allowed_signers file
	if err := client.Exec(vmName, "bash", "-c",
		fmt.Sprintf("echo '%s' > %s", allowedSignersContent, allowedSignersPath)); err != nil {
		return fmt.Errorf("failed to create allowed_signers: %w", err)
	}

	// Configure git to use allowed_signers
	if err := client.Exec(vmName, "git", "config", "--global", "gpg.ssh.allowedSignersFile", allowedSignersPath); err != nil {
		return fmt.Errorf("failed to set allowed_signers file: %w", err)
	}

	return nil
}

func GetSSHSigningKeyPath(cfg *config.Config) string {
	keyPath := config.ExpandPath(cfg.SSH.SigningKey)
	if keyPath == "" {
		home, _ := os.UserHomeDir()
		keyPath = filepath.Join(home, ".ssh", "id_ed25519.pub")
	}
	return keyPath
}
