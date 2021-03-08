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
	exit()
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
	vt      actor
}

// Open opens a new terminal.
func Open() (Terminal, error) {
	var err error

	var in, out int

	// open in, out
	// TODO: https://github.com/nsf/termbox-go/blob/master/api.go opens /dev/tty as rdwr on bsd?
	in, err = unix.Open("/dev/tty", unix.O_RDONLY, 0)
	if err != nil {
		return nil, &IOError{"opening stdin", err}
	}

	out, err = unix.Open("/dev/tty", unix.O_WRONLY, 0)
	if err != nil {
		unix.Close(out)
		return nil, &IOError{"opening stdout", err}
	}

	// save old termios
	var mode *unix.Termios
	mode, err = unix.IoctlGetTermios(in, reqGetTermios)

	if err != nil {
		unix.Close(in)
		unix.Close(out)

		return nil, &IOError{"getting current termios", err}
	}

	sCh := make(chan os.Signal, 1)
	closer := make(chan bool, 1)

	signal.Notify(sCh, unix.SIGWINCH)

	var t nixTerm = nixTerm{
		in, out, true, *mode, nil, make([]byte, ioBufSize),
		TermSize{0, 0},
		sCh, closer, nil,
	}
	t.vt = &vt{&t}

	p, err := newParser(&t)
	if err != nil {
		return nil, err
	}

	t.p = p

	p.open()
	t.readSize()

	go t.signalHandler()

	return &t, nil
}

func (t *nixTerm) SetRaw(raw bool) error {
	var mode *unix.Termios
	var err error

	mode, err = unix.IoctlGetTermios(t.in, reqGetTermios)
	if err != nil {
		return &IOError{"reading termios", err}
	}

	if raw {
		// apply raw mode flags:
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

	// set termios:
	err = unix.IoctlSetTermios(t.in, reqSetTermios, mode)
	if err != nil {
		unix.Close(t.in)  //nolint:errcheck // we've already failed
		unix.Close(t.out) //nolint:errcheck // we've already failed

		return &IOError{"setting termios", err}
	}

	err = unix.IoctlSetTermios(t.out, reqSetTermios, mode)
	if err != nil {
		unix.IoctlSetTermios(t.in, reqSetTermios, mode) //nolint:errcheck // we've already failed
		unix.Close(t.in)                                //nolint:errcheck // we've already failed
		unix.Close(t.out)                               //nolint:errcheck // we've already failed

		return &IOError{"setting termios", err}
	}

	return nil
}

func (t *nixTerm) signalHandler() {
	var sig os.Signal
	var doClose bool

	for {
		select {
		case sig = <-t.sCh:
			switch sig {
			case unix.SIGWINCH:
				t.readSize()
			case unix.SIGTSTP:
				t.Close()
			}
		case doClose = <-t.closer:
			if doClose {
				return
			}
		}
	}
}

func (t *nixTerm) readSize() {
	size, err := unix.IoctlGetWinsize(t.in, unix.TIOCGWINSZ)
	if err == nil {
		t.size.Width = size.Col
		t.size.Height = size.Row
	}
}

func (t *nixTerm) GetSize() TermSize {
	return t.size
}

func (t *nixTerm) Read() ([]Key, error) {
	n, err := unix.Read(t.in, t.inBuf)
	if err != nil {
		return nil, &IOError{"reading terminal input", err}
	}

	return t.p.asKey(t.inBuf[:n]), nil
}

// SetStyle sets the specified style onto the terminal.
// On terminals that only support some colors, the closest color is approximated.
// If an unsupported color is given, the function panics.
func (t *nixTerm) SetStyle(s Style) error {
	return t.vt.setStyle(s)
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
	if !t.ready {
		return
	}

	t.closer <- true
	close(t.sCh)
	close(t.closer)

	t.ready = false

	t.p.exit()

	unix.IoctlSetTermios(t.in, reqSetTermios, &t.oldMode)  //nolint:errcheck // we're closing anyways
	unix.IoctlSetTermios(t.out, reqSetTermios, &t.oldMode) //nolint:errcheck // we're closing anyways

	unix.Close(t.in)  //nolint:errcheck // we're closing anyways
	unix.Close(t.out) //nolint:errcheck // we're closing anyways

	t.in = -1
	t.out = -1
}

func (t *nixTerm) readback(p []byte) (int, error) {
	oldmode, err := unix.IoctlGetTermios(t.in, reqGetTermios)
	if err != nil {
		return 0, &IOError{"reading terminal input", err}
	}

	rawmode := *oldmode
	rawmode.Lflag &^= (unix.ICANON | unix.ECHO)

	err = unix.IoctlSetTermios(t.in, reqSetTermios, &rawmode)
	if err != nil {
		return 0, &IOError{"setting termios", err}
	}

	err = unix.IoctlSetTermios(t.in, reqSetTermios, &rawmode)
	if err != nil {
		return 0, &IOError{"setting termios", err}
	}

	n, err := unix.Read(t.in, p)
	if err != nil {
		return 0, &IOError{"reading from terminal", err}
	}

	err = unix.IoctlSetTermios(t.in, reqSetTermios, oldmode)
	if err != nil {
		return 0, &IOError{"setting termios", err}
	}

	err = unix.IoctlSetTermios(t.in, reqSetTermios, oldmode)
	if err != nil {
		return 0, &IOError{"setting termios", err}
	}

	return n, nil
}

func newParser(parent *nixTerm) (unixParser, error) {
	if os.Getenv("TERM") == "xterm" {
		panic("todo") // TODO
	}

	if os.Getenv("TERM") == "xterm" {
		return &xtermParser{parent}, nil
	}

	i, err := loadTerminfo()
	if err != nil {
		return nil, err
	}

	return &linuxParser{parent, i}, nil
}

func (t *nixTerm) ClearScreen(c ClearType) error {
	return t.vt.clearScreen(c)
}

func (t *nixTerm) ClearLine(c ClearType) error {
	return t.vt.clearLine(c)
}

func (t *nixTerm) Move(m *Movement) error {
	return t.vt.move(m)
}

func (t *nixTerm) GetPosition() (*Position, error) {
	return t.vt.getPosition()
}
