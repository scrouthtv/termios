// +build linux

package termios

import "golang.org/x/sys/unix"

const (
	reqGetTermios = unix.TCGETS
	reqSetTermios = unix.TCSETS
)
