package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize skills",
	Long:  `Synchronize skills from the central repository to local directories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Syncing skills...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
