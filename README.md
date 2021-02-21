goterm
======

Construction In Progress
------------------------

I originally adopted this library from `creack` on GitHub. The original project had to functions: setting / reading terminal size and (un-) setting raw mode.

Because of the current work over at `RythenGlyth/gosh`, we needed an API that abstracts away all OS-individual raw I/O things.

This library currently can:
 - (Un-) Set Raw mode (every single keypress is sent to the application instead of only full lines)
 - Parse a limited set of input sequences on a limited number of terminals (currently only Linux)

What is tbd:
 - [ ] Take another look at terminfo files. I trashed them during development because they're so many wrong and missing entries in every single terminfo file. Instead, I created my own built-in table which should be able to parse input on any sane terminal. Terminfo might nevertheless still be needed for writing escape sequences (movement, screen clear, ...)
 - [ ] Put everything together into a single, easy-to-use API.

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

Known issues
------------

 - Windows: A-F1 through A-F12, A-arrow, A-enter, A-escape, A-tab, C-A-Entf are not send
 - Windows/Terminal does not send many keys becauseof default keybindings: A-Enter, F11, C-Tab, C-S-Tab, S-Ins
 - Windows/Cmder does not send C-ArrowUp, C-ArrowDown
 - Windows/ConEmu does not send C-PgUp, C-PgDown by default (bound to scroll up / down)

termios
-------

This is a complete rework of creak's original termios.

I rewrote the code to provide this functionality:
 - `termios.Open()` opens the console (`stdin` for reading and `stdout` for writing) and returns a `Terminal` that has these capabilites:

 - `t.Close()` closes everything that requires closing and resets the terminal to the original state. Any encountered error is discarded.
 - `t.SetRaw(raw bool)` (un-) sets raw I/O capabilites\*, returns any encountered error.
 - `Read()` and `Write()` works according to `io` specifications
 - `IsOpen()` reports whether the programmer can currently use the terminal for I/O
 - `IsRaw()` reports whether the terminal is in raw mode

If any of these functions fail during execution, we try to undo all changes to the point before calling.

\*: *raw* means that everything that the user inputs gets transferred as that to the application. In *cooked* mode, another software layer prepares the user's input (handles `<C-c>`, waits for `\n`)

The library also exposes two file handles: `In` and `Out` which corresponds to the platform's input and output streams.

This fork makes use of Go's new, platform-specific, syscall wrappers and all "unsafe" code was removed. The entire library is just a very thin wrapper around the new `x/sys` packages which (hopefully) uses the correct constants for the correct calls.

As an example, `basic_test.go` is provided. When called, it expects the user to
1) input a line and press enter (cooked mode)
2) input a single character (raw mode)
3) again inputting a line and pressing enter (reverting to cooked mode)

The demo has been tested under
 - Windows/amd64 10, Windows 7 might not be working at all, since the console modes changed somewhere along the way of Windows 10
 - Linux/amd64 5.4 and linux/arm 5.4
 - FreeBsd/amd64 12.2, OpenBSD/amd64 6.8, NetBSD/amd64 5.8, DragonFly/amd64 5.8.3

darwin has not been tested since I can't get it to work.
From what I can tell, darwin is just BSD with a fancy-pants Kernel so it should work?

 - The new `sys` package isn't available (yet?) for Plan9, so there will be no Plan9 support for now
 - Solaris licensing is weird at the moment, so there will be no Solaris support

Interpreting byte sequences
---------------------------

In cooked I/O mode, the parent process preprocesses user input and post-processes program output.
If e. g. the user presses the backspace key, the last character is removed.
Only the full line is transmitted as a single string once the user presses enter.

In raw I/O mode, all key presses are directly transferred to the developer.
Every call to Terminal.Read() yields a byte sequence that corresponds to one or more keypresses (buffered I/O).

This has the advantage that the client application does not have to wait for the user to press enter.
However, the client application has to, in turn, interpret the raw byte sequences. For this purpose, the library `keys` has been created.

The method `keys.ParseKeys()` takes a sequence of bytes (doesn't matter if there's a single or multiple keys in it) and transforms those keys into an array of high-level data while keeping all information.

Since v2.1, the key parser has been reworked to parse normal letters using Go's runes and escape codes using termcap on Unix (local database or built-in backup) and using a built-in table on Windows.
Every pressed key has a
 - Type: either KeyLetter or KeySpecial
 - Modifier: for KeyLetter optionally ModCtrl or ModAlt, for KeySpecial one of Special\*
 - Value: for KeyLetter the full rune

Support for parsing escape sequences (cursor keys, F keys, ...) will soon be introduced.

Random things
-------------

Print all values for a capability for different terminals:
```
 ~ find /usr/share/terminfo/* -type f -printf "%f\n" | xargs -I {} infocmp {} | grep -oE "kf1=[^,]*," | sort -u
```
