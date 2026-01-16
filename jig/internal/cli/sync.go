package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/jig/internal/sync"
)

var (
	syncDryRun  bool
	syncForce   bool
	syncGlobal  bool
	syncProject bool
	syncSource  string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize skills",
	Long:  `Synchronize skills from the central repository to local directories.`,
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Preview sync without making changes")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "Overwrite existing files")
	syncCmd.Flags().BoolVar(&syncGlobal, "global", false, "Sync to ~/.claude/skills")
	syncCmd.Flags().BoolVar(&syncProject, "project", false, "Sync to .claude/skills")
	syncCmd.Flags().StringVar(&syncSource, "source", "", "Source skills directory (overrides default)")

	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	if !syncGlobal && !syncProject {
		return fmt.Errorf("either --global or --project must be specified")
	}

	sourceDir := syncSource
	if sourceDir == "" {
		// Try to find skills directory relative to the binary or current repo
		// For now, let's check if we are in the skills repo
		cwd, _ := os.Getwd()
		potentialSource := filepath.Join(cwd, "skills")
		if _, err := os.Stat(potentialSource); err == nil {
			sourceDir = potentialSource
		} else {
			// Fallback: check one level up (if running from jig/ or similar)
			potentialSource = filepath.Join(cwd, "..", "skills")
			if _, err := os.Stat(potentialSource); err == nil {
				sourceDir = potentialSource
			}
		}
	}

	if sourceDir == "" {
		return fmt.Errorf("could not find skills source directory. Use --source to specify it")
	}

	var targetDir string
	if syncGlobal {
		home, _ := os.UserHomeDir()
		targetDir = filepath.Join(home, ".claude", "skills")
	} else if syncProject {
		cwd, _ := os.Getwd()
		targetDir = filepath.Join(cwd, ".claude", "skills")
	}

	fmt.Printf("Syncing skills from %s to %s\n", sourceDir, targetDir)
	if syncDryRun {
		fmt.Println("Running in DRY-RUN mode")
	}

	opts := sync.Options{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		DryRun:    syncDryRun,
		Force:     syncForce,
	}

	if err := sync.Sync(opts); err != nil {
		return err
	}

	if !syncDryRun {
		fmt.Println("Successfully synchronized skills.")
	} else {
		fmt.Println("Dry-run complete.")
	}

	return nil
}
