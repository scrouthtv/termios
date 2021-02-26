termios
=======

For documentation, see https://pkg.go.dev/github.com/scrouthtv/termios.

Under construction
------------------

Tag `v3.0` showcases this libraries' functionality and already provides basic functionality for all major platforms. Still to be done:
 - Testing (different Linux terminals, bsd, darwin)
 - Linux: C-Special, A-Special, S-Special -- they neither have terminfo entries nor are identical for all terminals
 - Linux: the builtin terminfo is currently empty
 - Linux: xterm parser
 - Linux: parser#open() should return an error

Tag `v3.2` adds the `GetSize()` functionality. It's return value is implementation- and terminal-dependant. The line width should always be reported correctly, however
 - The old Windows Terminal reports the height of the underlying buffer (lines that aren't visible)
 - The new Windows Terminal reports the visible height
 - The Linux implementation is fairly consistent in that it always returns the visible size.

In `v3.3`, the `SetRaw()` functionality was reintroduced. 

Tag `v4.0` introduced the `Style` and `Color` type. Every terminal is able set its style to one that as closely as possible resembles the specified style, but may not be able to display exactly the specified style (e. g. before Windows 10 there were only 16 colors).
 - 8 and 16 color constants are defined in the `style.go`. 
 - 256 color constants are named in color256.go. They look like [this](https://jonasjacek.github.io/colors/).

The unix implementation waits for the signal SIGWINCH and reads the new window size using ioctl.

The windows implementation directly reads the current window size from the console info.
Waiting for a `WINDOW_SIZE_CHANGE` event isn't applicable as this would require the developer
reading user input every time just before the window size is required or else it'd get desynced.

I originally adopted this library from *creack* on GitHub. The original project had these functions: setting / reading terminal size and (un-) setting raw mode.

The key parsing API supports these keys on all supported terminals:
 - Letters: a-z, A-Z, 0-9, Extended Latin (U+0100 - U+FFFF)
 - Symbols: + - * # ~ , . - ; : _ < > | ^ ° ! " § $ % & / ( ) = ? { } [ ] \ ` ´
 - C-[a-z], for C-[A-Z] the lower case variant C-[a-z] should be returned
 - A-letter, A-Letter, A-symbol
 - F1 through F12, C-Fx, S-Fx, C-S-Fx
 - Special keys: Delete, Backspace, Enter, Insert, Home, End, PgUp, PgDown, Arrow Keys
 - C-Special, A-Special, S-Special

Combinations of different modifiers are only partly supported. 

For non-special keys, C-A-key is *explicitely* not supported and will always be replaced by key.

Supported terminals
-------------------

 - Windows:
    * Windows Terminal >= 1.6
    * ConHost >= Windows 7
    * Cmder >= 191012
    * ConEmu >= 210206
    * Fluent Terminal >= 0.7.5
 - Linux:
    * Windows Terminal >= 1.6 via WSL

Known issues
------------

 - Windows A-arrow, A-enter, A-escape, A-tab, C-A-Entf are not send
 - The `Write([]byte)` method does not work well with extended latin characters. Test before usage.
 - Windows/Terminal does not send many keys becauseof default keybindings: A-Enter, F11, C-Tab, C-S-Tab, S-Ins
 - Windows/Cmder does not send C-ArrowUp, C-ArrowDown
 - Windows/ConEmu does not send C-PgUp, C-PgDown by default (bound to scroll up / down)
 - Windows/Fluent Terminal hides a lot of keys: https://github.com/felixse/FluentTerminal/issues/885
 - Linux sends C-m instead of enter, C-j instead of C-enter, ...

termios
-------

This is a complete rework of creak's original termios.

I rewrote the code to provide this functionality:
 - `termios.Open()` opens the console (`stdin` for reading and `stdout` for writing) and returns a `Terminal`. The terminal is always opened in raw mode and has these capabilites:
 - `t.Close()` closes everything that requires closing and resets the terminal to the original state. Any encountered error is discarded.
 - `Read()` reads a sequence of keystrokes in the order that they were typed by the user
 - `Write()` writes a string to the terminal
 - `IsOpen()` reports whether the developer can currently use the terminal for I/O

If any of these functions fail during execution, we try to undo all changes to the point before calling.

`Read()` and `Write()` will (try to) use the underlying devices, even if the terminal isn't opened (properly) or has already been closed. 

This fork makes use of Go's new, platform-specific, syscall wrappers and all "unsafe" code was removed. The entire library is just a very thin wrapper around the new `x/sys` packages which (hopefully) uses the correct constants for the correct calls.

For testing, `basic_test.go` is provided. When called, it prints every raw data the library reads and reads 10 recognized keystrokes from the user.

 - The new `sys` package isn't available (yet?) for Plan9, so there will be no Plan9 support for now
 - Solaris licensing is weird at the moment, so there will be no Solaris support

Reading Keys
------------

Every pressed key has a
 - Type: either KeyLetter or KeySpecial
 - Modifier: for KeyLetter optionally ModCtrl or ModAlt, for KeySpecial one of Special\*
 - Value: for KeyLetter the full rune

 Implementation
 --------------

 `Terminal` is implemented by platform-specific terminal implementations (nixTerm, winTerm).
 The implementation is responsible for opening and closing as well as reading and writing byte sequences.

 Each read byte sequence is passed to an even more specific parser to be converted to a `[]Key`:
  - On Windows, bytes are compared to a built-in table (see `parse_win.go`).
  - On Linux, consoles that are compatible to `xterm`s advanced input mode (`altSendsEscape`) is parsed in `parse_xterm.go`. For other consoles, they are compared to either a terminfo file on the disk (`terminfo.go`) or a built-in terminfo table (see `terminfo_builtin.go`).
