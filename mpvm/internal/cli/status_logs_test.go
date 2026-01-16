package cli

import (
	"strings"
	"testing"

	"github.com/takoeight0821/skills/skills-cli/internal/config"
	"github.com/takoeight0821/skills/skills-cli/internal/logging"
	"github.com/takoeight0821/skills/skills-cli/internal/multipass"
)

func TestStatus_Running(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true

	// Capture output? runStatus prints to stdout.
	// We can't easily assert stdout without redirecting it, which is complex in parallel tests.
	// But we can ensure it doesn't error.
	err := runStatus(nil, []string{})
	if err != nil {
		t.Errorf("Status failed: %v", err)
	}
}

func TestStatus_NotExist(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	err := runStatus(nil, []string{})
	if err != nil {
		t.Errorf("Status should not fail if VM does not exist, just log info. Got: %v", err)
	}
}

func TestLogs(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true

	err := runLogs(nil, []string{})
	if err != nil {
		t.Errorf("Logs failed: %v", err)
	}

	// Verify execution logic for logs
	if len(mockClient.ExecCalls) == 0 {
		t.Error("Expected exec call for logs")
	}

	lastCall := mockClient.ExecCalls[len(mockClient.ExecCalls)-1]
	if !strings.Contains(lastCall, "tail") {
		t.Errorf("Expected tail command, got: %s", lastCall)
	}
}
