package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSpec_ValidBundle(t *testing.T) {
	// Create temp bundle
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"root": {
			"path": "rootfs",
			"readonly": true
		},
		"process": {
			"args": ["/bin/sh"],
			"cwd": "/"
		}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseSpec(bundle)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if spec.Version != "1.0.2" {
		t.Errorf("expected version 1.0.2, got %s", spec.Version)
	}
	if spec.Root.Path != "rootfs" {
		t.Errorf("expected root.path rootfs, got %s", spec.Root.Path)
	}
	if !spec.Root.Readonly {
		t.Error("expected root.readonly true")
	}
	if len(spec.Process.Args) != 1 || spec.Process.Args[0] != "/bin/sh" {
		t.Errorf("expected args [/bin/sh], got %v", spec.Process.Args)
	}
}

func TestLoadSpec_MinimalConfig(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	// Minimal valid config: version, root, process with args and cwd
	config := `{
		"ociVersion": "1.0.2",
		"root": {"path": "rootfs"},
		"process": {"args": ["/bin/echo", "hello"], "cwd": "/"}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseSpec(bundle)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if len(spec.Process.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(spec.Process.Args))
	}
}

func TestLoadSpec_MissingBundle(t *testing.T) {
	_, err := ParseSpec("/nonexistent/bundle")
	if err == nil {
		t.Fatal("expected error for missing bundle")
	}
}

func TestLoadSpec_MissingConfig(t *testing.T) {
	bundle := t.TempDir()
	// Don't create config.json

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing config.json")
	}
}

func TestLoadSpec_InvalidJSON(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	if err := os.WriteFile(configPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadSpec_MissingVersion(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"root": {"path": "rootfs"},
		"process": {"args": ["/bin/sh"], "cwd": "/"}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing ociVersion")
	}
}

func TestLoadSpec_MissingRoot(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"process": {"args": ["/bin/sh"], "cwd": "/"}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing root")
	}
}

func TestLoadSpec_MissingProcess(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"root": {"path": "rootfs"}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing process")
	}
}

func TestLoadSpec_MissingArgs(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"root": {"path": "rootfs"},
		"process": {"cwd": "/"}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing process.args")
	}
}

func TestLoadSpec_MissingCwd(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"root": {"path": "rootfs"},
		"process": {"args": ["/bin/sh"]}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSpec(bundle)
	if err == nil {
		t.Fatal("expected error for missing process.cwd")
	}
}

func TestWriteSpec(t *testing.T) {
	bundle := t.TempDir()
	spec := DefaultSpec()

	if err := WriteSpec(spec, bundle); err != nil {
		t.Fatalf("WriteSpec failed: %v", err)
	}

	// Load it back
	loaded, err := ParseSpec(bundle)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if loaded.Version != spec.Version {
		t.Errorf("version mismatch: %s != %s", loaded.Version, spec.Version)
	}
}

func TestLoadSpec_FullConfig(t *testing.T) {
	bundle := t.TempDir()
	configPath := filepath.Join(bundle, "config.json")

	config := `{
		"ociVersion": "1.0.2",
		"root": {
			"path": "rootfs",
			"readonly": true
		},
		"process": {
			"terminal": true,
			"user": {
				"uid": 1000,
				"gid": 1000,
				"additionalGids": [100, 200]
			},
			"args": ["/bin/bash", "-c", "echo hello"],
			"env": ["PATH=/usr/bin", "HOME=/home/user"],
			"cwd": "/home/user",
			"capabilities": {
				"bounding": ["CAP_NET_BIND_SERVICE"],
				"effective": ["CAP_NET_BIND_SERVICE"],
				"permitted": ["CAP_NET_BIND_SERVICE"]
			},
			"noNewPrivileges": true
		},
		"hostname": "test-container",
		"mounts": [
			{
				"destination": "/tmp",
				"type": "tmpfs",
				"source": "tmpfs",
				"options": ["nosuid", "nodev", "mode=1777"]
			}
		],
		"linux": {
			"namespaces": [
				{"type": "pid"},
				{"type": "network"},
				{"type": "mount"}
			],
			"resources": {
				"memory": {"limit": 104857600},
				"cpu": {"shares": 512},
				"pids": {"limit": 100}
			}
		}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseSpec(bundle)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if spec.Hostname != "test-container" {
		t.Errorf("expected hostname test-container, got %s", spec.Hostname)
	}
	if spec.Process.User.UID != 1000 {
		t.Errorf("expected uid 1000, got %d", spec.Process.User.UID)
	}
	if len(spec.Process.User.AdditionalGids) != 2 {
		t.Errorf("expected 2 additional gids, got %d", len(spec.Process.User.AdditionalGids))
	}
	if spec.Linux == nil || spec.Linux.Resources == nil {
		t.Fatal("expected linux.resources")
	}
	if spec.Linux.Resources.Memory == nil || *spec.Linux.Resources.Memory.Limit != 104857600 {
		t.Error("expected memory limit 104857600")
	}
	if spec.Linux.Resources.Pids == nil || *spec.Linux.Resources.Pids.Limit != 100 {
		t.Error("expected pids limit 100")
	}
}
