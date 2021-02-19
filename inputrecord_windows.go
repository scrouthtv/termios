// +build windows

package termios

type InputRecord struct {
	Type uint16
	_ [2]byte
	Data [6]uint16
}
