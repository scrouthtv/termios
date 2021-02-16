package utf8

type Key struct {
	Type  byte
	Value byte
}

const (
	// KeyLetter is a single logical key. Value is a keycode out of the Basic Latin keymap.
	KeyLetter = (1 << iota)
	// KeyCtrl is or'd to the modifier list if the ctrl key was pressed.
	// For technical reasons, C-(A-Z) is reported as C-(a-z).
	KeyCtrl
	// KeyAlt is or'd to the modifier list if the alt key was pressed.
	KeyAlt
	// KeySymbol is the Latin-1 keymap that adds many symbols.
	// 0xc2 has to be added in front of the value for printing
	KeySymbol
	// KeyUml is the Latin Extended-A keymap that adds basic umlauts.
	// 0xc3 has to be added in front of the value for printing
	KeyUml
	// KeySpecial indicates that this key should not be printed, but be interpreted instead
	KeySpecial
)

// ParseUTF8 splits the inputted bytes into logical keypresses
func ParseUTF8(in []byte) []Key {
	var i int
	var keys []Key

	for i = 0; i < len(in); i++ {
		if in[i] == 0x08 || in[i] == 0x7F {
			keys = append(keys, Key{KeySpecial, SpecialBackspace})
		} else if in[i] == 0x0D {
			keys = append(keys, Key{KeySpecial, SpecialEnter})
		} else if in[i] >= 0x01 && in[i] <= 0x1A {
			// C-key / C-Key:
			keys = append(keys, Key{KeyCtrl, in[i] - 0x01 + 0x61})
		} else if in[i] == 0xC2 {
			// parse symbol:
			i++
			keys = append(keys, Key{KeySymbol, in[i]})
		} else if in[i] == 0xC3 {
			// parse umlaut:
			i++
			keys = append(keys, Key{KeyUml, in[i]})
		} else if in[i] == 0x1B {
			i++
			if in[i] >= 0x41 && in[i] <= 0x5A || in[i] >= 0x61 && in[i] <= 0x7A {
				// A-Key / A-key:
				keys = append(keys, Key{KeyAlt, in[i]})
			}
		} else {
			keys = append(keys, Key{KeyLetter, in[i]})
		}
	}

	return keys
}
