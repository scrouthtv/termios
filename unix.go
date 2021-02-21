// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import (
	"golang.org/x/sys/unix"
)

const ioBufSize int = 128

type nixTerm struct {
	in      int
	out     int
	ready   bool
	isRaw   bool
	oldMode unix.Termios
	p       *linuxParser
	inBuf   []byte
}

// Open opens a new terminal for raw i/o
func Open() (Terminal, error) {
	var err error

	var p *linuxParser
	p, err = newParser()
	if err != nil {
		return nil, err
	}

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
	mode, err = unix.IoctlGetTermios(in, reqGetTermios)
	if err != nil {
		unix.Close(in)
		unix.Close(out)
		return nil, err
	}

	var t nixTerm = nixTerm{in, out, true, false, *mode, p, make([]byte, ioBufSize)}

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

func (t *nixTerm) IsRaw() bool {
	return t.isRaw
}

func (t *nixTerm) Close() {
	t.ready = false

	unix.Write(t.out, t.p.formatSimpleAction(ActionExit))

	unix.IoctlSetTermios(t.in, reqSetTermios, &t.oldMode)
	unix.IoctlSetTermios(t.out, reqSetTermios, &t.oldMode)

	unix.Close(t.in)
	unix.Close(t.out)

	t.in = -1
	t.out = -1
}

func (t *nixTerm) SetRaw(raw bool) error {
	var mode *unix.Termios
	var err error

	mode, err = unix.IoctlGetTermios(t.in, reqGetTermios)
	if err != nil {
		return err
	}

	if raw {
		mode.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
		mode.Oflag &^= unix.OPOST
		mode.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
		mode.Cflag &^= unix.CSIZE | unix.PARENB
		mode.Cflag |= unix.CS8
		mode.Cc[unix.VMIN] = 1
		mode.Cc[unix.VTIME] = 0
	} else {
		mode = &t.oldMode
	}

	err = unix.IoctlSetTermios(t.in, reqSetTermios, mode)
	if err != nil {
		return err
	}
	err = unix.IoctlSetTermios(t.out, reqSetTermios, mode)
	if err != nil {
		unix.IoctlSetTermios(t.in, reqSetTermios, mode)
		return err
	}

	if raw {
		unix.Write(t.out, t.p.formatSimpleAction(ActionInit))
	} else {
		unix.Write(t.out, t.p.formatSimpleAction(ActionExit))
	}

	t.isRaw = raw
	return nil
}
