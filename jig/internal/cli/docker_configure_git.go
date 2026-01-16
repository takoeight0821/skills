package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dockerConfigureGitCmd = &cobra.Command{
	Use:   "configure-git",
	Short: "Configure git in container",
	Long:  `Re-configure git settings in the running container.`,
	RunE:  runDockerConfigureGit,
}

func init() {
	dockerCmd.AddCommand(dockerConfigureGitCmd)
}

func runDockerConfigureGit(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()

	running, err := dockerClient.ContainerRunning(containerName)
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("container '%s' is not running", containerName)
	}

	return configureGitInContainer(containerName)
}
