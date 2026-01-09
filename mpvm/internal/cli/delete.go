package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"destroy", "rm"},
	Short:   "Delete VM",
	Long:    `Delete a Multipass VM. Requires confirmation unless --force is used.`,
	RunE:    runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	vmName := getVMName()

	exists, err := client.VMExists(vmName)
	if err != nil {
		return err
	}

	if !exists {
		log.Warn("VM '%s' does not exist.", vmName)
		return nil
	}

	// Confirm deletion
	if !deleteForce {
		fmt.Printf("Are you sure you want to delete VM '%s'? [y/N] ", vmName)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			log.Info("Deletion cancelled.")
			return nil
		}
	}

	log.Info("Deleting VM '%s'...", vmName)
	if err := client.Delete(vmName, true); err != nil {
		return err
	}

	log.Success("VM '%s' deleted and purged.", vmName)
	return nil
}
