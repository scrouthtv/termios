// +build linux freebsd dragonfly openbsd netbsd darwin

package termios

import (
	"golang.org/x/sys/unix"
)

type nixTerm struct {
	in int
	out int
	ready bool
	isRaw bool
	oldMode unix.Termios
}

func Open() (*nixTerm, error) {
	var err error

	var in, out int

	// TODO: https://github.com/nsf/termbox-go/blob/master/api.go opens /dev/tty as rdwr on bsd?
	in, err = unix.Open("/dev/tty", unix.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	out, err = unix.Open("/dev/tty", unix.O_WRONLY, 0)
	if err != nil {
		unix.Close(out)
		return nil, err
	}

	var mode *unix.Termios
	mode, err = unix.IoctlGetTermios(in, unix.TCGETS)
	if err != nil {
		unix.Close(in)
		unix.Close(out)
		return nil, err
	}

	var t nixTerm = nixTerm{in, out, true, false, *mode}

	return &t, nil
}

func (t *nixTerm) Read(p []byte) (int, error) {
	return unix.Read(t.in, p)
}

func (t *nixTerm) Write(p []byte) (int, error) {
	return unix.Write(t.out, p)
}

func (t *nixTerm) IsOpen() bool {
	return t.ready && t.in != -1 && t.out != -1
}

func (t *nixTerm) IsRaw() bool {
	return t.isRaw
}

func (t *nixTerm) Close() {
	t.ready = false

	unix.IoctlSetTermios(t.in, unix.TCSETS, &t.oldMode)
	unix.IoctlSetTermios(t.out, unix.TCSETS, &t.oldMode)

	unix.Close(t.in)
	unix.Close(t.out)

	t.in = -1
	t.out = -1
}

func (t *nixTerm) SetRaw(raw bool) error {
	var mode *unix.Termios
	var err error

	mode, err = unix.IoctlGetTermios(t.in, unix.TCGETS)
	if err != nil {
		return err
	}

	mode.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	mode.Oflag &^= unix.OPOST
	mode.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	mode.Cflag &^= unix.CSIZE | unix.PARENB
	mode.Cflag |= unix.CS8
	mode.Cc[unix.VMIN] = 1
	mode.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(t.in, unix.TCSETS, mode)
	if err != nil {
		return err
	}
	err = unix.IoctlSetTermios(t.out, unix.TCSETS, mode)
	return err
}
