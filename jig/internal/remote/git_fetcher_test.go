package remote

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGitFetcher_Fetch(t *testing.T) {
	t.Run("successfully fetches a public repo", func(t *testing.T) {
		// Create a temp dir for the test
		tempDir, err := os.MkdirTemp("", "git_fetcher_test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		fetcher := NewGitFetcher()
		
		// We use the project's own repo as a test case since it's public and valid.
		dest := filepath.Join(tempDir, "cloned_repo")
		err = fetcher.Fetch(context.Background(), "https://github.com/takoeight0821/skills", dest)
		
		if err != nil {
			t.Errorf("Fetch returned error: %v", err)
		}
		
		// Check if destination exists.
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			t.Errorf("Destination directory was not created")
		}
	})

	t.Run("fails with invalid URL", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "git-fetch-fail-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		fetcher := NewGitFetcher()
		dest := filepath.Join(tempDir, "dest")
		
		// Use an invalid URL
		err = fetcher.Fetch(context.Background(), "https://github.com/takoeight0821/non-existent-repo-12345", dest)
		
		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
		}
	})
}