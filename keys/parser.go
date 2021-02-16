package keys

import "unicode/utf8"

type Parser struct {
	s specialParser
}

type specialParser interface {
	// ParseFirst is expected to read and return the first escape sequence and its length
	// If an error occurs, the implementation shall return InvalidKey and 1
	ParseFirst([]byte) (Key, int)
}

type Key struct {
	Type  byte
	Mod   byte
	Value rune
}

var InvalidKey = Key{KeyInvalid, 0, utf8.RuneError}

const (
	// KeyLetter is a single letter. Value is a keycode out of the Basic Latin keymap.
	KeyLetter = iota
	// KeySpecial indicates that this key should not be printed, but be interpreted instead
	KeySpecial
	// KeyInvalid is an invalid key, e. g. if an error occured during parsing
	KeyInvalid
)

const (
	// KeyCtrl is or'd to the modifier list if the ctrl key was pressed.
	// For technical reasons, C-(A-Z) is reported as C-(a-z).
	ModCtrl = (1 << iota)
	// KeyAlt is or'd to the modifier list if the alt key was pressed.
	ModAlt
)

// Init initializes the parser.
// Either the returned parser or the returned error is 0.
func Init() (*Parser, error) {
	var s specialParser
	var err error
	// newSpecialParser should initialize the special parser
	s, err = newSpecialParser()
	if err != nil {
		return nil, err
	}
	return &Parser{s}, nil
}

func parseX1B(in []byte) Key {
	return Key{0, 0, 0}
}

// ParseUTF8 splits the inputted bytes into logical keypresses
func ParseUTF8(in []byte) []Key {
	var i, n int
	var r rune
	var keys []Key

	for i = 0; i < len(in); i++ {
		if in[i] == 0x08 || in[i] == 0x7F {
			keys = append(keys, Key{KeySpecial, SpecialBackspace, utf8.RuneError})
		} else if in[i] == 0x0D {
			keys = append(keys, Key{KeySpecial, SpecialEnter, utf8.RuneError})
		} else if in[i] >= 0x01 && in[i] <= 0x1A {
			// C-key / C-Key:
			r, n = utf8.DecodeRune(in[1:])
			i += n
			keys = append(keys, Key{KeyLetter, ModCtrl, r})
		} else if in[i] == 0x1B {
			i++
			if in[i] >= 0x41 && in[i] <= 0x5A || in[i] >= 0x61 && in[i] <= 0x7A {
				// A-Key / A-key:
				r, n = utf8.DecodeRune(in[1:])
				i += n - 1
				keys = append(keys, Key{KeyLetter, ModAlt, r})
			} else if in[i] == 0x4F {
				keys = append(keys, Key{KeySpecial, SpecialF1, utf8.RuneError})
			}
		} else {
			r, n = utf8.DecodeRune(in[1:])
			i += n - 1
			keys = append(keys, Key{KeyLetter, 0, r})
		}
	}

	return keys
}
