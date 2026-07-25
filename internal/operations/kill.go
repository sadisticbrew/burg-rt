package operations

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// Kill sends the specified signal to the container's init process.
// If all is true, the signal is sent to all processes in the container.
func Kill(root, containerID string, signal unix.Signal, all bool) error {
	return fmt.Errorf("kill: not yet implemented")
}
