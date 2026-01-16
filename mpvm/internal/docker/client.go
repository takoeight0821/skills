package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CreateOptions contains options for creating a container
type CreateOptions struct {
	Name     string
	Image    string
	Hostname string
	CPUs     string
	Memory   string
	Volumes  []string // format: "source:dest" or "source:dest:ro"
	Env      map[string]string
	WorkDir  string
	Command  []string
}

// BuildOptions contains options for building an image
type BuildOptions struct {
	Tag        string
	Dockerfile string
	Context    string
	BuildArgs  map[string]string
}

// RunOptions contains options for running a temporary container
type RunOptions struct {
	Image    string
	Name     string
	Hostname string
	Volumes  []string
	Env      map[string]string
	WorkDir  string
	Command  []string
	Remove   bool // --rm flag
}

// SystemClient implements the Docker client interface using system commands
type SystemClient struct {
	binary string
}

func NewClient() Client {
	return &SystemClient{
		binary: "docker",
	}
}

func (c *SystemClient) IsInstalled() bool {
	_, err := exec.LookPath(c.binary)
	return err == nil
}

func (c *SystemClient) ContainerExists(name string) (bool, error) {
	output, err := c.runOutput("ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == name {
			return true, nil
		}
	}
	return false, nil
}

func (c *SystemClient) ContainerRunning(name string) (bool, error) {
	output, err := c.runOutput("ps", "--format", "{{.Names}}")
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == name {
			return true, nil
		}
	}
	return false, nil
}

func (c *SystemClient) Create(opts CreateOptions) error {
	args := []string{"create", "--name", opts.Name}

	if opts.Hostname != "" {
		args = append(args, "--hostname", opts.Hostname)
	}
	if opts.CPUs != "" {
		args = append(args, "--cpus", opts.CPUs)
	}
	if opts.Memory != "" {
		args = append(args, "--memory", opts.Memory)
	}
	for _, v := range opts.Volumes {
		args = append(args, "-v", v)
	}
	for k, v := range opts.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	if opts.WorkDir != "" {
		args = append(args, "-w", opts.WorkDir)
	}
	args = append(args, "-it", opts.Image)
	args = append(args, opts.Command...)

	return c.run(args...)
}

func (c *SystemClient) Start(name string) error {
	return c.run("start", name)
}

func (c *SystemClient) Stop(name string) error {
	return c.run("stop", name)
}

func (c *SystemClient) Remove(name string, force bool) error {
	if force {
		return c.run("rm", "-f", name)
	}
	return c.run("rm", name)
}

func (c *SystemClient) ImageExists(name string) (bool, error) {
	err := c.run("image", "inspect", name)
	if err != nil {
		if strings.Contains(err.Error(), "No such image") {
			return false, nil
		}
		return false, nil // Treat errors as "not exists" for safety
	}
	return true, nil
}

func (c *SystemClient) Build(opts BuildOptions) error {
	args := []string{"build", "-t", opts.Tag}

	if opts.Dockerfile != "" {
		args = append(args, "-f", opts.Dockerfile)
	}
	for k, v := range opts.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, opts.Context)

	return c.runWithOutput(args...)
}

func (c *SystemClient) RemoveImage(name string) error {
	return c.run("rmi", name)
}

func (c *SystemClient) VolumeCreate(name string) error {
	// Create silently succeeds if volume exists
	return c.run("volume", "create", name)
}

func (c *SystemClient) VolumeRemove(name string) error {
	return c.run("volume", "rm", name)
}

func (c *SystemClient) Exec(name string, command ...string) error {
	args := append([]string{"exec", name}, command...)
	return c.run(args...)
}

func (c *SystemClient) ExecOutput(name string, command ...string) (string, error) {
	args := append([]string{"exec", name}, command...)
	return c.runOutput(args...)
}

func (c *SystemClient) ExecInteractive(name string, command ...string) error {
	args := append([]string{"exec", "-it", name}, command...)
	cmd := exec.Command(c.binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) ExecInteractiveWithEnv(name string, env map[string]string, command ...string) error {
	args := []string{"exec", "-it"}
	for k, v := range env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, name)
	args = append(args, command...)

	cmd := exec.Command(c.binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) Run(opts RunOptions) error {
	args := c.buildRunArgs(opts)
	return c.run(args...)
}

func (c *SystemClient) RunInteractive(opts RunOptions) error {
	args := c.buildRunArgs(opts)
	cmd := exec.Command(c.binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) buildRunArgs(opts RunOptions) []string {
	args := []string{"run"}

	if opts.Remove {
		args = append(args, "--rm")
	}
	args = append(args, "-it")

	if opts.Name != "" {
		args = append(args, "--name", opts.Name)
	}
	if opts.Hostname != "" {
		args = append(args, "--hostname", opts.Hostname)
	}
	for _, v := range opts.Volumes {
		args = append(args, "-v", v)
	}
	for k, v := range opts.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	if opts.WorkDir != "" {
		args = append(args, "-w", opts.WorkDir)
	}
	args = append(args, opts.Image)
	args = append(args, opts.Command...)

	return args
}

func (c *SystemClient) Logs(name string, lines int) (string, error) {
	return c.runOutput("logs", "--tail", fmt.Sprintf("%d", lines), name)
}

func (c *SystemClient) Inspect(name string) (string, error) {
	return c.runOutput("inspect", "--format",
		`Container: {{.Name}}
State: {{.State.Status}}
Created: {{.Created}}
Image: {{.Config.Image}}`, name)
}

// GetSSHAgentMount returns the appropriate SSH agent socket mount for the current OS
func GetSSHAgentMount() (volume, envVar string, ok bool) {
	switch runtime.GOOS {
	case "darwin":
		// macOS: Docker Desktop's built-in SSH agent forwarding
		sockPath := "/run/host-services/ssh-auth.sock"
		if _, err := os.Stat(sockPath); err == nil {
			return sockPath + ":/ssh-agent.sock", "/ssh-agent.sock", true
		}
	case "linux":
		// Linux: Mount the SSH agent socket directly
		if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
			if _, err := os.Stat(sock); err == nil {
				return sock + ":/ssh-agent.sock", "/ssh-agent.sock", true
			}
		}
	}
	return "", "", false
}

func (c *SystemClient) run(args ...string) error {
	cmd := exec.Command(c.binary, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return fmt.Errorf("%s: %s", err, errMsg)
		}
		return err
	}
	return nil
}

func (c *SystemClient) runOutput(args ...string) (string, error) {
	cmd := exec.Command(c.binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return "", fmt.Errorf("%s: %s", err, errMsg)
		}
		return "", err
	}
	return stdout.String(), nil
}

func (c *SystemClient) runWithOutput(args ...string) error {
	cmd := exec.Command(c.binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
