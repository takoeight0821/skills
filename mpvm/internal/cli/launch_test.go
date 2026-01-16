package cli

import (
	"testing"

	"github.com/takoeight0821/skills/skills-cli/internal/config"
	"github.com/takoeight0821/skills/skills-cli/internal/logging"
	"github.com/takoeight0821/skills/skills-cli/internal/multipass"
)

func TestLaunch_NewVM(t *testing.T) {
	// Setup
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	// Override config to avoid side effects (though these are mostly read)
	cfg.SSH.SigningKey = "mock_key.pub" // typo handling: Config might have SigningKey or similar

	// Execute
	// We call runLaunch directly or via command?
	// runLaunch expects *cobra.Command which we can mock or pass nil if unused
	err := runLaunch(nil, []string{})

	// Verify
	if err != nil {
		t.Errorf("Launch failed: %v", err)
	}

	if len(mockClient.Launched) != 1 {
		t.Errorf("Expected 1 VM launched, got %d", len(mockClient.Launched))
	}

	if mockClient.Launched[0] != "coding-agent" {
		t.Errorf("Expected VM name 'coding-agent', got '%s'", mockClient.Launched[0])
	}
}

func TestLaunch_ExistingStoppedVM(t *testing.T) {
	// Setup
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	// Pre-create VM in stopped state
	mockClient.VMs["coding-agent"] = false // exists, not running

	// Execute
	err := runLaunch(nil, []string{})

	// Verify
	if err != nil {
		t.Errorf("Launch failed: %v", err)
	}

	if len(mockClient.Launched) != 0 {
		t.Errorf("Expected 0 VMs launched (should start existing), got %d", len(mockClient.Launched))
	}

	running, _ := mockClient.VMRunning("coding-agent")
	if !running {
		t.Error("VM should be running after launch")
	}
}

func TestLaunch_ExistingRunningVM(t *testing.T) {
	// Setup
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	// Pre-create VM in running state
	mockClient.VMs["coding-agent"] = true // exists, running

	// Execute
	err := runLaunch(nil, []string{})

	// Verify
	if err != nil {
		t.Errorf("Launch failed: %v", err)
	}

	if len(mockClient.Launched) != 0 {
		t.Errorf("Expected 0 VMs launched, got %d", len(mockClient.Launched))
	}
}
