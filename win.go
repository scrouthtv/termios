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

func Open() (Terminal, error) {
	var err error

	var p *winParser
	p, err = newParser()
	if err != nil {
		return nil, err
	}

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

	var t winTerm = winTerm{in, out, true, false, inMode, outMode, p}

	return &t, nil
}

func (t *winTerm) IsOpen() bool {
	return t.ready && t.in != windows.InvalidHandle && t.out != windows.InvalidHandle
}

func (t *winTerm) IsRaw() bool {
	return t.isRaw
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
		windows.Write(t.out, []byte(iR.String()))
		k = t.p.asKey(iR)
	}

	return []Key{*t.p.asKey(iR)}, nil
}

func (t *winTerm) Write(p string) (int, error) {
	ae := 'ä'
	windows.Write(t.out, []byte(string(ae)))
	return windows.Write(t.out, []byte(p))
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

	if raw {
		inMode = windows.ENABLE_WINDOW_INPUT
		outMode = windows.ENABLE_PROCESSED_OUTPUT // this is needed so that new line works??????
	} else {
		inMode = t.oldInMode
		outMode = t.oldOutMode
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
