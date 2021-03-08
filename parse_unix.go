// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

// parse_linux.go contains functionality for reading a Key from a []byte input on Unix.

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

// linuxParser is the default parser for linux terminals that should always work.
// It compares entered key sequences to a terminfo (either on disk or built-in).
type linuxParser struct {
	parent *nixTerm
	i      *info
}

func (p *linuxParser) open() {
	p.parent.Write(p.formatSimpleAction(ActionInit)) //nolint:errcheck // nothing to do about it
}

func (p *linuxParser) exit() {
	// FIXME: reset to the mode we were in when we first started
	p.parent.Write(p.formatSimpleAction(ActionExit)) //nolint:errcheck // nothing to do about it
}

// ParseUTF8 splits the inputted bytes into logical keypresses.
func (p *linuxParser) asKey(in []byte) []Key {
	var keys []Key

	var position int
	var l int
	var k Key
	var r rune

	// Here we look for:
	//  - a-z, A-Z, 0-9, ext latin
	//  - symbols
	//  - C-[a-z], C-[A-Z]
	//  - A-letter, A-Letter, A-symbol
	// Special keys starting with x1b are delegated to the info implementation

	if doDebug {
		os.Stdout.WriteString("Have to parse [ ")

		for _, b := range in {
			os.Stdout.WriteString(fmt.Sprintf("0x%x ", b))
		}

		os.Stdout.WriteString("]\r\n")
	}

	for position < len(in) {
		// is escape code maybe?
		k, l = p.i.readSpecialKey(in[position:])

		if doDebug {
			os.Stdout.WriteString("It's a special key: ")
			os.Stdout.WriteString(k.String())
			os.Stdout.WriteString(fmt.Sprintf("\r\nIt's %d long\r\n", l))
		}

		if k != InvalidKey {
			keys = append(keys, k)
			position += l

			continue
		}

		if in[position] == 0x7f { // somehow, they always get this wrong
			keys = append(keys, Key{KeySpecial, 0, SpecialBackspace})
			position++
		} else if in[position] == 0xd { // this one as well
			keys = append(keys, Key{KeySpecial, 0, SpecialEnter})
			position++
		} else if in[position] >= 0x01 && in[position] <= 0x1a {
			// C-key
			r := rune(in[position]-0x01) + 'a'
			keys = append(keys, Key{KeyLetter, ModCtrl, r})
			position++
		} else if in[position] == 0x1b {
			// Else try A-Letter, A-letter, A-symbol
			// TODO in which terminals does this work and why???
			keys = append(keys, Key{KeyLetter, ModAlt, rune(in[position+1])})
			position += 2
		} else if r, l = utf8.DecodeRune(in[position:]); unicode.IsGraphic(r) {
			keys = append(keys, Key{KeyLetter, 0, r})
			position += l
		}
	}

	return keys
}

func (p *linuxParser) formatSimpleAction(a int) []byte {
	switch a {
	case ActionInit:
		return p.i.actions[ActionInit]
	case ActionExit:
		return p.i.actions[ActionExit]
	default:
		return nil
	}
}
