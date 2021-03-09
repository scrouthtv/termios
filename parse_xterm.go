// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// xterm supports advanced input through the altSendsEscape resource
// This library implements it like this:
// Upon initialization, enable all modify*:
// CSI > i ; v m where i is 0, 1, 2, 4 and v is the value
// As well as altSendsEscape:
// CSI ? 1039 h
// Now Alt-<anything> will be sent as escape code:
// I don't know yet what to do about ctrl modifier.

// The documentation is very poor to say the least:
//  - Dickey has an article bashing Evans: https://invisible-island.net/xterm/modified-keys.html
//  - Evans has an article bashing Dickey: http://www.leonerd.org.uk/hacks/fixterms/

// Then there's a manual of 49 pages that has some technical details:
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.pdf
// Key words: modifyKeyboard, modifyOtherKeys, altSendsEscape

type xtermParser struct {
	parent *nixTerm
}

func (p *xtermParser) open() error {
	var s strings.Builder // buffer the opening sequence

	modify := make(map[byte][]byte, 0)

	// modifyKeyboard: modify numeric, editing, function and special keys
	// https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyKeyboard
	modify['0'] = []byte{ '1', '5' }

	// modifyCursorKeys: prefix modified sequences with csi
	// https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyCursorKeys
	modify['1'] = []byte{ '2' }

	// modifyFunctionKeys: prefix modified sequences with csi
	// https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyFunctionKeys
	modify['2'] = []byte{ '2' }

	// modifyOtherKeys: modify all other keys
	// https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
	modify['4'] = []byte{ '2' }

	for k, v := range modify {
		s.Write(append([]byte{0x1b, '[', '>', k, ';'}))
		s.Write(v)
		s.Write([]byte{ 'm' })
	}

	// Enable altSendsEscape
	s.Write([]byte{0x1b, '[', '?', '1', '0', '3', '9', 'p'})

	_, err := p.parent.WriteString(s.String())
	return err
}

func (p *xtermParser) exit() {
	var s strings.Builder

	for _, i := range []byte{'0', '1', '2', '4'} {
		s.Write([]byte{0x1b, '[', '>', i, ';', '0', 'm'})
	}

	s.Write([]byte{0x1b, '[', '?', '1', '0', '3', '9', 'l'})

	p.parent.WriteString(s.String()) //nolint:errcheck // nothing to do about it
}

func (p *xtermParser) asKey(in []byte) []Key {
	var keys []Key

	if doDebug {
		os.Stdout.WriteString("Have to parse [ ")

		var pos int
		for pos < len(in) {
			os.Stdout.WriteString(fmt.Sprintf("0x%x ", in[pos]))

			if in[pos] == 0x08 {
				keys = append(keys, Key{KeySpecial, 0, SpecialBackspace})
				pos++
			} else if in[pos] == 0x7f {
				keys = append(keys, Key{KeySpecial, ModCtrl, SpecialBackspace})
				pos++
			} else if in[pos] == 0x0d {
				keys = append(keys, Key{KeySpecial, 0, SpecialEnter})
				pos++
			} else if in[pos] >= 0x20 && in[pos] <= 0x7e {
				keys = append(keys, Key{KeyLetter, 0, rune(in[pos])})
				pos++
			} else if in[pos] == 0x1b {
				pos++
				if in[pos] == '[' { // CSI
					if in[pos+1] == '2' && in[pos+2] == '7' {
						// ESC [ 27 ; mod ; key ~
						k, i := p.parseCSI27(in[pos+3:])
						keys = append(keys, k)
						pos += i + 3
					} else if len(in) > pos+2 && in[pos+2] == '~' {
						// ESC [ key ~
						switch in[pos+1] {
						case '2':
							keys = append(keys, Key{KeySpecial, 0, SpecialIns})
						case '3':
							keys = append(keys, Key{KeySpecial, 0, SpecialDelete})
						case '5':
							keys = append(keys, Key{KeySpecial, 0, SpecialPgUp})
						case '6':
							keys = append(keys, Key{KeySpecial, 0, SpecialPgDown})
						}
						pos += 3
					} else {
						switch in[pos+1] {
						case 'A': // ESC [ A
							keys = append(keys, Key{KeySpecial, 0, SpecialArrowUp})
						case 'B': // ESC [ B
							keys = append(keys, Key{KeySpecial, 0, SpecialArrowDown})
						case 'C': // ESC [ C
							keys = append(keys, Key{KeySpecial, 0, SpecialArrowRight})
						case 'D': // ESC [ D
							keys = append(keys, Key{KeySpecial, 0, SpecialArrowLeft})
						case 'F': // ESC [ F
							keys = append(keys, Key{KeySpecial, 0, SpecialEnd})
						case 'H': // ESC [ H
							keys = append(keys, Key{KeySpecial, 0, SpecialHome})
						}
						pos += 2
					}
				}
			} else if r, l := utf8.DecodeRune(in[pos:]); unicode.IsGraphic(r) {
				keys = append(keys, Key{KeyLetter, 0, r})
				pos += l
			} else {
				keys = append(keys, InvalidKey)
				pos++
			} // TODO: [CS]*-F1-F12

		}

		// TODO: [CS]*-special

		os.Stdout.WriteString("]\r\n")
	}

	return keys
}

func (p *xtermParser) parseCSI27(buf []byte) (Key, int) {
	if len(buf) < 5 || buf[0] != ';' || buf[2] != ';' || buf[len(buf) - 1] != '~' {
		return InvalidKey, 1
	}

	//mod := buf[1]
	//key := buf[3:len(buf)-2]

	var mod rune
	var key int
	n, err := fmt.Sscanf(string(buf), ";%c;%d~", &mod, &key)
	if err != nil || n < 2 {
		return InvalidKey, 1
	}

	var kmod byte

	switch mod {
	case '2':
		kmod = ModShift
	case '3':
		kmod = ModAlt
	case '4':
		kmod = ModShift | ModAlt
	case '5':
		kmod = ModCtrl
	case '6':
		kmod = ModCtrl | ModShift
	case '7':
		kmod = 0 // C-A-key is replaecd by key
	case '8':
		kmod = ModShift // C-A-S-key is replaced by S-key
	default:
		return InvalidKey, 1
	}

	if key >= 0x20 && key <= 0x7e {
		if kmod | ModShift != 0 {
			if kmod == ModShift | ModCtrl &&  key >= 0x41 && key <= 0x5a {
				// send C-Q as C-q
				key = key - 'A' + 'a'
			}

			kmod &^= ModShift // remove shift from the modifier list as it's only allowed for special keys
		}

		// calculate the length of the whole sequence
		l := 6
		if key > 99 {
			l++
		}

		return Key{KeyLetter, kmod, rune(key)}, l
	} else if key == 0x0d {
		return Key{KeySpecial, kmod, SpecialEnter}, 6
	} else if key == 167 {
		return Key{KeyLetter, kmod, '§'}, 7 // TODO: maybe these two fixes are
	} else if key == 180 {
		return Key{KeyLetter, kmod, '´'}, 7 // specific to german keyboard layout?
	}

	return InvalidKey, 1
}
