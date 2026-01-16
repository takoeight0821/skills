package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/docker"
)

var (
	dockerClient docker.Client
)

// dockerCmd represents the docker command group for Docker container management
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker containers",
	Long: `Docker commands for managing coding agent containers.
Supports launch, start, stop, delete, ssh, claude, gemini, status, and logs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Run parent's PersistentPreRunE first
		if rootCmd.PersistentPreRunE != nil {
			if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
		}
		// Initialize docker client
		return initDockerClient()
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}

// initDockerClient initializes the docker client for docker commands
func initDockerClient() error {
	dockerClient = docker.NewClient()
	if !dockerClient.IsInstalled() {
		return fmt.Errorf("docker is not installed. Please install Docker from https://www.docker.com/products/docker-desktop/")
	}
	return nil
}

// Docker configuration helpers
func getContainerName() string {
	if cfg != nil && cfg.Docker.ContainerName != "" {
		return cfg.Docker.ContainerName
	}
	return "coding-agent-docker"
}

func getImageName() string {
	if cfg != nil && cfg.Docker.ImageName != "" {
		return cfg.Docker.ImageName
	}
	return "coding-agent:latest"
}

func getVolumeName() string {
	return getContainerName() + "-claude-data"
}
