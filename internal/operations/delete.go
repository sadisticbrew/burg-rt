package operations

import "fmt"

// Delete removes all resources associated with a container.
// If force is true, a running container will be killed before deletion.
func Delete(root, containerID string, force bool) error {
	return fmt.Errorf("delete: not yet implemented")
}
