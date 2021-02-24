// +build windows

package bwin

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var procSetConsoleTextAttributeW = modkernel32.NewProc("SetConsoleTextAttribute")

// SetConsoleTextAttribute sets the specified attributes to the given console.
// attr may be any combination of the Attribute constants.
func SetConsoleTextAttribute(console windows.Handle, attr Attribute) (err error) {
	ua := uint16(attr)
	r1, _, e1 := syscall.Syscall6(procSetConsoleTextAttributeW.Addr(), 2,
		uintptr(console), uintptr(ua), 0, 0, 0, 0)

	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
