package cli

import (
	"testing"

	"github.com/takoeight0821/skills/mpvm/internal/config"
	"github.com/takoeight0821/skills/mpvm/internal/logging"
	"github.com/takoeight0821/skills/mpvm/internal/multipass"
)

func TestStart_StoppedVM(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = false // stopped

	err := runStart(nil, []string{})
	if err != nil {
		t.Errorf("Start failed: %v", err)
	}

	running, _ := mockClient.VMRunning("coding-agent")
	if !running {
		t.Error("VM should be running")
	}
}

func TestStart_RunningVM(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true // running

	err := runStart(nil, []string{})
	if err != nil {
		t.Errorf("Start failed: %v", err)
	}

	// Should remain running
	running, _ := mockClient.VMRunning("coding-agent")
	if !running {
		t.Error("VM should be running")
	}
}

func TestStart_NotExist(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	err := runStart(nil, []string{})
	if err == nil {
		t.Error("Start should fail if VM does not exist")
	}
}
