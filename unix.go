// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import (
	"golang.org/x/sys/unix"
)

const ioBufSize int = 128

type unixParser interface {
	open()
	close()
	asKey(in []byte) []Key
}

type nixTerm struct {
	in      int
	out     int
	ready   bool
	oldMode unix.Termios
	p       unixParser
	inBuf   []byte
}

// Open opens a new terminal for raw i/o
func Open() (Terminal, error) {
	var err error

	var in, out int

	// open in, out
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

	// save old termios
	var mode *unix.Termios
	mode, err = unix.IoctlGetTermios(in, reqGetTermios)
	if err != nil {
		unix.Close(in)
		unix.Close(out)
		return nil, err
	}
	var oldMode unix.Termios = *mode

	// apply raw mode flags:
	mode.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	mode.Oflag &^= unix.OPOST
	mode.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	mode.Cflag &^= unix.CSIZE | unix.PARENB
	mode.Cflag |= unix.CS8
	mode.Cc[unix.VMIN] = 1
	mode.Cc[unix.VTIME] = 0

	// set raw termios:
	err = unix.IoctlSetTermios(in, reqSetTermios, mode)
	if err != nil {
		unix.Close(in)
		unix.Close(out)
		return nil, err
	}
	err = unix.IoctlSetTermios(out, reqSetTermios, mode)
	if err != nil {
		unix.IoctlSetTermios(in, reqSetTermios, mode)
		unix.Close(in)
		unix.Close(out)
		return nil, err
	}

	var t nixTerm = nixTerm{in, out, true, oldMode, nil, make([]byte, ioBufSize)}

	var p unixParser
	p, err = newParser(&t)
	if err != nil {
		return nil, err
	}
	t.p = p

	p.open()

	return &t, nil
}

func (t *nixTerm) Read() ([]Key, error) {
	var err error
	var n int
	n, err = unix.Read(t.in, t.inBuf)
	if err != nil {
		return nil, err
	} else {
		return t.p.asKey(t.inBuf[:n]), nil
	}
}

func (t *nixTerm) Write(s string) (int, error) {
	return unix.Write(t.out, []byte(s))
}

func (t *nixTerm) IsOpen() bool {
	return t.ready && t.in != -1 && t.out != -1
}

func (t *nixTerm) Close() {
	t.ready = false

	t.p.close()

	unix.IoctlSetTermios(t.in, reqSetTermios, &t.oldMode)
	unix.IoctlSetTermios(t.out, reqSetTermios, &t.oldMode)

	unix.Close(t.in)
	unix.Close(t.out)

	t.in = -1
	t.out = -1
}
