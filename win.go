// +build windows

package termios

import (
	"github.com/scrouthtv/termios/bwin"
	"golang.org/x/sys/windows"
)

type winTerm struct {
	in         windows.Handle
	out        windows.Handle
	ready      bool
	isRaw      bool
	oldInMode  uint32
	oldOutMode uint32
	p          *winParser
	a          actor
}

// Open opens a new terminal for raw i/o
func Open() (Terminal, error) {
	var err error

	// open I/O:
	var in, out windows.Handle
	in, err = windows.Open("CONIN$", windows.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	out, err = windows.Open("CONOUT$", windows.O_RDWR, 0)
	if err != nil {
		windows.Close(in)
		return nil, err
	}

	// store the old modes:
	var inMode, outMode uint32
	err = windows.GetConsoleMode(in, &inMode)
	if err != nil {
		windows.Close(in)
		windows.Close(out)
		return nil, err
	}

	err = windows.GetConsoleMode(out, &outMode)
	if err != nil {
		windows.Close(in)
		windows.Close(out)
		return nil, err
	}

	var t winTerm = winTerm{in, out, true, false, inMode, outMode, nil, nil}

	if t.tryVT() {
		t.a = &vt{&t}
		t.WriteString("Using vt")
	} else {
		t.a = &wincon{&t}
		t.WriteString("Using wincon")
	}

	var p *winParser
	p, err = newParser(&t)
	if err != nil {
		return nil, err
	}
	t.p = p

	return &t, nil
}

func (t *winTerm) tryVT() bool {
	m := t.oldOutMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err := windows.SetConsoleMode(t.out, m)
	return err == nil
}

// SetRaw (un-) sets the terminal's stdin & stdout to raw mode.
// Keep in mind that the winTerm.Read() functionality reads single characters
// in both modes.
func (t *winTerm) SetRaw(raw bool) error {
	var inMode, outMode, oldInMode uint32
	var err error

	windows.GetConsoleMode(t.in, &oldInMode)

	inMode = windows.ENABLE_WINDOW_INPUT      // enable input
	outMode = windows.ENABLE_PROCESSED_OUTPUT // parse new line
	err = windows.SetConsoleMode(t.in, inMode)
	if err != nil {
		return err
	}
	err = windows.SetConsoleMode(t.out, outMode)
	if err != nil {
		windows.SetConsoleMode(t.in, oldInMode)
		return err
	}

	return nil
}

func (t *winTerm) GetSize() TermSize {
	var info windows.ConsoleScreenBufferInfo
	windows.GetConsoleScreenBufferInfo(t.out, &info)
	return TermSize{Width: uint16(info.Size.X), Height: uint16(info.Size.Y)}
}

func (t *winTerm) IsOpen() bool {
	return t.ready && t.in != windows.InvalidHandle && t.out != windows.InvalidHandle
}

func (t *winTerm) Read() ([]Key, error) {
	var iR bwin.InputRecord
	var err error
	var read uint32
	var k *Key

	// wait for the first valid input:
	for k == nil {
		err = bwin.ReadConsoleInput(t.in, &iR, 1, &read)
		if err != nil {
			return nil, err
		}
		if doDebug {
			t.WriteString(iR.String())
		}
		k = t.p.asKey(iR)
	}

	return []Key{*t.p.asKey(iR)}, nil
}

func (t *winTerm) readback(p []byte) (int, error) {
	panic("todo")
	return 0, nil
}

func (t *winTerm) SetStyle(s Style) error {
	return t.a.setStyle(s)
}

// The Write method on Windows does not work well with extended latin characters.
// Use with caution.
func (t *winTerm) Write(p []byte) (int, error) {
	return windows.Write(t.out, p)
}

func (t *winTerm) WriteString(p string) (int, error) {
	var written uint32 = 0
	var err error
	var rs []rune = []rune(p)
	var code []uint16 = make([]uint16, len(rs))

	for i, r := range rs {
		code[i] = uint16(r)
	}

	err = windows.WriteConsole(t.out, &code[0], uint32(len(rs)), &written, nil)
	return int(written), err
}

func (t *winTerm) Close() {
	t.ready = false

	windows.SetConsoleMode(t.in, t.oldInMode)
	windows.SetConsoleMode(t.out, t.oldOutMode)

	windows.Close(t.in)
	windows.Close(t.out)

	t.in = windows.InvalidHandle
	t.out = windows.InvalidHandle
}

func (t *winTerm) GetPosition() (*Position, error) {
	return t.a.getPosition()
}

func (t *winTerm) Move(m *Movement) error {
	return t.a.move(m)
}

func (t *winTerm) ClearLine(c ClearType) error {
	return t.a.clearLine(c)
}

func (t *winTerm) ClearScreen(c ClearType) error {
	return t.a.clearScreen(c)
}
