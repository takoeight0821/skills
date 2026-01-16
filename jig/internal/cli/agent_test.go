package cli

import (
	"testing"

	"github.com/takoeight0821/skills/jig/internal/config"
	"github.com/takoeight0821/skills/jig/internal/logging"
	"github.com/takoeight0821/skills/jig/internal/multipass"
)

func TestClaude_Run(t *testing.T) {
	// similar limitation to exec: calls ssh via os/exec
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true
	mockClient.IPs["coding-agent"] = "1.2.3.4"

	// Just validity check
	// We can't actually run it because it tries to exec ssh
}

func TestJoinArgs(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"hello", "world"}, "hello world"},
		{[]string{"hello", "world space"}, "hello 'world space'"},
		{[]string{"-v", "--flag"}, "-v --flag"},
	}

	for _, test := range tests {
		result := joinArgs(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, result)
		}
	}
}

func TestContainsSpace(t *testing.T) {
	if !containsSpace("hello world") {
		t.Error("Should contain space")
	}
	if containsSpace("helloworld") {
		t.Error("Should not contain space")
	}
}
