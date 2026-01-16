package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/config"
	"github.com/takoeight0821/skills/skills-cli/internal/logging"
)

func TestSyncCommandExists(t *testing.T) {
	// Setup
	cfg = config.DefaultConfig()
	log = logging.Default()

	// Find sync command in rootCmd
	var found bool
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "sync" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'sync' command to be registered in rootCmd")
	}
}

func TestSyncFlags(t *testing.T) {
	// Find sync command
	var syncCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "sync" {
			syncCmd = cmd
			break
		}
	}

	if syncCmd == nil {
		t.Fatal("sync command not found")
	}

	// Check flags
	flags := []string{"dry-run", "force", "global", "project", "source"}
	for _, f := range flags {
		if syncCmd.Flags().Lookup(f) == nil {
			t.Errorf("Expected flag --%s to be defined", f)
		}
	}
}
