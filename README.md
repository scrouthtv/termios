goterm
======

This is a complete rework of creak's original termios.

I rewrote the code to provide this functionality:
 - `termios.Open()` opens the console (`stdin` for reading and `stdout` for writing) and returns a `Terminal` that has these capabilites:

 - `t.Close()` closes everything that requires closing and resets the terminal to the original state. Any encountered error is discarded.
 - `t.SetRaw(raw bool)` (un-) sets raw I/O capabilites\*, returns any encountered error.
 - `Read()` and `Write()` works according to `io` specifications
 - `IsOpen()` reports whether the programmer can currently use the terminal for I/O
 - `IsRaw()` reports whether the terminal is in raw mode

If any of these functions fail during execution, we try to undo all changes to the point before calling.

\*: *raw* means that everything that the user inputs gets transferred as that to the application. In *cooked* mode, another software layer prepares the user's input (handles `<C-c>`, waits for `\`)

The library also exposes two file handles: `In` and `Out` which corresponds to the platform's input and output streams.

This fork makes use of Go's new, platform-specific, syscall wrappers and all "unsafe" code was removed. The entire library is just a very thin wrapper around the new `x/sys` packages which (hopefully) uses the correct constants for the correct calls.

As an example, `basic_test.go` is provided. When called, it expects the user to
1) input a line and press enter (cooked mode)
2) input a single character (raw mode)
3) again inputting a line and pressing enter (reverting to cooked mode)

The demo has been tested under
 - Windows/amd64 10
 - Linux/amd64 5.4 and linux/arm 5.4
 - FreeBsd/amd64 12.2, OpenBSD/amd64 6.8, NetBSD/amd64 5.8, DragonFly/amd64 5.8.3

darwin has not been tested since I can't get it to work.
From what I can tell, darwin is just BSD with a fancy-pants Kernel so it should work?

 - The new `sys` package isn't available (yet?) for Plan9, so there will be no Plan9 support for now
 - Solaris licensing is weird at the moment, so there will be no Solaris support
