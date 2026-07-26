package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func PivotRoot(newRoot string) error {
	newRoot, err := filepath.Abs(newRoot)
	if err != nil {
		fmt.Println("Could not get absolute path:", err)
		return err
	}

	putOld := filepath.Join(newRoot, "put_old")
	err = os.MkdirAll(putOld, 0755)
	if err != nil {
		fmt.Println("Could not create put_old directory:", err)
		return err
	}

	err = unix.Mount(newRoot, newRoot, "", unix.MS_BIND|unix.MS_REC, "")
	if err != nil {
		fmt.Println("Could not mount new root:", err)
		return err
	}

	err = unix.PivotRoot(newRoot, putOld)
	if err != nil {
		fmt.Println("Could not pivot root:", err)
		return err
	}

	err = unix.Chdir("/")
	if err != nil {
		fmt.Println("Could not change directory to root:", err)
		return err
	}

	err = unix.Unmount("/put_old", unix.MNT_DETACH)
	if err != nil {
		fmt.Println("Could not unmount put_old:", putOld, err)
		return err
	}

	return nil
}
