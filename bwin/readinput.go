// +build windows

package bwin

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var procReadConsoleInputW = modkernel32.NewProc("ReadConsoleInputW")

func ReadConsoleInput(console windows.Handle, rec *InputRecord, toread uint32, read *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadConsoleInputW.Addr(), 4,
		uintptr(console), uintptr(unsafe.Pointer(rec)), uintptr(toread),
		uintptr(unsafe.Pointer(read)), 0, 0)

	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
