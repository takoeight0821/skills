package docker

// Client defines the interface for Docker container operations
type Client interface {
	// IsInstalled checks if Docker is available on the system
	IsInstalled() bool

	// Container lifecycle
	ContainerExists(name string) (bool, error)
	ContainerRunning(name string) (bool, error)
	Create(opts CreateOptions) error
	Start(name string) error
	Stop(name string) error
	Remove(name string, force bool) error

	// Image management
	ImageExists(name string) (bool, error)
	Build(opts BuildOptions) error
	RemoveImage(name string) error

	// Volume management
	VolumeCreate(name string) error
	VolumeRemove(name string) error

	// Container execution
	Exec(name string, command ...string) error
	ExecOutput(name string, command ...string) (string, error)
	ExecInteractive(name string, command ...string) error
	ExecInteractiveWithEnv(name string, env map[string]string, command ...string) error

	// Run a new container (one-shot)
	Run(opts RunOptions) error
	RunInteractive(opts RunOptions) error

	// Container info
	Logs(name string, lines int) (string, error)
	Inspect(name string) (string, error)
}
