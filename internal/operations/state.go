package operations

import (
	"fmt"

	"burg/internal/state"
)

// GetState returns the OCI runtime state for the given container.
func GetState(root, containerID string) (*state.State, error) {
	return nil, fmt.Errorf("container %q not found", containerID)
}
