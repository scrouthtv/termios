// +build linux freebsd netbsd openbsd dragonfly darwin

package keys

import (
	"fmt"
	"unicode/utf8"
)

type linuxParser struct {
	specialKeys map[int][][]byte
}

func newSpecialParser() (*linuxParser, error) {

	// Maps a special key to all possible representations:
	var specialKeys map[int][][]byte = make(map[int][][]byte)

	// Since the terminfo files are simply wrong,
	// I'll be using hardcoded values.
	// All values must be tested with linux, screen, xterm, urxvt, termite, eterm

	specialKeys[SpecialEnter] = [][]byte{[]byte{0x0D}}
	specialKeys[SpecialBackspace] = [][]byte{[]byte{0x7f}, []byte{0x08}}

	specialKeys[SpecialArrowUp] = [][]byte{[]byte{0x1b, 0x5b, 0x41}}
	specialKeys[SpecialArrowDown] = [][]byte{[]byte{0x1b, 0x5b, 0x42}}
	specialKeys[SpecialArrowRight] = [][]byte{[]byte{0x1b, 0x5b, 0x43}}
	specialKeys[SpecialArrowLeft] = [][]byte{[]byte{0x1b, 0x5b, 0x44}}
	specialKeys[SpecialDelete] = [][]byte{[]byte{0x1b, 0x5b, 0x33, 0x7e}}
	specialKeys[SpecialHome] = [][]byte{[]byte{0x1b, 0x5b, 0x31, 0x7e}, // linux, screen
		[]byte{0x1b, 0x5b, 0x37, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x5b, 0x48}}       // xterm, termite
	specialKeys[SpecialEnd] = [][]byte{[]byte{0x1b, 0x5b, 0x34, 0x7e}, // linux, screen
		[]byte{0x1b, 0x5b, 0x38, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x5b, 0x46}}       // xterm, termite
	specialKeys[SpecialPgUp] = [][]byte{[]byte{0x1b, 0x5b, 0x35, 0x7e}}
	specialKeys[SpecialPgDown] = [][]byte{[]byte{0x1b, 0x5b, 0x36, 0x7e}}
	specialKeys[SpecialIns] = [][]byte{[]byte{0x1b, 0x5b, 0x32, 0x7e}}

	specialKeys[SpecialF1] = [][]byte{[]byte{0x1b, 0x5b, 0x5b, 0x41}, // linux
		[]byte{0x1b, 0x5b, 0x31, 0x31, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x4f, 0x50}}             // xterm, screen, termite
	specialKeys[SpecialF2] = [][]byte{[]byte{0x1b, 0x5b, 0x5b, 0x42}, // linux
		[]byte{0x1b, 0x5b, 0x31, 0x32, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x4f, 0x51}}             // xterm, screen, termite
	specialKeys[SpecialF3] = [][]byte{[]byte{0x1b, 0x5b, 0x5b, 0x43}, // linux
		[]byte{0x1b, 0x5b, 0x31, 0x33, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x4f, 0x52}}             // xterm, screen, termite
	specialKeys[SpecialF4] = [][]byte{[]byte{0x1b, 0x5b, 0x5b, 0x44}, // linux
		[]byte{0x1b, 0x5b, 0x31, 0x34, 0x7e}, // urxvt, eterm
		[]byte{0x1b, 0x4f, 0x53}}             // xterm, screen, termite

	// Brain damage big times ahead: Starting from F5, every terminal works like xterm.
	// Also there's a jump in between F5 and F6.
	specialKeys[SpecialF5] = [][]byte{[]byte{0x1b, 0x5b, 0x5b, 0x45}, // linux
		[]byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}} // xterm, screen, termite
	specialKeys[SpecialF6] = [][]byte{[]byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}}
	specialKeys[SpecialF7] = [][]byte{[]byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}}
	specialKeys[SpecialF8] = [][]byte{[]byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}}

	// Mind the gap?
	specialKeys[SpecialF9] = [][]byte{[]byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}}
	specialKeys[SpecialF10] = [][]byte{[]byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}}
	specialKeys[SpecialF11] = [][]byte{[]byte{0x1b, 0x5b, 0x32, 0x33, 0x7e}}
	specialKeys[SpecialF12] = [][]byte{[]byte{0x1b, 0x5b, 0x32, 0x34, 0x7e}}

	return &linuxParser{specialKeys}, nil
}

func specialKeyFromSpecial(special byte) Key {
	return Key{KeySpecial, special, utf8.RuneError}
}

func (l *linuxParser) ParseFirst(in []byte) (int, int) {
	var k int
	var reps [][]byte
	var rep []byte
	var i int

	for k, reps = range l.specialKeys {
nextCap:
		for _, rep = range reps {
			if len(in) >= len(rep) {
				for i = 0; i < len(rep); i++ {
					if in[i] != rep[i] {
						continue nextCap
					}
				}
				return k, len(rep)
			}
		}
	}
	fmt.Println("No applicable sequence found")
	return -1, 1
}
