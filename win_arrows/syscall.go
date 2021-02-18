// +build syscall

package main

import (
	"syscall"
	"unsafe"
)

var in syscall.Handle
var oldMode dword

const enable_window_input = 0x8

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var proc_get_console_mode = kernel32.NewProc("GetConsoleMode")
var proc_set_console_mode = kernel32.NewProc("SetConsoleMode")
var proc_read_console_input = kernel32.NewProc("ReadConsoleInputW")

func SetRaw() error {
	var err error
	in, err = syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		return err
	}

	err = get_console_mode(in, &oldMode)

	err = set_console_mode(in, enable_window_input)
	if err != nil {
		return err
	}

	return nil
}

// Read reads a single keypress
func Read(p []uint16) (int, error) {
	//return windows.Read(in, p)
	var r input_record
	err := read_console_input(in, &r)
	return 0, err
}

func Close() {
	set_console_mode(in, oldMode)
	syscall.Close(in)
}

// INTERNALS:
type word uint16
type dword uint32

type input_record struct {
	event_type word
	_          [2]byte
	event      [16]byte
}

var tmp_arg dword

func get_console_mode(h syscall.Handle, mode *dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_mode.Addr(),
	2, uintptr(h), uintptr(unsafe.Pointer(mode)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func set_console_mode(h syscall.Handle, mode dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_mode.Addr(),
	2, uintptr(h), uintptr(mode), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func read_console_input(h syscall.Handle, record *input_record) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_read_console_input.Addr(),
	4, uintptr(h), uintptr(unsafe.Pointer(record)), 1, uintptr(unsafe.Pointer(&tmp_arg)), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
