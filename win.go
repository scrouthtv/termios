// +build windows

package termios

import (
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
}

// Open opens a new terminal for raw i/o
func Open() (Terminal, error) {
	var err error

	var p *winParser
	p, err = newParser()
	if err != nil {
		return nil, err
	}

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

	var oldInMode, oldOutMode uint32 = inMode, outMode

	// open as raw:
	inMode = windows.ENABLE_WINDOW_INPUT
	outMode = windows.ENABLE_PROCESSED_OUTPUT // parse new line
	err = windows.SetConsoleMode(in, inMode)
	if err != nil {
		windows.Close(in)
		windows.Close(out)
		return nil, err
	}
	err = windows.SetConsoleMode(out, outMode)
	if err != nil {
		windows.SetConsoleMode(in, oldInMode)
		windows.Close(in)
		windows.Close(out)
		return nil, err
	}

	var t winTerm = winTerm{in, out, true, false, oldInMode, oldOutMode, p}

	return &t, nil
}

func (t *winTerm) IsOpen() bool {
	return t.ready && t.in != windows.InvalidHandle && t.out != windows.InvalidHandle
}

func (t *winTerm) Read() ([]Key, error) {
	var iR InputRecord
	var err error
	var read uint32
	var k *Key

	// wait for the first valid input:
	for k == nil {
		err = ReadConsoleInput(t.in, &iR, 1, &read)
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
