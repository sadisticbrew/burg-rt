package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

var cloneFlagsMap = map[specs.LinuxNamespaceType]int{
	specs.PIDNamespace:     unix.CLONE_NEWPID,
	specs.UTSNamespace:     unix.CLONE_NEWUTS,
	specs.MountNamespace:   unix.CLONE_NEWNS,
	specs.IPCNamespace:     unix.CLONE_NEWIPC,
	specs.UserNamespace:    unix.CLONE_NEWUSER,
	specs.NetworkNamespace: unix.CLONE_NEWNET,
	specs.TimeNamespace:    unix.CLONE_NEWTIME,
	specs.CgroupNamespace:  unix.CLONE_NEWCGROUP,
}

func CloneFlags(spec *specs.Spec) uintptr {
	if spec.Linux == nil {
		return 0
	}

	var cloneFlags uintptr
	for _, ns := range spec.Linux.Namespaces {
		if ns.Path == "" {
			cloneFlags |= uintptr(cloneFlagsMap[ns.Type])
		}
	}
	return cloneFlags
}

// must be called inside a UTS namespace
func SetHostname(name string) error {
	return unix.Sethostname([]byte(name))
}
