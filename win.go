// +build windows

package termios

import (
	"errors"
	
	"golang.org/x/sys/windows"
)

var ErrorClosed error = errors.New("I/O error: terminal is closed")

type winTerm struct {
	in windows.Handle
	out windows.Handle
	ready bool
	isRaw bool
}

func Open() (*winTerm, error) {
	var err error

	var in, out windows.Handle
	in, err = windows.Open("CONIN$", windows.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	out, err = windows.Open("CONOUT$", windows.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	var t winTerm = winTerm{in, out, true, false}

	return &t, nil
}

func (t *winTerm) IsOpen() bool {
	return t.ready && t.in != windows.InvalidHandle && t.out != windows.InvalidHandle
}

func (t *winTerm) IsRaw() bool {
	return t.isRaw
}

func (t *winTerm) Read(p []byte) (int, error) {
	return windows.Read(t.in, p)
}

func (t *winTerm) Write(p []byte) (int, error) {
	return windows.Write(t.out, p)
}

func (t *winTerm) Close() {
	t.ready = false
	windows.Close(t.in)
	windows.Close(t.out)
	t.in = windows.InvalidHandle
	t.out = windows.InvalidHandle
}

func (t *winTerm) MakeRaw() error {

	return nil
}

func (t *winTerm) MakeCooked() error {

	return nil
}
