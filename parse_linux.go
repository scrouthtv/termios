// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "os"
import "fmt"
import "unicode"
import "unicode/utf8"

// linuxParser is the default parser for linux terminals that should always work.
// It compares entered key sequences to a terminfo (either on disk or built-in)
type linuxParser struct {
	parent *nixTerm
	i      *info
}

func newParser(parent *nixTerm) (unixParser, error) {
	/*if os.Getenv("TERM") == "xterm" {
		return &xtermParser{parent}, nil
	}*/

	i, err := loadTerminfo()
	if err != nil {
		return nil, err
	}
	return &linuxParser{parent, i}, nil
}

func (p *linuxParser) open() {
	p.parent.Write(string(p.formatSimpleAction(ActionInit)))
}

func (p *linuxParser) close() {
	// FIXME: reset to the mode we were in when we first started
	p.parent.Write(string(p.formatSimpleAction(ActionExit)))
}

// ParseUTF8 splits the inputted bytes into logical keypresses
func (p *linuxParser) asKey(in []byte) []Key {
	var keys []Key

	var position int = 0
	var l int
	var k Key
	var r rune

	// What already works:
	//  - a-z, A-Z, 0-9, ext latin
	//  - symbols
	//  - C-[a-z], C-[A-Z]
	//  -

	os.Stdout.WriteString("Have to parse [ ")
	for _, b := range in {
		os.Stdout.WriteString(fmt.Sprintf("0x%x ", b))
	}
	os.Stdout.WriteString("]\r\n")

	for position < len(in) {
		if in[position] >= 0x01 && in[position] <= 0x1a {
			// C-key
			var r rune = rune(in[position]-0x01) + 'a'
			keys = append(keys, Key{KeyLetter, ModCtrl, r})
			position++
		} else if in[position] == 0x7f {
			keys = append(keys, Key{KeySpecial, 0, SpecialBackspace})
			position++
		} else if in[position] == 0x1b {

			// is escape code maybe?
			k, l = p.i.readSpecialKey(in[position:])
			os.Stdout.WriteString("It's a special key: ")
			os.Stdout.WriteString(k.String())
			os.Stdout.WriteString(fmt.Sprintf("\r\nIt's %d long\r\n", l))
			if k != InvalidKey {
				keys = append(keys, k)
				position += l
				continue
			}

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
