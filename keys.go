package termios

import (
	"unicode/utf8"
)

// Key is an abstract key combination on the keyboard.
type Key struct {
	Type  byte
	Mod   byte
	Value rune
}

// InvalidKey is returned if an error ocured during reading.
var InvalidKey = Key{KeyInvalid, 0, utf8.RuneError}

const (
	// KeyLetter is a single letter. Value is a keycode out of the Basic Latin keymap.
	KeyLetter = iota
	// KeySpecial indicates that this key should not be printed, but be interpreted instead.
	KeySpecial
	// KeyInvalid is an invalid key, e. g. if an error ocured during parsing.
	KeyInvalid
)

const (
	// ModCtrl is or'd to the modifier list if the ctrl key was pressed.
	// For technical reasons, C-(A-Z) is reported as C-(a-z).
	ModCtrl = (1 << iota)
	// ModAlt is or'd to the modifier list if the alt key was pressed.
	ModAlt
	// ModShift is or'd to the modifier list if the shift key was pressed.
	// It is only applicable for special keys.
	ModShift
)

var (
	keyNames = []string{"letter", "special", "invalid"}
	modNames = []string{"ctrl", "alt", "shift"}
)

func (k *Key) String() string {
	if *k == InvalidKey {
		return "Invalid Key"
	}

	var s string = keyNames[k.Type] + ":"

	for mod, i := range []byte{ModCtrl, ModAlt, ModShift} {
		if (k.Mod & i) != 0 {
			s += " " + modNames[mod]
		}
	}

	if k.Type == KeyLetter {
		s += " " + string(k.Value)
	} else if k.Type == KeySpecial {
		s += " " + specialNames[k.Value]
	}

	return s
}

// Equal compares this key with another one and returns
// wether they correspond to the same combination of keys.
// Keep in mind that a key might not always have all modifiers set.
func (k *Key) Equal(other *Key) bool {
	return k.Type == other.Type && k.Mod == other.Mod && k.Value == other.Value
}
