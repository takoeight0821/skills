package multipass

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SystemClient struct {
	binary string
}

func NewClient() Client {
	return &SystemClient{
		binary: "multipass",
	}
}

type LaunchOptions struct {
	Name      string
	CPUs      int
	Memory    string
	Disk      string
	CloudInit string
}

func (c *SystemClient) Launch(opts LaunchOptions) error {
	args := []string{
		"launch",
		"--name", opts.Name,
		"--cpus", fmt.Sprintf("%d", opts.CPUs),
		"--memory", opts.Memory,
		"--disk", opts.Disk,
	}

	if opts.CloudInit != "" {
		args = append(args, "--cloud-init", opts.CloudInit)
	}

	// Show output in real-time for launch (it takes a long time)
	return c.runWithOutput(args...)
}

func (c *SystemClient) runWithOutput(args ...string) error {
	cmd := exec.Command(c.binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) VMExists(name string) (bool, error) {
	output, err := c.runOutput("list", "--format", "csv")
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) > 0 && fields[0] == name {
			return true, nil
		}
	}
	return false, nil
}

func (c *SystemClient) VMRunning(name string) (bool, error) {
	output, err := c.runOutput("list", "--format", "csv")
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) > 1 && fields[0] == name {
			return fields[1] == "Running", nil
		}
	}
	return false, nil
}

func (c *SystemClient) GetIP(name string) (string, error) {
	output, err := c.runOutput("list", "--format", "csv")
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) > 2 && fields[0] == name {
			ip := strings.TrimSpace(fields[2])
			if ip != "" && ip != "--" {
				return ip, nil
			}
		}
	}
	return "", fmt.Errorf("IP address not found for VM %s", name)
}

func (c *SystemClient) Start(name string) error {
	return c.run("start", name)
}

func (c *SystemClient) Stop(name string) error {
	return c.run("stop", name)
}

func (c *SystemClient) Delete(name string, purge bool) error {
	if err := c.run("delete", name); err != nil {
		return err
	}
	if purge {
		return c.run("purge")
	}
	return nil
}

func (c *SystemClient) Exec(name string, command ...string) error {
	args := append([]string{"exec", name, "--"}, command...)
	return c.run(args...)
}

func (c *SystemClient) ExecOutput(name string, command ...string) (string, error) {
	args := append([]string{"exec", name, "--"}, command...)
	return c.runOutput(args...)
}

func (c *SystemClient) ExecInteractive(name string, command ...string) error {
	args := append([]string{"exec", name, "--"}, command...)
	cmd := exec.Command(c.binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) Shell(name string) error {
	cmd := exec.Command(c.binary, "shell", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *SystemClient) Mount(source, target string) error {
	return c.run("mount", source, target)
}

func (c *SystemClient) Umount(target string) error {
	return c.run("umount", target)
}

func (c *SystemClient) Transfer(source, dest string) error {
	return c.run("transfer", source, dest)
}

func (c *SystemClient) TransferRecursive(source, dest string) error {
	return c.run("transfer", "-r", source, dest)
}

func (c *SystemClient) WaitForCloudInit(name string, timeoutSeconds int, logFunc func(string, ...interface{})) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	checkCount := 0
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cloud-init timed out after %ds", timeoutSeconds)
		case <-ticker.C:
			checkCount++
			output, err := c.ExecOutput(name, "cloud-init", "status")
			if err != nil {
				if logFunc != nil && checkCount%3 == 0 {
					logFunc("Waiting for VM to be ready... (%ds)", checkCount*10)
				}
				continue // VM might not be ready yet
			}

			outputLower := strings.ToLower(output)
			if strings.Contains(outputLower, "done") {
				return nil
			}
			if strings.Contains(outputLower, "error") {
				return fmt.Errorf("cloud-init failed: %s", strings.TrimSpace(output))
			}
			if logFunc != nil && checkCount%3 == 0 {
				status := strings.TrimSpace(output)
				if status != "" {
					logFunc("Cloud-init status: %s", status)
				}
			}
		}
	}
}

func (c *SystemClient) Info(name string) (string, error) {
	return c.runOutput("info", name)
}

func (c *SystemClient) List() (string, error) {
	return c.runOutput("list")
}

func (c *SystemClient) IsMounted(vmName, sourcePath string) (bool, error) {
	output, err := c.runOutput("info", vmName)
	if err != nil {
		return false, err
	}

	// Parse mounts from info output
	return strings.Contains(output, sourcePath), nil
}

func (c *SystemClient) GetMountPoint(vmName, sourcePath string) string {
	// Convert source path to mount point
	// e.g., /home/user/project -> /mnt/home-user-project
	clean := strings.ReplaceAll(sourcePath, "/", "-")
	clean = strings.TrimPrefix(clean, "-")
	return filepath.Join("/mnt", clean)
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

func (c *SystemClient) IsInstalled() bool {
	_, err := exec.LookPath(c.binary)
	return err == nil
}
