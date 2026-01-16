package multipass

import (
	"fmt"
	"strings"
)

type MockClient struct {
	// State
	VMs       map[string]bool     // name -> running
	IPs       map[string]string   // name -> ip
	Mounts    map[string][]string // name -> [src]
	ExecCalls []string
	Launched  []string
	Deleted   []string

	// Error injection
	LaunchErr    error
	StartErr     error
	StopErr      error
	DeleteErr    error
	ExecErr      error
	VMExistsErr  error
	VMRunningErr error
}

func NewMockClient() *MockClient {
	return &MockClient{
		VMs:    make(map[string]bool),
		IPs:    make(map[string]string),
		Mounts: make(map[string][]string),
	}
}

func (m *MockClient) Launch(opts LaunchOptions) error {
	if m.LaunchErr != nil {
		return m.LaunchErr
	}
	m.VMs[opts.Name] = true
	m.IPs[opts.Name] = "192.168.64.2" // Default mock IP
	m.Launched = append(m.Launched, opts.Name)
	return nil
}

func (m *MockClient) VMExists(name string) (bool, error) {
	if m.VMExistsErr != nil {
		return false, m.VMExistsErr
	}
	_, exists := m.VMs[name]
	return exists, nil
}

func (m *MockClient) VMRunning(name string) (bool, error) {
	if m.VMRunningErr != nil {
		return false, m.VMRunningErr
	}
	running, exists := m.VMs[name]
	return exists && running, nil
}

func (m *MockClient) GetIP(name string) (string, error) {
	if ip, ok := m.IPs[name]; ok {
		return ip, nil
	}
	return "", fmt.Errorf("VM not found")
}

func (m *MockClient) Start(name string) error {
	if m.StartErr != nil {
		return m.StartErr
	}
	if _, exists := m.VMs[name]; !exists {
		return fmt.Errorf("VM does not exist")
	}
	m.VMs[name] = true
	return nil
}

func (m *MockClient) Stop(name string) error {
	if m.StopErr != nil {
		return m.StopErr
	}
	if _, exists := m.VMs[name]; !exists {
		return fmt.Errorf("VM does not exist")
	}
	m.VMs[name] = false
	return nil
}

func (m *MockClient) Delete(name string, purge bool) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.VMs, name)
	delete(m.IPs, name)
	m.Deleted = append(m.Deleted, name)
	return nil
}

func (m *MockClient) Exec(name string, command ...string) error {
	if m.ExecErr != nil {
		return m.ExecErr
	}
	m.ExecCalls = append(m.ExecCalls, strings.Join(command, " "))
	return nil
}

func (m *MockClient) ExecOutput(name string, command ...string) (string, error) {
	if m.ExecErr != nil {
		return "", m.ExecErr
	}
	m.ExecCalls = append(m.ExecCalls, strings.Join(command, " "))

	// Mock specific commands
	cmdStr := strings.Join(command, " ")
	if strings.Contains(cmdStr, "cloud-init status") {
		return "status: done", nil
	}

	return "", nil
}

func (m *MockClient) ExecInteractive(name string, command ...string) error {
	return m.Exec(name, command...)
}

func (m *MockClient) Shell(name string) error {
	return nil
}

func (m *MockClient) Mount(source, target string) error {
	// Simplified mock
	return nil
}

func (m *MockClient) Umount(target string) error {
	return nil
}

func (m *MockClient) Transfer(source, dest string) error {
	return nil
}

func (m *MockClient) TransferRecursive(source, dest string) error {
	return nil
}

func (m *MockClient) WaitForCloudInit(name string, timeoutSeconds int, logFunc func(string, ...interface{})) error {
	return nil
}

func (m *MockClient) Info(name string) (string, error) {
	return "Name: " + name + "\nState: Running\nIPv4: 192.168.64.2\n", nil
}

func (m *MockClient) List() (string, error) {
	return "", nil
}

func (m *MockClient) IsMounted(vmName, sourcePath string) (bool, error) {
	return false, nil
}

func (m *MockClient) GetMountPoint(vmName, sourcePath string) string {
	return "/mnt/mock"
}

func (m *MockClient) IsInstalled() bool {
	return true
}
