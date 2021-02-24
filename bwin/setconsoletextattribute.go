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
var procSetConsoleTextAttributeW = modkernel32.NewProc("SetConsoleTextAttributeW")

// SetConsoleTextAttribute sets the specified attributes to the given console.
// attr may be any combination of the Attribute constants.
func SetConsoleTextAttribute(console windows.Handle, attr Attribute) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetConsoleTextAttributeW.Addr(), 2,
		uintptr(console), uintptr(unsafe.Pointer(attr)), 0, 0, 0, 0)

	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
