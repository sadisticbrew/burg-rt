package cmd

import (
	"burg/internal/filesystem"
	"burg/internal/namespace"
	"burg/internal/spec"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var childInitCmd = &cobra.Command{
	Use:    "child-init <bundle>",
	Hidden: true,
	Args:   cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle := args[0]
		containerSpec, err := spec.ParseSpec(bundle)
		if err != nil {
			return err
		}

		if err = namespace.SetHostname(containerSpec.Hostname); err != nil {
			return err
		}
		slog.Info("Hostname set successfully.")

		rootfs := filepath.Join(bundle, containerSpec.Root.Path)

		if err = filesystem.MountStandard(rootfs, containerSpec.Mounts); err != nil {
			return err
		}
		slog.Info("Standard filesystem mounted successfully.")

		if err = filesystem.PivotRoot(rootfs); err != nil {
			return err
		}
		slog.Info("Pivot root successful.")

		if err = os.Chdir(containerSpec.Process.Cwd); err != nil {
			return err
		}
		slog.Info("Changed directory successfully.")

		fmt.Println(containerSpec.Process.Env)

		argPath, err := resolvePath(containerSpec.Process.Args[0], containerSpec.Process.Env)
		if err != nil {
			slog.Error("Failed to resolve path:", "error", err)
			return err
		}
		slog.Info("Resolved path:", "path", argPath)

		return unix.Exec(argPath, containerSpec.Process.Args, containerSpec.Process.Env)

	},
}

func init() {
	rootCmd.AddCommand(childInitCmd)
}

func resolvePath(name string, env []string) (string, error) {

	if strings.Contains(name, "/") {
		return name, nil
	}
	var pathVar string
	for _, v := range env {
		if strings.HasPrefix(v, "PATH=") {
			pathVar = v[5:]
			break
		}
	}
	if pathVar == "" {
		return "", fmt.Errorf("PATH not found in environment")
	}

	dirs := strings.Split(pathVar, ":")
	for _, dir := range dirs {
		path := filepath.Join(dir, name)

		if _, err := os.Stat(path); err == nil {
			if isExec, err := isExecutable(path); err == nil && isExec {
				return path, nil
			}
		}
	}
	return "", fmt.Errorf("command not found: %s", name)
}

func isExecutable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return false, nil
	}

	// Check if any executable bit is set (0111 octal mask)
	return info.Mode().Perm()&0111 != 0, nil
}
