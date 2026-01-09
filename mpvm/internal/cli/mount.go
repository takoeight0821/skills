package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var mountCmd = &cobra.Command{
	Use:   "mount [path]",
	Short: "Mount directory to VM",
	Long: `Mount a local directory to the VM.
If no path is specified, the current directory is mounted.

The mount point in the VM is automatically determined based on the source path.`,
	RunE: runMount,
}

var umountCmd = &cobra.Command{
	Use:     "umount [path]",
	Aliases: []string{"unmount"},
	Short:   "Unmount directory from VM",
	Long:    `Unmount a previously mounted directory from the VM.`,
	RunE:    runUmount,
}

func init() {
	rootCmd.AddCommand(mountCmd)
	rootCmd.AddCommand(umountCmd)
}

func runMount(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	// Determine source path
	var sourcePath string
	if len(args) > 0 {
		sourcePath = args[0]
	} else {
		var err error
		sourcePath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	// Check if already mounted
	mounted, _ := client.IsMounted(vmName, absPath)
	if mounted {
		log.Info("'%s' is already mounted.", absPath)
		return nil
	}

	// Determine mount point
	mountPoint := client.GetMountPoint(vmName, absPath)

	log.Info("Mounting '%s' to '%s:%s'...", absPath, vmName, mountPoint)
	if err := client.Mount(absPath, vmName+":"+mountPoint); err != nil {
		return fmt.Errorf("failed to mount: %w", err)
	}

	log.Success("Mounted successfully.")
	return nil
}

func runUmount(cmd *cobra.Command, args []string) error {
	if err := checkVMRunning(); err != nil {
		return err
	}

	vmName := getVMName()

	// Determine source path
	var sourcePath string
	if len(args) > 0 {
		sourcePath = args[0]
	} else {
		var err error
		sourcePath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Get mount point
	mountPoint := client.GetMountPoint(vmName, absPath)
	target := vmName + ":" + mountPoint

	log.Info("Unmounting '%s'...", target)
	if err := client.Umount(target); err != nil {
		return fmt.Errorf("failed to unmount: %w", err)
	}

	log.Success("Unmounted successfully.")
	return nil
}
