package multipass

type Client interface {
	Launch(opts LaunchOptions) error
	VMExists(name string) (bool, error)
	VMRunning(name string) (bool, error)
	GetIP(name string) (string, error)
	Start(name string) error
	Stop(name string) error
	Delete(name string, purge bool) error
	Exec(name string, command ...string) error
	ExecOutput(name string, command ...string) (string, error)
	ExecInteractive(name string, command ...string) error
	Shell(name string) error
	Mount(source, target string) error
	Umount(target string) error
	Transfer(source, dest string) error
	TransferRecursive(source, dest string) error
	WaitForCloudInit(name string, timeoutSeconds int, logFunc func(string, ...interface{})) error
	Info(name string) (string, error)
	List() (string, error)
	IsMounted(vmName, sourcePath string) (bool, error)
	GetMountPoint(vmName, sourcePath string) string
	IsInstalled() bool
}
