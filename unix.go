// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

// unix.go implements the Terminal interface on Unix-like platforms.

import (
	"os"
	"os/signal"

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
	size    TermSize
	sCh     chan os.Signal
	closer  chan bool // send true through closer to indicate that this terminal is going down
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

	var sCh chan os.Signal = make(chan os.Signal, 1)

	signal.Notify(sCh, unix.SIGWINCH)

	var closer chan bool = make(chan bool, 1)

	var t nixTerm = nixTerm{in, out, true, oldMode, nil, make([]byte, ioBufSize),
		TermSize{0, 0}, sCh, closer}

	var p unixParser
	p, err = newParser(&t)
	if err != nil {
		return nil, err
	}
	t.p = p

	p.open()

	t.readSize()

	return &t, nil
}

func (t *nixTerm) sizeReadLoop() {
	var signal os.Signal
	var doClose bool

	for {
		select {
		case signal = <-t.sCh:
			switch signal {
			case unix.SIGWINCH:
				t.readSize()
			case unix.SIGTSTP:
				t.Close()
			}
		case doClose = <-t.closer:
			return
		}
	}
}

func (t *nixTerm) readSize() {
	size, err := unix.IoctlGetWinsize(t.in, unix.TIOCGWINSZ)
	if err == nil {
		t.size.Width = size.Cols
		t.size.Height = size.Rows
	}
}

func (t *nixTerm) GetSize() TermSize {
	return t.size
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

func (t *nixTerm) Write(p []byte) (int, error) {
	return unix.Write(t.out, p)
}

func (t *nixTerm) WriteString(s string) (int, error) {
	return unix.Write(t.out, []byte(s))
}

func (t *nixTerm) IsOpen() bool {
	return t.ready && t.in != -1 && t.out != -1
}

func (t *nixTerm) Close() {
	t.closer <- true
	close(t.sCh)
	close(t.closer)

	t.ready = false

	t.p.close()

	unix.IoctlSetTermios(t.in, reqSetTermios, &t.oldMode)
	unix.IoctlSetTermios(t.out, reqSetTermios, &t.oldMode)

	unix.Close(t.in)
	unix.Close(t.out)

	t.in = -1
	t.out = -1
}
