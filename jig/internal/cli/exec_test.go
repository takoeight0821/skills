package cli

import (
	"testing"

	"github.com/takoeight0821/skills/jig/internal/config"
	"github.com/takoeight0821/skills/jig/internal/logging"
	"github.com/takoeight0821/skills/jig/internal/multipass"
)

// Note: exec command uses os/exec to call ssh, which maps to system calls.
// Since we are mocking multipass.Client but NOT os/exec for SSH, we can only verify
// parts of the logic or need to refactor runExec to use client for SSH too (if client supported it).
// The current implementation of runExec calls 'client.GetIP' but then does its own 'exec.Command("ssh"...)'.
// This makes it hard to test without refactoring runExec to delegate the SSH call to the client interface.

// For now, we will test the validation logic and IP retrieval.

func TestExec_Validation(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = true
	mockClient.IPs["coding-agent"] = "1.2.3.4"

	// Test no args
	err := runExec(nil, []string{})
	if err == nil || err.Error() != "no command specified" {
		t.Errorf("Expected 'no command specified', got %v", err)
	}

	// Test only double dash
	err = runExec(nil, []string{"--"})
	if err == nil || err.Error() != "no command specified" {
		t.Errorf("Expected 'no command specified' after dash, got %v", err)
	}
}

func TestExec_VMNotRunning(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	mockClient.VMs["coding-agent"] = false

	err := runExec(nil, []string{"ls"})
	if err == nil {
		t.Error("Exec should fail if VM is not running")
	}
}
