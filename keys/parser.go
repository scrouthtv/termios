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

// ParseUTF8 splits the inputted bytes into logical keypresses
func (p *Parser) ParseUTF8(in []byte) []Key {
	var i, n int
	var r rune
	var key Key
	var keys []Key

	var runs int

	for i = 0; i < len(in); i++ {
		runs++
		if runs > 10000 {
			panic("endless loop")
		}
		if in[i] >= 0x01 && in[i] <= 0x1A {
			// C-key / C-Key:
			// FIXME
			keys = append(keys, Key{KeyLetter, ModCtrl, r})
		} else if in[i] == 0x1B {
			if (in[i+1] >= 0x41 && in[i+1] <= 0x5A) || (in[i+1] >= 0x61 && in[i+1] <= 0x7A) {
				// A-Key / A-key: decode remaining characters using utf8 library
				r, n = utf8.DecodeRune(in[i+1:])
				i += n - 1
				keys = append(keys, Key{KeyLetter, ModAlt, r})
			} else {
				key, n = p.s.ParseFirst(in[i+1:])
				i += n - 1
				keys = append(keys, key)
			}
		} else {
			r, n = utf8.DecodeRune(in[i+1:])
			i += n - 1
			keys = append(keys, Key{KeyLetter, 0, r})
		}
	}

	return keys
}
