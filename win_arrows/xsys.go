// +build xsys

package main

import "golang.org/x/sys/windows"

var in windows.Handle
var oldMode uint32

func SetRaw() error {
	var err error
	in, err = windows.Open("CONIN$", windows.O_RDWR, 0)
	if err != nil {
		return err
	}

	err = windows.GetConsoleMode(in, &oldMode)

	err = windows.SetConsoleMode(in, windows.ENABLE_WINDOW_INPUT | windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
	if err != nil {
		return err
	}

	return nil
}

// Read reads a single keypress
func Read(p []uint16) (int, error) {
	//return windows.Read(in, p)
	var tmp_arg uint32
	var inputControl byte = 0
	err := windows.ReadConsole(in, &p[0], 1, &tmp_arg, &inputControl)
	return 0, err
}

func Close() {
	windows.SetConsoleMode(in, oldMode)
	windows.Close(in)
}
