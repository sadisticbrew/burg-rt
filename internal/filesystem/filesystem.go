package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

var MountFlags = map[string]uintptr{
	"ro":          unix.MS_RDONLY,
	"nosuid":      unix.MS_NOSUID,
	"nodev":       unix.MS_NODEV,
	"noexec":      unix.MS_NOEXEC,
	"sync":        unix.MS_SYNCHRONOUS,
	"remount":     unix.MS_REMOUNT,
	"mand":        unix.MS_MANDLOCK,
	"dirsync":     unix.MS_DIRSYNC,
	"noatime":     unix.MS_NOATIME,
	"nodiratime":  unix.MS_NODIRATIME,
	"bind":        unix.MS_BIND,
	"rbind":       unix.MS_BIND | unix.MS_REC,
	"move":        unix.MS_MOVE,
	"rec":         unix.MS_REC,
	"silent":      unix.MS_SILENT,
	"relatime":    unix.MS_RELATIME,
	"strictatime": unix.MS_STRICTATIME,
	"private":     unix.MS_PRIVATE,
	"rprivate":    unix.MS_PRIVATE | unix.MS_REC,
	"shared":      unix.MS_SHARED,
	"rshared":     unix.MS_SHARED | unix.MS_REC,
	"slave":       unix.MS_SLAVE,
	"rslave":      unix.MS_SLAVE | unix.MS_REC,
	"unbindable":  unix.MS_UNBINDABLE,
	"runbindable": unix.MS_UNBINDABLE | unix.MS_REC,
}

func parseMountOptions(opts []string) (flags uintptr, data string) {
	var dataOpts []string
	flags = uintptr(0)
	for _, opt := range opts {
		if flag, ok := MountFlags[opt]; ok {
			flags |= flag
		} else {
			dataOpts = append(dataOpts, opt)
		}
	}

	if len(dataOpts) > 0 {
		for _, opt := range dataOpts {
			data += opt + ","
		}
		data = data[:len(data)-1]
	}
	return flags, data
}

func MakePrivate(target string) error {
	if err := unix.Mount("", target, "", unix.MS_PRIVATE|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("make private %s: %w", target, err)
	}
	return nil
}

func MountStandard(rootfs string, mounts []specs.Mount) error {
	if err := MakePrivate("/"); err != nil {
		return fmt.Errorf("Failed to make private: %w", err)
	}

	for _, mount := range mounts {
		target := filepath.Join(rootfs, mount.Destination)

		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("could not create directory %s: %w", target, err)
		}

		flags, data := parseMountOptions(mount.Options)
		if err := unix.Mount(mount.Source, target, mount.Type, flags, data); err != nil {
			return fmt.Errorf("mount %s: %w", target, err)
		}
	}

	return nil
}
