package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage skills vm configuration.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long: `Create a default configuration file.
The configuration file is created at ~/.config/skills vm/config.toml`,
	RunE: runConfigInit,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	RunE:  runConfigPath,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE:  runConfigShow,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	configPath := config.ConfigPath()

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		log.Warn("Configuration file already exists at %s", configPath)
		return nil
	}

	// Create default config
	defaultCfg := config.DefaultConfig()

	// Prompt for git user info if not set
	fmt.Printf("Enter your git user name: ")
	fmt.Scanln(&defaultCfg.Git.UserName)

	fmt.Printf("Enter your git user email: ")
	fmt.Scanln(&defaultCfg.Git.UserEmail)

	// Save config
	if err := defaultCfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	log.Success("Configuration file created at %s", configPath)
	return nil
}

func runConfigPath(cmd *cobra.Command, args []string) error {
	fmt.Println(config.ConfigPath())
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Configuration file: %s\n\n", config.ConfigPath())
	fmt.Printf("[vm]\n")
	fmt.Printf("  name = %q\n", cfg.VM.Name)
	fmt.Printf("  cpus = %d\n", cfg.VM.CPUs)
	fmt.Printf("  memory = %q\n", cfg.VM.Memory)
	fmt.Printf("  disk = %q\n", cfg.VM.Disk)
	fmt.Printf("\n[docker]\n")
	fmt.Printf("  container_name = %q\n", cfg.Docker.ContainerName)
	fmt.Printf("  image_name = %q\n", cfg.Docker.ImageName)
	fmt.Printf("  cpus = %q\n", cfg.Docker.CPUs)
	fmt.Printf("  memory = %q\n", cfg.Docker.Memory)
	fmt.Printf("\n[git]\n")
	fmt.Printf("  user_name = %q\n", cfg.Git.UserName)
	fmt.Printf("  user_email = %q\n", cfg.Git.UserEmail)
	fmt.Printf("\n[ssh]\n")
	fmt.Printf("  signing_key = %q\n", cfg.SSH.SigningKey)
	fmt.Printf("\n[claude]\n")
	fmt.Printf("  marketplaces = %v\n", cfg.Claude.Marketplaces)

	return nil
}
