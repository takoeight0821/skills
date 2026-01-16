package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/docker"
)

var dockerGeminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Run Gemini CLI in container",
	Long:  `Run Gemini CLI in a new container with current directory mounted.`,
	RunE:  runDockerGemini,
}

func init() {
	dockerCmd.AddCommand(dockerGeminiCmd)
}

func runDockerGemini(cmd *cobra.Command, args []string) error {
	imageName := getImageName()

	// Check image exists
	exists, _ := dockerClient.ImageExists(imageName)
	if !exists {
		return fmt.Errorf("image '%s' does not exist. Run 'skills docker launch' first", imageName)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Prepare volumes
	volumes := []string{
		fmt.Sprintf("%s:/home/agent/.claude", getVolumeName()),
		fmt.Sprintf("%s:/workspace", cwd),
	}

	// Add SSH agent mount if available
	if sshVol, _, ok := docker.GetSSHAgentMount(); ok {
		volumes = append(volumes, sshVol)
	}

	// Add config mounts
	volumes = append(volumes, getConfigMounts()...)

	log.Info("Working directory: /workspace (mounted from %s)", cwd)

	// Build command: bash -ic 'gemini [args...]'
	geminiArgs := "gemini"
	for _, arg := range args {
		geminiArgs += " " + arg
	}

	return dockerClient.RunInteractive(docker.RunOptions{
		Image:    imageName,
		Hostname: "coding-agent",
		Volumes:  volumes,
		Env:      getDockerEnv(),
		WorkDir:  "/workspace",
		Command:  []string{"bash", "-ic", geminiArgs},
		Remove:   true,
	})
}
