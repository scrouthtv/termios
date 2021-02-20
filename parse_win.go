// +build windows

package termios

import (
	"fmt"
	"strings"
	"unicode"
)

type winParser struct{}

func newParser() (*winParser, error) {
	return &winParser{}, nil
}

const (
	// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	vkCancel    = 0x03
	vkBackspace = 0x08
	vkDelete    = 0x2E
	vkTab       = 0x09
	vkClear     = 0x0c
	vkReturn    = 0x0d
	vkEscape    = 0x1b
	vkPgUp      = 0x21
	vkPgDown    = 0x22
	vkEnd       = 0x23
	vkHome      = 0x24
	vkCLeft     = 0x25
	vkCUp       = 0x26
	vkCRight    = 0x27
	vkCDown     = 0x28

	vkF1  = 0x70
	vkF2  = 0x71
	vkF3  = 0x72
	vkF4  = 0x73
	vkF5  = 0x74
	vkF6  = 0x75
	vkF7  = 0x76
	vkF8  = 0x77
	vkF9  = 0x78
	vkF10 = 0x79
	vkF11 = 0x7a
	vkF12 = 0x7b
	vkF13 = 0x7c
	vkF14 = 0x7d
	vkF15 = 0x7e
	vkF16 = 0x7f
	vkF17 = 0x80
	vkF18 = 0x81
	vkF19 = 0x82
	vkF20 = 0x83
	vkF21 = 0x84
	vkF22 = 0x85
	vkF23 = 0x86
	vkF24 = 0x87
)

func (p *winParser) asKey(i InputRecord) *Key {
	if i.Type != 0x1 {
		// ignore everything but keys
		return nil
	}
	if i.Data[1] != 0x1 {
		// ignore everything but keypresses
		return nil
	}

	var mods byte = 0

	/*if i.Data[5] >= 0x01 && i.Data[5] <= 0x1A {
		FIXME
		mods |= ModCtrl
	}*/

	if i.Data[7]&0b11 > 0 {
		mods |= ModAlt
	}

	// here we use the unicode codepoint which is at position 6:
	var r rune = rune(i.Data[6])

	// here we use the virtual keyboard code which is at position 4:
	switch i.Data[4] {
	case vkCancel:
		// TODO
	case vkBackspace:
		return &Key{KeySpecial, mods, SpecialBackspace}
	case vkDelete:
		return &Key{KeySpecial, mods, SpecialDelete}
	case vkTab:
		return &Key{KeySpecial, mods, SpecialTab}
	case vkClear:
		// TODO
	case vkReturn:
		return &Key{KeySpecial, mods, SpecialEnter}
	case vkEscape:
		return &Key{KeySpecial, mods, SpecialEscape}
	case vkPgUp:
		return &Key{KeySpecial, mods, SpecialPgUp}
	case vkPgDown:
		return &Key{KeySpecial, mods, SpecialPgDown}
	case vkEnd:
		return &Key{KeySpecial, mods, SpecialEnd}
	case vkHome:
		return &Key{KeySpecial, mods, SpecialHome}
	case vkCLeft:
		return &Key{KeySpecial, mods, SpecialArrowLeft}
	case vkCUp:
		return &Key{KeySpecial, mods, SpecialArrowUp}
	case vkCRight:
		return &Key{KeySpecial, mods, SpecialArrowRight}
	case vkCDown:
		return &Key{KeySpecial, mods, SpecialArrowDown}
	case vkF1:
		return &Key{KeySpecial, mods, SpecialF1}
	case vkF2:
		return &Key{KeySpecial, mods, SpecialF2}
	case vkF3:
		return &Key{KeySpecial, mods, SpecialF3}
	case vkF4:
		return &Key{KeySpecial, mods, SpecialF4}
	case vkF5:
		return &Key{KeySpecial, mods, SpecialF5}
	case vkF6:
		return &Key{KeySpecial, mods, SpecialF6}
	case vkF7:
		return &Key{KeySpecial, mods, SpecialF7}
	case vkF8:
		return &Key{KeySpecial, mods, SpecialF8}
	case vkF9:
		return &Key{KeySpecial, mods, SpecialF9}
	case vkF10:
		return &Key{KeySpecial, mods, SpecialF10}
	case vkF11:
		return &Key{KeySpecial, mods, SpecialF11}
	case vkF12:
		return &Key{KeySpecial, mods, SpecialF12}
	default:
		//return &Key{KeyLetter, 0, 'r'}
		// TODO f keys, letter keys, umlauts?
	}

	if unicode.IsGraphic(r) {
		return &Key{KeyLetter, mods, r}
	}

	return nil
}

// FIXME: remove this method
func (r *InputRecord) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Type: %d\n\r", r.Type))
	sb.WriteString("Data: ")
	for _, b := range r.Data {
		sb.WriteString(fmt.Sprintf("0x%x ", b))
	}
	sb.WriteString("\n\r")

	return sb.String()
}
