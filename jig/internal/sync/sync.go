package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options defines the parameters for the synchronization process.
type Options struct {
	SourceDir string
	TargetDir string
	DryRun    bool
	Force     bool
}

// Sync synchronizes files from SourceDir to TargetDir.
func Sync(opts Options) error {
	if opts.SourceDir == "" {
		return fmt.Errorf("source directory is required")
	}
	if opts.TargetDir == "" {
		return fmt.Errorf("target directory is required")
	}

	// Ensure target directory exists
	if !opts.DryRun {
		if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	return filepath.Walk(opts.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		relPath, err := filepath.Rel(opts.SourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		// Skip .git directory and other hidden files/dirs if needed
		// For now, let's keep it simple and just sync everything in the skills dir
		
		targetPath := filepath.Join(opts.TargetDir, relPath)

		if info.IsDir() {
			if opts.DryRun {
				fmt.Printf("[dry-run] Would create directory: %s\n", targetPath)
				return nil
			}
			return os.MkdirAll(targetPath, 0755)
		}

		// It's a file
		if opts.DryRun {
			fmt.Printf("[dry-run] Would copy file: %s -> %s\n", path, targetPath)
			return nil
		}

		if err := copyFile(path, targetPath, opts.Force); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w", path, targetPath, err)
		}

		return nil
	})
}

func copyFile(src, dst string, force bool) error {
	if !force {
		if _, err := os.Stat(dst); err == nil {
			// File exists, skip if not forced
			return nil
		}
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Ensure the parent directory of dst exists (in case it wasn't created)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
