package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/mpvm/internal/config"
	"github.com/takoeight0821/skills/mpvm/internal/logging"
	"github.com/takoeight0821/skills/mpvm/internal/multipass"
	"github.com/takoeight0821/skills/mpvm/pkg/version"
)

var (
	cfgFile string
	vmName  string
	cfg     *config.Config
	log     *logging.Logger
	client  *multipass.Client
)

var rootCmd = &cobra.Command{
	Use:   "mpvm",
	Short: "Multipass VM manager for coding agents",
	Long: `mpvm manages Multipass VMs configured for Claude Code and Gemini CLI.
Supports SSH Agent Forwarding for secure git commit signing.`,
	Version: version.Short(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override VM name if provided via flag
		if vmName != "" {
			cfg.VM.Name = vmName
		}

		log = logging.Default()
		client = multipass.NewClient()

		// Check if multipass is installed
		if !client.IsInstalled() {
			return fmt.Errorf("multipass is not installed. Please install it from https://multipass.run")
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&vmName, "vm-name", "n", "", "VM name (default: coding-agent)")

	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func getVMName() string {
	if cfg != nil && cfg.VM.Name != "" {
		return cfg.VM.Name
	}
	return "coding-agent"
}

func checkVMExists() error {
	exists, err := client.VMExists(getVMName())
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("VM '%s' does not exist. Run 'mpvm launch' to create it", getVMName())
	}
	return nil
}

func checkVMRunning() error {
	if err := checkVMExists(); err != nil {
		return err
	}

	running, err := client.VMRunning(getVMName())
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("VM '%s' is not running. Run 'mpvm start' to start it", getVMName())
	}
	return nil
}

func getVMTerm() string {
	term := os.Getenv("TERM")
	if term == "" {
		return "xterm-256color"
	}

	switch term {
	case "xterm-ghostty", "ghostty":
		return "xterm-256color"
	default:
		return term
	}
}
