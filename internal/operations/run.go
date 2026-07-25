package operations

import "fmt"

// Run creates and immediately starts a container from an OCI bundle.
// This is equivalent to calling Create followed by Start.
func Run(root, containerID, bundle string, detach bool, pidFile, consoleSocket string) error {
	return fmt.Errorf("run: not yet implemented")
}
