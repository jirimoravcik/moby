//go:build linux && cgo

// Package loopback provides utilities to work with loopback devices.
//
// Deprecated: this package is deprecated and will be removed in the next release.

package loopback // import "github.com/docker/docker/pkg/loopback"

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/log"
	"golang.org/x/sys/unix"
)

func getLoopbackBackingFile(file *os.File) (uint64, uint64, error) {
	loopInfo, err := unix.IoctlLoopGetStatus64(int(file.Fd()))
	if err != nil {
		log.G(context.TODO()).Errorf("Error get loopback backing file: %s", err)
		return 0, 0, ErrGetLoopbackBackingFile
	}
	return loopInfo.Device, loopInfo.Inode, nil
}

// SetCapacity reloads the size for the loopback device.
//
// Deprecated: the loopback package is deprected and will be removed in the next release.
func SetCapacity(file *os.File) error {
	if err := unix.IoctlSetInt(int(file.Fd()), unix.LOOP_SET_CAPACITY, 0); err != nil {
		log.G(context.TODO()).Errorf("Error loopbackSetCapacity: %s", err)
		return ErrSetCapacity
	}
	return nil
}

// FindLoopDeviceFor returns a loopback device file for the specified file which
// is backing file of a loop back device.
//
// Deprecated: the loopback package is deprected and will be removed in the next release.
func FindLoopDeviceFor(file *os.File) *os.File {
	var stat unix.Stat_t
	err := unix.Stat(file.Name(), &stat)
	if err != nil {
		return nil
	}
	targetInode := stat.Ino
	targetDevice := uint64(stat.Dev) //nolint: unconvert // the type is 32bit on mips

	for i := 0; true; i++ {
		path := fmt.Sprintf("/dev/loop%d", i)

		file, err := os.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}

			// Ignore all errors until the first not-exist
			// we want to continue looking for the file
			continue
		}

		dev, inode, err := getLoopbackBackingFile(file)
		if err == nil && dev == targetDevice && inode == targetInode {
			return file
		}
		file.Close()
	}

	return nil
}
