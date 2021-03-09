package main

import "syscall"
import "unsafe"

import "golang.org/x/sys/windows"

var _ unsafe.Pointer

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procSetConsoleMode = modkernel32.NewProc("SetConsoleMode")

func setConsoleMode(console windows.Handle, mode uint32) (uintptr, uintptr, error) {
	r1, r2, e1 := syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(console), uintptr(mode), 0)
	return r1, r2, e1
}
