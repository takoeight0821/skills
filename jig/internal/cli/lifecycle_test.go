package cli

import (
	"testing"

	"github.com/takoeight0821/skills/jig/internal/config"
	"github.com/takoeight0821/skills/jig/internal/logging"
	"github.com/takoeight0821/skills/jig/internal/multipass"
)

func TestStop_RunningVM(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true // running

	err := runStop(nil, []string{})
	if err != nil {
		t.Errorf("Stop failed: %v", err)
	}

	running, _ := mockClient.VMRunning("coding-agent")
	if running {
		t.Error("VM should be stopped")
	}
}

func TestStop_StoppedVM(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = false // stopped

	err := runStop(nil, []string{})
	if err != nil {
		t.Errorf("Stop failed: %v", err)
	}

	// Should remain stopped
	running, _ := mockClient.VMRunning("coding-agent")
	if running {
		t.Error("VM should be stopped")
	}
}

func TestDelete_Force(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	// Set force flag
	deleteForce = true
	defer func() { deleteForce = false }()

	mockClient.VMs["coding-agent"] = false

	err := runDelete(nil, []string{})
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	exists, _ := mockClient.VMExists("coding-agent")
	if exists {
		t.Error("VM should be deleted")
	}

	if len(mockClient.Deleted) != 1 {
		t.Errorf("Expected 1 VM deleted, got %d", len(mockClient.Deleted))
	}
}

func TestDelete_NotExist(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()
	deleteForce = true
	defer func() { deleteForce = false }()

	err := runDelete(nil, []string{}) // coding-agent doesn't exist in mock
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Should log warning but not fail
}
