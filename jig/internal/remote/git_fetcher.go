package remote

import (
	"context"
	"fmt"
	"os/exec"
)

type GitFetcher struct{}

func NewGitFetcher() *GitFetcher {
	return &GitFetcher{}
}

func (f *GitFetcher) Fetch(ctx context.Context, url string, dest string) error {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}

	// Prepare git clone command
	// We use --depth 1 to minimize download size since we only need the latest version
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", url, dest)
	
	// Capture output for debugging if needed
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s: %w", string(output), err)
	}

	return nil
}
