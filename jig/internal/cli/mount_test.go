package cli

import (
	"path/filepath"
	"testing"

	"github.com/takoeight0821/skills/jig/internal/config"
	"github.com/takoeight0821/skills/jig/internal/logging"
	"github.com/takoeight0821/skills/jig/internal/multipass"
)

func TestMount_CurrentDir(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true

	// Test mount current directory
	err := runMount(nil, []string{})
	if err != nil {
		t.Errorf("Mount failed: %v", err)
	}

	// Since we mock the client, we can't easily check if Mount was called with exact args
	// unless we inspect the mock state (which we didn't fully implement for Mount args in MockClient).
	// But we can ensure no error was returned.
}

func TestMount_SpecificDir(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true

	// Test mount specific directory (must exist)
	// We use the current directory as a valid path
	absPath, _ := filepath.Abs(".")

	err := runMount(nil, []string{absPath})
	if err != nil {
		t.Errorf("Mount failed: %v", err)
	}
}

func TestUmount_SpecificDir(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true

	absPath, _ := filepath.Abs(".")

	err := runUmount(nil, []string{absPath})
	if err != nil {
		t.Errorf("Umount failed: %v", err)
	}
}
