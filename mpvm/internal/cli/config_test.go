package cli

import (
	"testing"

	"github.com/takoeight0821/skills/skills-cli/internal/config"
	"github.com/takoeight0821/skills/skills-cli/internal/logging"
	"github.com/takoeight0821/skills/skills-cli/internal/multipass"
)

func TestConfig_Show(t *testing.T) {
	mockClient := multipass.NewMockClient()
	client = mockClient
	cfg = config.DefaultConfig()
	log = logging.Default()

	err := runConfigShow(nil, []string{})
	if err != nil {
		t.Errorf("Config show failed: %v", err)
	}
}

func TestConfig_Path(t *testing.T) {
	err := runConfigPath(nil, []string{})
	if err != nil {
		t.Errorf("Config path failed: %v", err)
	}
}
