package main

import "golang.org/x/sys/windows"
import "log"

func main() {
	in, _ := windows.Open("CONIN$", windows.O_RDWR, 0)
	out, _ := windows.Open("CONOUT$", windows.O_RDWR, 0)

	var inMode, outMode uint32
	windows.GetConsoleMode(in, &inMode)
	windows.GetConsoleMode(out, &outMode)

	oim, oom := inMode, outMode

	inMode |= windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	inMode &^= windows.ENABLE_LINE_INPUT
	windows.SetConsoleMode(in, inMode)
	windows.SetConsoleMode(out, outMode)

	windows.Write(out, []byte("\x1b[6n"))

	p := make([]uint16, 1024)
	var read uint32 = 0
	err := windows.ReadConsole(in, &(p[0]), 1, &read, nil)
	log.Println(read, err)
	log.Printf("%q", p[:read])

	windows.SetConsoleMode(in, oim)
	windows.SetConsoleMode(out, oom)

	windows.Close(in)
	windows.Close(out)
}
