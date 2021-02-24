// +build windows

package bwin

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

const errnoERROR_IO_PENDING = 997

var errERROR_IO_PENDING error = syscall.Errno(997)
var errERROR_EINVAL error = syscall.EINVAL

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	return e
}

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
