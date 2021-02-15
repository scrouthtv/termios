// +build darwin freebsd dragonfly openbsd netbsd

package termios

import "syscall"

const (
	getTermios = syscall.TIOCGETA
	setTermios = syscall.TIOCSETA
)
