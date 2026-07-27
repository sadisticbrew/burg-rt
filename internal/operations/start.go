package operations

import (
	"burg/internal/namespace"
	"burg/internal/spec"
	"burg/internal/state"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

// Start executes the user-defined process in a previously created container.
// The container transitions from "created" to "running".
func Start(root, containerID string) error {
	stateStore := state.NewStateStore(root)
	if !stateStore.Exists(containerID) {
		return fmt.Errorf("Container with id: %s not found", containerID)
	}

	st, err := stateStore.Load(containerID)
	if err != nil {
		return err
	}

	// if st.Status != state.StateCreated {
	// 	return fmt.Errorf("Container is not in created state")
	// }

	containerSpec, err := spec.ParseSpec(st.Bundle)
	if err != nil {
		return err
	}

	childCmd := exec.Command("/proc/self/exe", "child-init", st.Bundle)

	childCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: namespace.CloneFlags(containerSpec),
	}

	if containerSpec.Process.Terminal {
		childCmd.Stdout, childCmd.Stderr, childCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	} else {
		childCmd.Stdout, childCmd.Stderr, childCmd.Stdin = nil, nil, nil
	}

	if err := childCmd.Start(); err != nil {
		return err
	}

	st.Pid = childCmd.Process.Pid
	st.Status = state.StateRunning
	stateStore.Save(*st, containerID)

	slog.Info("Saved PID and updated status. Waiting.")

	if err := childCmd.Wait(); err != nil {
		return err
	}

	return nil
}
