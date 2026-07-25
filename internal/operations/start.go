package operations

import "fmt"

// Start executes the user-defined process in a previously created container.
// The container transitions from "created" to "running".
func Start(root, containerID string) error {
	return fmt.Errorf("start: not yet implemented")
}
