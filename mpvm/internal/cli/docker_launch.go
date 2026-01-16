package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/takoeight0821/skills/skills-cli/internal/docker"
)

var dockerLaunchCmd = &cobra.Command{
	Use:     "launch",
	Aliases: []string{"create"},
	Short:   "Build and start container",
	Long: `Build the Docker image and start a coding agent container.
If the container already exists, it will be started instead.`,
	RunE: runDockerLaunch,
}

func init() {
	dockerCmd.AddCommand(dockerLaunchCmd)
}

func runDockerLaunch(cmd *cobra.Command, args []string) error {
	containerName := getContainerName()
	imageName := getImageName()
	volumeName := getVolumeName()

	// Check if image exists, if not build it
	imageExists, _ := dockerClient.ImageExists(imageName)
	if !imageExists {
		log.Info("Building Docker image '%s'...", imageName)

		// Find Dockerfile path - look in docker/ subdirectory relative to skills repo
		dockerDir := getDockerDir()
		if dockerDir == "" {
			return fmt.Errorf("cannot find docker directory with Dockerfile")
		}

		// Get user/group ID for Linux compatibility
		buildArgs := map[string]string{
			"USER_ID":  fmt.Sprintf("%d", os.Getuid()),
			"GROUP_ID": fmt.Sprintf("%d", os.Getgid()),
		}

		if err := dockerClient.Build(docker.BuildOptions{
			Tag:       imageName,
			Context:   dockerDir,
			BuildArgs: buildArgs,
		}); err != nil {
			return fmt.Errorf("failed to build image: %w", err)
		}
		log.Success("Image built successfully")
	} else {
		log.Info("Docker image '%s' already exists.", imageName)
	}

	// Create volume for persistent data
	if err := dockerClient.VolumeCreate(volumeName); err != nil {
		log.Warn("Failed to create volume: %v", err)
	}

	// Check if container exists
	exists, _ := dockerClient.ContainerExists(containerName)
	if exists {
		log.Warn("Container '%s' already exists.", containerName)
		running, _ := dockerClient.ContainerRunning(containerName)
		if !running {
			log.Info("Starting existing container...")
			if err := dockerClient.Start(containerName); err != nil {
				return err
			}
		} else {
			log.Info("Container is already running.")
		}
	} else {
		log.Info("Creating container '%s'...", containerName)
		log.Info("  CPUs: %s", getDockerCPUs())
		log.Info("  Memory: %s", getDockerMemory())

		// Prepare volumes
		volumes := []string{
			fmt.Sprintf("%s:/home/agent/.claude", volumeName),
		}

		// Add SSH agent mount if available
		if sshVol, sshEnv, ok := docker.GetSSHAgentMount(); ok {
			volumes = append(volumes, sshVol)
			log.Info("SSH agent forwarding enabled")
			_ = sshEnv // Will be used in env
		}

		// Add config mounts
		volumes = append(volumes, getConfigMounts()...)

		// Prepare environment
		env := getDockerEnv()

		if err := dockerClient.Create(docker.CreateOptions{
			Name:     containerName,
			Image:    imageName,
			Hostname: "coding-agent",
			CPUs:     getDockerCPUs(),
			Memory:   getDockerMemory(),
			Volumes:  volumes,
			Env:      env,
			Command:  []string{"tail", "-f", "/dev/null"},
		}); err != nil {
			return fmt.Errorf("failed to create container: %w", err)
		}

		if err := dockerClient.Start(containerName); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}
	}

	// Configure git in container
	if err := configureGitInContainer(containerName); err != nil {
		log.Warn("Git configuration failed: %v", err)
	}

	fmt.Println()
	log.Success("Container '%s' is ready!", containerName)
	fmt.Println()
	fmt.Println("Use:")
	fmt.Println("  skills docker ssh      # Interactive shell")
	fmt.Println("  skills docker claude   # Run Claude Code")
	fmt.Println("  skills docker gemini   # Run Gemini CLI")

	return nil
}

func getDockerDir() string {
	// Try to find the docker directory
	// First check relative to executable or working directory
	candidates := []string{
		"docker",
		"../docker",
		filepath.Join(os.Getenv("HOME"), "ghq/github.com/takoeight0821/skills/docker"),
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			dockerfile := filepath.Join(path, "Dockerfile")
			if _, err := os.Stat(dockerfile); err == nil {
				return path
			}
		}
	}
	return ""
}

func getDockerCPUs() string {
	if cfg != nil && cfg.Docker.CPUs != "" {
		return cfg.Docker.CPUs
	}
	return "2"
}

func getDockerMemory() string {
	if cfg != nil && cfg.Docker.Memory != "" {
		return cfg.Docker.Memory
	}
	return "4g"
}

func getConfigMounts() []string {
	var mounts []string
	home := os.Getenv("HOME")

	// Mount ~/.gemini if exists
	geminiDir := filepath.Join(home, ".gemini")
	if info, err := os.Stat(geminiDir); err == nil && info.IsDir() {
		mounts = append(mounts, fmt.Sprintf("%s:/home/agent/.gemini:ro", geminiDir))
	}

	// Mount ~/.aws if exists
	awsDir := filepath.Join(home, ".aws")
	if info, err := os.Stat(awsDir); err == nil && info.IsDir() {
		mounts = append(mounts, fmt.Sprintf("%s:/home/agent/.aws:ro", awsDir))
	}

	return mounts
}

func getDockerEnv() map[string]string {
	env := map[string]string{
		"TERM":      getVMTerm(),
		"COLORTERM": getColorTerm(),
	}

	if cfg != nil {
		if cfg.Git.UserName != "" {
			env["GIT_USER_NAME"] = cfg.Git.UserName
		}
		if cfg.Git.UserEmail != "" {
			env["GIT_USER_EMAIL"] = cfg.Git.UserEmail
		}
	}

	// Add SSH agent socket env if available
	if _, sshEnv, ok := docker.GetSSHAgentMount(); ok {
		env["SSH_AUTH_SOCK"] = sshEnv
	}

	return env
}

func configureGitInContainer(containerName string) error {
	log.Info("Configuring Git in container...")

	if cfg != nil && cfg.Git.UserName != "" {
		if err := dockerClient.Exec(containerName, "git", "config", "--global", "user.name", cfg.Git.UserName); err != nil {
			return err
		}
		log.Info("Git user.name set to: %s", cfg.Git.UserName)
	}

	if cfg != nil && cfg.Git.UserEmail != "" {
		if err := dockerClient.Exec(containerName, "git", "config", "--global", "user.email", cfg.Git.UserEmail); err != nil {
			return err
		}
		log.Info("Git user.email set to: %s", cfg.Git.UserEmail)
	}

	return nil
}
