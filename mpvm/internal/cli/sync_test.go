package cli

import (
	"testing"

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
