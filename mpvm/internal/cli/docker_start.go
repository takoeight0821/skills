package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start container",
	Long:  `Start an existing coding agent container.`,
	RunE:  runDockerStart,
}

func init() {
	dockerCmd.AddCommand(dockerStartCmd)
}

func runDockerStart(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()

	exists, err := dockerClient.ContainerExists(containerName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("container '%s' does not exist. Run 'skills docker launch' first", containerName)
	}

	running, _ := dockerClient.ContainerRunning(containerName)
	if running {
		log.Info("Container '%s' is already running.", containerName)
		return nil
	}

	log.Info("Starting container '%s'...", containerName)
	if err := dockerClient.Start(containerName); err != nil {
		return err
	}

	// Re-configure git after start
	if err := configureGitInContainer(containerName); err != nil {
		log.Warn("Git configuration failed: %v", err)
	}

	log.Success("Container started.")
	return nil
}
