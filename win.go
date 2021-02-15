// +build windows

package termios

import (
	"golang.org/x/sys/windows"
)

var (
	rawInFlags []uint32 = []uint32{
		windows.ENABLE_PROCESSED_INPUT, windows.ENABLE_LINE_INPUT, windows.ENABLE_ECHO_INPUT,
	}
	rawOutFlags []uint32 = []uint32{
		windows.ENABLE_PROCESSED_OUTPUT, windows.ENABLE_WRAP_AT_EOL_OUTPUT,
	}
)

type winTerm struct {
	in windows.Handle
	out windows.Handle
	ready bool
	isRaw bool
	oldInMode uint32
	oldOutMode uint32
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
		windows.Close(in)
		return nil, err
	}

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

	var t winTerm = winTerm{in, out, true, false, inMode, outMode}

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

	windows.SetConsoleMode(t.in, t.oldInMode)
	windows.SetConsoleMode(t.out, t.oldOutMode)

	windows.Close(t.in)
	windows.Close(t.out)

	t.in = windows.InvalidHandle
	t.out = windows.InvalidHandle
}

func (t *winTerm) SetRaw(raw bool) error {
	var inMode, outMode uint32
	var err error

	err = windows.GetConsoleMode(t.in, &inMode)
	if err != nil {
		return err
	}

	err = windows.GetConsoleMode(t.out, &outMode)
	if err != nil {
		return err
	}

	// see https://docs.microsoft.com/en-us/windows/console/high-level-console-modes

	var flag uint32

	for _, flag = range rawInFlags {
		if raw {
			// unset processed input
			inMode ^= flag
		} else {
			// set processed input
			inMode |= flag
		}
	}

	for _, flag = range rawOutFlags {
		if raw {
			// unset processed output
			outMode ^= flag
		} else {
			// set processed output
			outMode |= flag
		}
	}

	err = windows.SetConsoleMode(t.in, inMode)
	if err != nil {
		return err
	}
	err = windows.SetConsoleMode(t.out, outMode)
	if err != nil {
		windows.SetConsoleMode(t.in, t.oldInMode)
		return err
	}

	t.isRaw = raw
	return nil
}
