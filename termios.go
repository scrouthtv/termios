package termios

import "errors"

var ErrorClosed error = errors.New("I/O error: terminal is closed")

type Terminal interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	IsOpen() bool
	IsRaw() bool

	Close()
	SetRaw(bool) error
}
