package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dockerStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"info"},
	Short:   "Show container status",
	Long:    `Show status information about the coding agent container.`,
	RunE:    runDockerStatus,
}

func init() {
	dockerCmd.AddCommand(dockerStatusCmd)
}

func runDockerStatus(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()

	exists, _ := dockerClient.ContainerExists(containerName)
	if !exists {
		log.Info("Container '%s' does not exist.", containerName)
		return nil
	}

	info, err := dockerClient.Inspect(containerName)
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}
