goterm
======

I originally adopted this library from `creack` on GitHub. The original project had to functions: setting / reading terminal size and (un-) setting raw mode.

The key parsing API supports these keys on all supported terminals:
 - Letters: a-z, A-Z, 0-9, Extended Latin (U+0100 - U+FFFF)
 - Symbols: + - * # ~ , . - ; : _ < > | ^ ° ! " § $ % & / ( ) = ? { } [ ] \ ` ´
 - C-[a-z]
 - A-letter, A-Letter, A-symbol
 - F1 through F12, C-Fx, A-Fx
 - Special keys: Delete, Backspace, Enter, Insert, Home, End, PgUp, PgDown, Arrow Keys
 - C-Special, A-Special, S-Special, all combinations of these

For non-special keys, C-A-key is *explicitely* not supported and will always be replaced by key.

Supported terminals
-------------------

 - Windows:
    * Windows Terminal >= 1.6
    * ConHost >= Windows 7
    * Cmder >= 191012
    * ConEmu >= 210206
    * Fluent Terminal >= 0.7.5

Known issues
------------

 - Windows: A-F1 through A-F12, A-arrow, A-enter, A-escape, A-tab, C-A-Entf are not send
 - Windows/Terminal does not send many keys becauseof default keybindings: A-Enter, F11, C-Tab, C-S-Tab, S-Ins
 - Windows/Cmder does not send C-ArrowUp, C-ArrowDown
 - Windows/ConEmu does not send C-PgUp, C-PgDown by default (bound to scroll up / down)
 - Windows/Fluent Terminal hides a lot of keys: https://github.com/felixse/FluentTerminal/issues/885

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