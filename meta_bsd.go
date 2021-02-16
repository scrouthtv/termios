// +build freebsd openbsd netbsd dragonfly darwin

package termios

import "golang.org/x/sys/unix"

const (
	reqGetTermios = unix.TIOCGETA
	reqSetTermios = unix.TIOCSETA
)
