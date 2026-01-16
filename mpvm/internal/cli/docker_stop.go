package cli

import (
	"github.com/spf13/cobra"
)

var dockerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop container",
	Long:  `Stop the running coding agent container.`,
	RunE:  runDockerStop,
}

func init() {
	dockerCmd.AddCommand(dockerStopCmd)
}

func runDockerStop(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()

	exists, _ := dockerClient.ContainerExists(containerName)
	if !exists {
		log.Warn("Container '%s' does not exist.", containerName)
		return nil
	}

	log.Info("Stopping container '%s'...", containerName)
	if err := dockerClient.Stop(containerName); err != nil {
		return err
	}

	log.Success("Container stopped.")
	return nil
}
