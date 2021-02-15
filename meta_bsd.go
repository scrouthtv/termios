// +build freebsd openbsd netbsd

package termios

import "golang.org/x/sys/unix"

const (
	reqGetTermios = unix.TIOCGETA
	reqSetTermios = unix.TIOCSETA
)
