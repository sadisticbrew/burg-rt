package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type StateStore struct {
	rootDir string
}

func NewStateStore(rootDir string) *StateStore {
	if rootDir == "" {
		rootDir = DefaultRoot()
	}
	return &StateStore{rootDir: rootDir}
}

func (s *StateStore) Save(state State, containerID string) error {
	containerDir := filepath.Join(s.rootDir, containerID)

	if err := os.MkdirAll(containerDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, " ", "  ")
	if err != nil {
		return err
	}

	err = s.AtomicWrite(containerDir, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StateStore) Load(containerID string) (*State, error) {
	containerDir := filepath.Join(s.rootDir, containerID)

	data, err := os.ReadFile(s.getStateDir(containerDir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("state not found")
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (s *StateStore) Delete(containerID string) error {
	containerDir := filepath.Join(s.rootDir, containerID)
	if err := os.RemoveAll(containerDir); err != nil {
		return err
	}
	return nil
}

func (s *StateStore) AtomicWrite(containerDir string, data []byte) error {
	tmpFile, err := os.CreateTemp(containerDir, "state.json.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpFile.Name(), s.getStateDir(containerDir)); err != nil {
		return err
	}
	return nil
}

func (s *StateStore) Exists(containerID string) bool {
	containerDir := filepath.Join(s.rootDir, containerID)
	_, err := os.Stat(s.getStateDir(containerDir))
	return !os.IsNotExist(err)
}

func (s *StateStore) ChangeStatus(containerID string, status ContainerState) error {
	state, err := s.Load(containerID)
	if err != nil {
		return err
	}
	state.Status = status
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return s.AtomicWrite(filepath.Join(s.rootDir, containerID), data)
}

func (s *StateStore) getStateDir(containerDir string) string {
	return filepath.Join(containerDir, "state.json")
}
