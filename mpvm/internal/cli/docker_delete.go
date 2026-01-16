package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dockerDeleteForce bool
)

var dockerDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"destroy"},
	Short:   "Delete container and image",
	Long:    `Delete the coding agent container, image, and data volume.`,
	RunE:    runDockerDelete,
}

func init() {
	dockerDeleteCmd.Flags().BoolVarP(&dockerDeleteForce, "force", "f", false, "Skip confirmation prompt")
	dockerCmd.AddCommand(dockerDeleteCmd)
}

func runDockerDelete(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()
	imageName := getImageName()
	volumeName := getVolumeName()

	containerExists, _ := dockerClient.ContainerExists(containerName)
	imageExists, _ := dockerClient.ImageExists(imageName)

	if !containerExists && !imageExists {
		log.Warn("Container '%s' and image '%s' do not exist.", containerName, imageName)
		return nil
	}

	// Confirmation prompt
	if !dockerDeleteForce {
		log.Warn("This will delete:")
		if containerExists {
			fmt.Printf("  - Container: %s\n", containerName)
		}
		if imageExists {
			fmt.Printf("  - Image: %s\n", imageName)
		}
		fmt.Printf("  - Volume: %s (Claude credentials)\n", volumeName)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure? [y/N] ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			log.Info("Cancelled.")
			return nil
		}
	}

	// Stop container if running
	if containerExists {
		running, _ := dockerClient.ContainerRunning(containerName)
		if running {
			log.Info("Stopping container...")
			_ = dockerClient.Stop(containerName)
		}

		log.Info("Removing container...")
		if err := dockerClient.Remove(containerName, true); err != nil {
			log.Warn("Failed to remove container: %v", err)
		}
	}

	// Remove image
	if imageExists {
		log.Info("Removing image...")
		if err := dockerClient.RemoveImage(imageName); err != nil {
			log.Warn("Failed to remove image: %v", err)
		}
	}

	// Remove volume
	log.Info("Removing volume...")
	_ = dockerClient.VolumeRemove(volumeName)

	log.Success("Deleted.")
	return nil
}
