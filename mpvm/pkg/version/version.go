package version

import (
	"fmt"
	"runtime"
)

// Set via ldflags at build time
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func Info() string {
	return fmt.Sprintf("mpvm %s (commit: %s, built: %s, go: %s)",
		Version, Commit, Date, runtime.Version())
}

func Short() string {
	return Version
}
