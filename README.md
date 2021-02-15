goterm
======

This is a complete rework of creak's original termios.

I rewrote the code to provide four functions:
 - `termios.Open()` opens a console (`stdin` for reading and `stdout` for writing)
 - `termios.Close()` closes everything that requires closing
 - `termios.MakeRaw()` sets raw I/O capabilites\*
 - `termios.MakeCooked()` unsets raw I/O capabilites

If any of these functions fail during execution, we try to undo all changes to the point before calling.

\*: *raw* means that everything that the user inputs gets transferred as that to the application. In *cooked* mode, another software layer prepares the user's input (handles `<C-c>`, waits for `\`)

The library also exposes two file handles: `In` and `Out` which corresponds to the platform's input and output streams.

This fork also makes use of Go's new, platform-specific, syscall wrappers and all "unsafe" code should be removed.
