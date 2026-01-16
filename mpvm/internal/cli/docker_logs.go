package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dockerLogsLines int
)

var dockerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show container logs",
	Long:  `Show logs from the coding agent container.`,
	RunE:  runDockerLogs,
}

func init() {
	dockerLogsCmd.Flags().IntVarP(&dockerLogsLines, "lines", "n", 100, "Number of lines to show")
	dockerCmd.AddCommand(dockerLogsCmd)
}

func runDockerLogs(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()

	exists, _ := dockerClient.ContainerExists(containerName)
	if !exists {
		return fmt.Errorf("container '%s' does not exist", containerName)
	}

	logs, err := dockerClient.Logs(containerName, dockerLogsLines)
	if err != nil {
		return err
	}

	fmt.Print(logs)
	return nil
}
