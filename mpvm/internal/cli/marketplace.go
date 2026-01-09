package cli

import (
	"github.com/spf13/cobra"
)

var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Manage Claude plugin marketplaces",
	Long:  `Manage Claude Code plugin marketplaces in the VM.`,
}

var marketplaceAddCmd = &cobra.Command{
	Use:   "add [marketplace...]",
	Short: "Add plugin marketplaces",
	Long: `Add Claude Code plugin marketplaces to the VM.
If no marketplace is specified, adds all marketplaces from config.`,
	RunE: runMarketplaceAdd,
}

func init() {
	marketplaceCmd.AddCommand(marketplaceAddCmd)
	rootCmd.AddCommand(marketplaceCmd)
}

func runMarketplaceAdd(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	// Determine which marketplaces to add
	var marketplaces []string
	if len(args) > 0 {
		marketplaces = args
	} else {
		marketplaces = cfg.Claude.Marketplaces
	}

	if len(marketplaces) == 0 {
		log.Info("No marketplaces to add. Specify marketplace names or configure them in config.toml")
		return nil
	}

	// Add each marketplace
	for _, marketplace := range marketplaces {
		log.Info("Adding Claude plugin marketplace: %s", marketplace)
		if err := client.Exec(vmName, "claude", "-p", "plugin", "marketplace", "add", marketplace); err != nil {
			log.Warn("Failed to add marketplace %s: %v", marketplace, err)
		} else {
			log.Success("Added marketplace: %s", marketplace)
		}
	}

	return nil
}
