package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
)

var (
	ErrBundleNotFound = errors.New("bundle not found")

	ErrConfigNotFound = errors.New("config.json not found")

	ErrInvalidSpec = errors.New("invalid spec")
)

func checkIfExistOrErr(path string, return_err error) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return return_err
	}
	if err != nil {
		return err
	}
	return nil
}

func DefaultSpec() *specs.Spec {
	return &specs.Spec{
		Version: specs.Version,
		Root: &specs.Root{
			Path:     "rootfs",
			Readonly: false,
		},
		Process: &specs.Process{
			Terminal: false,
			User: specs.User{
				UID: 0,
				GID: 0,
			},
			Args: []string{"/bin/sh"},
			Env: []string{
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				"TERM=xterm",
			},
			Cwd:             "/",
			NoNewPrivileges: true,
		},
		Hostname: "container",
		Mounts: []specs.Mount{
			{
				Destination: "/proc",
				Type:        "proc",
				Source:      "proc",
				Options:     []string{"nosuid", "noexec", "nodev"},
			},
			{
				Destination: "/dev",
				Type:        "tmpfs",
				Source:      "tmpfs",
				Options:     []string{"nosuid", "strictatime", "mode=755", "size=65536k"},
			},
			{
				Destination: "/dev/pts",
				Type:        "devpts",
				Source:      "devpts",
				Options:     []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620"},
			},
			{
				Destination: "/dev/shm",
				Type:        "tmpfs",
				Source:      "shm",
				Options:     []string{"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"},
			},
			{
				Destination: "/sys",
				Type:        "sysfs",
				Source:      "sysfs",
				Options:     []string{"nosuid", "noexec", "nodev", "ro"},
			},
		},
		Linux: &specs.Linux{
			Namespaces: []specs.LinuxNamespace{
				{Type: specs.PIDNamespace},
				{Type: specs.MountNamespace},
				{Type: specs.IPCNamespace},
				{Type: specs.UTSNamespace},
				{Type: specs.NetworkNamespace},
				{Type: specs.CgroupNamespace},
			},
			// MaskedPaths:   specs.DefaultMaskedPaths,
			// ReadonlyPaths: specs.DefaultReadonlyPaths,
		},
	}
}

func validateSpec(spec *specs.Spec) error {
	if spec.Version == "" {
		return fmt.Errorf("%w: specsVersion is required", ErrInvalidSpec)
	}
	if spec.Root == nil {
		return fmt.Errorf("%w: root is required", ErrInvalidSpec)
	}
	if spec.Root.Path == "" {
		return fmt.Errorf("%w: root.path is required", ErrInvalidSpec)
	}
	if spec.Process == nil {
		return fmt.Errorf("%w: process is required", ErrInvalidSpec)
	}
	if len(spec.Process.Args) == 0 {
		return fmt.Errorf("%w: process.args is required", ErrInvalidSpec)
	}
	if spec.Process.Cwd == "" {
		return fmt.Errorf("%w: process.cwd is required", ErrInvalidSpec)
	}
	return nil
}

func ParseSpec(bundlePath string) (*specs.Spec, error) {
	bundlePath, err := filepath.Abs(bundlePath)
	if err != nil {
		return nil, err
	}
	if err := checkIfExistOrErr(bundlePath, ErrBundleNotFound); err != nil {
		return nil, err
	}

	configPath := filepath.Join(bundlePath, "config.json")

	if err := checkIfExistOrErr(configPath, ErrConfigNotFound); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	spec := &specs.Spec{}
	if err = json.Unmarshal(data, spec); err != nil {
		return nil, err
	}
	if err := validateSpec(spec); err != nil {
		return nil, err
	}
	return spec, nil
}

func WriteSpec(spec *specs.Spec, bundlePath string) error {
	if err := os.MkdirAll(bundlePath, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(bundlePath, "config.json")
	return os.WriteFile(configPath, data, 0644)
}
