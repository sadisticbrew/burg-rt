package operations

import (
	"errors"
	"log/slog"
	"path/filepath"

	"burg/internal/spec"
	"burg/internal/state"
)

var (
	ErrInvalidID       = errors.New("invalid container id")
	ErrBundleNotFound  = errors.New("bundle not found")
	ErrContainerExists = errors.New("container already exists")
)

type CreateOptions struct {
	ID          string
	Bundle      string
	Root        string
	Annotations map[string]string
}

// Create sets up a container from an OCI bundle without starting the user process.
// The container will be in the "created" state after this call.
func Create(opts CreateOptions) error {
	slog.Info("Creating container:", "id", opts.ID)

	if opts.ID == "" {
		return ErrInvalidID
	}

	if !filepath.IsAbs(opts.Bundle) {
		bundlePath, err := filepath.Abs(opts.Bundle)
		if err != nil {
			return err
		}
		opts.Bundle = bundlePath
	}

	stateStore := state.NewStateStore(opts.Root)
	if stateStore.Exists(opts.ID) {
		return ErrContainerExists
	}

	slog.Info("Parsing config.json")
	spec, err := spec.ParseSpec(opts.Bundle)
	if err != nil {
		return err
	}
	slog.Info("Successfully parsed config.json")

	slog.Info("Creating state store")
	state := state.State{
		Version:     spec.Version,
		ID:          opts.ID,
		Status:      state.StateCreating,
		Bundle:      opts.Bundle,
		Annotations: opts.Annotations,
	}
	err = stateStore.Save(state, opts.ID)
	if err != nil {
		return err
	}
	slog.Info("Successfully saved state")

	

	return nil
}
