// +build linux

package termios

import "syscall"

const (
	getTermios = syscall.TCGETS
	setTermios = syscall.TCSETS
)
