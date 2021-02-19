// +build windows

package termios

import "syscall"
import "unsafe"
import "golang.org/x/sys/windows"

var _ unsafe.Pointer

const errnoERROR_IO_PENDING = 997
var errERROR_IO_PENDING error = syscall.Errno(997)
var errERROR_EINVAL error = syscall.EINVAL

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0: return errERROR_EINVAL
	case errnoERROR_IO_PENDING: return errERROR_IO_PENDING
	}
	return e
}

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
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
