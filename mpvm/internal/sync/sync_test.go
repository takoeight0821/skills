package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSync(t *testing.T) {
	// Setup temporary directories
	tmpDir, err := os.MkdirTemp("", "sync-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
	}

	for path, content := range files {
		if err := os.WriteFile(filepath.Join(srcDir, path), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("NormalSync", func(t *testing.T) {
		opts := Options{
			SourceDir: srcDir,
			TargetDir: dstDir,
		}

		if err := Sync(opts); err != nil {
			t.Errorf("Sync failed: %v", err)
		}

		// Verify files
		for path, expectedContent := range files {
			content, err := os.ReadFile(filepath.Join(dstDir, path))
			if err != nil {
				t.Errorf("Failed to read synced file %s: %v", path, err)
				continue
			}
			if string(content) != expectedContent {
				t.Errorf("Expected content %q for %s, got %q", expectedContent, path, string(content))
			}
		}
	})

	t.Run("DryRun", func(t *testing.T) {
		dstDirDry := filepath.Join(tmpDir, "dst-dry")
		opts := Options{
			SourceDir: srcDir,
			TargetDir: dstDirDry,
			DryRun:    true,
		}

		if err := Sync(opts); err != nil {
			t.Errorf("Sync failed: %v", err)
		}

		// Verify dst directory does NOT exist (except possibly the top level if we created it before walk)
		// Wait, Sync creates the target dir.
		// Let's check if files exist.
		for path := range files {
			if _, err := os.Stat(filepath.Join(dstDirDry, path)); err == nil {
				t.Errorf("File %s should NOT exist in dry-run", path)
			}
		}
	})

	t.Run("ForceSync", func(t *testing.T) {
		// Modify a file in dst
		existingFile := filepath.Join(dstDir, "file1.txt")
		if err := os.WriteFile(existingFile, []byte("old-content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Sync without force
		opts := Options{
			SourceDir: srcDir,
			TargetDir: dstDir,
			Force:     false,
		}
		if err := Sync(opts); err != nil {
			t.Fatal(err)
		}

		content, _ := os.ReadFile(existingFile)
		if string(content) != "old-content" {
			t.Error("File should NOT have been overwritten without force")
		}

		// Sync with force
		opts.Force = true
		if err := Sync(opts); err != nil {
			t.Fatal(err)
		}

		content, _ = os.ReadFile(existingFile)
		if string(content) != "content1" {
			t.Error("File SHOULD have been overwritten with force")
		}
	})
}
