// +build windows

package termios

// parse_win.go contains functionality for reading a Key from an InputRecord on Windows.

import (
	"fmt"
	"strings"
	"unicode"
)

type winParser struct {
	term *winTerm
}

func newParser(term *winTerm) (*winParser, error) {
	return &winParser{term}, nil
}

var vkCodes map[uint16]byte = make(map[uint16]byte)

func init() {
	// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	// they also specify F13 through F24, however I can't seem to trigger these keys
	// C-F1, A-F1, C-A-F1 does not work
	// Maybe they were special keys at some ancient point in time
	// Same with vkCancel and vkClear
	// If you own a keyboard with these keys feel free to @ me
	vkCodes[0x08] = SpecialBackspace
	vkCodes[0x09] = SpecialTab
	vkCodes[0x0d] = SpecialEnter
	vkCodes[0x1b] = SpecialEscape
	vkCodes[0x21] = SpecialPgUp
	vkCodes[0x22] = SpecialPgDown
	vkCodes[0x23] = SpecialEnd
	vkCodes[0x24] = SpecialHome
	vkCodes[0x25] = SpecialArrowLeft
	vkCodes[0x26] = SpecialArrowUp
	vkCodes[0x27] = SpecialArrowRight
	vkCodes[0x28] = SpecialArrowDown
	vkCodes[0x2d] = SpecialIns
	vkCodes[0x2e] = SpecialDelete

	vkCodes[0x70] = SpecialF1
	vkCodes[0x71] = SpecialF2
	vkCodes[0x72] = SpecialF3
	vkCodes[0x73] = SpecialF4
	vkCodes[0x74] = SpecialF5
	vkCodes[0x75] = SpecialF6
	vkCodes[0x76] = SpecialF7
	vkCodes[0x77] = SpecialF8
	vkCodes[0x78] = SpecialF9
	vkCodes[0x79] = SpecialF10
	vkCodes[0x7a] = SpecialF11
	vkCodes[0x7b] = SpecialF12
}

func (p *winParser) asKey(i InputRecord) *Key {
	if i.Type != 0x1 {
		// ignore everything but keys
		// I tested using the WindowBufferSizeChange event, however it gets not sent
		// until the user resizes the window at least once. That's why I'm using
		// GetConsoleScreenBufferInfo instead.
		// TODO: maybe it'd be better to read the first size with GetConsoleScreenBufferInfo
		// store that value and change it whenever we receive 0x4
		return nil
	}
	if i.Data[1] != 0x1 {
		// ignore everything but keypresses
		return nil
	}

	var mods byte = 0

	// Detecting modifiers:
	if i.Data[7]&0b11 > 0 {
		mods |= ModAlt
	}
	if i.Data[7]&0b1100 > 0 {
		mods |= ModCtrl
	}

	// First some random stuff that does not work out of the box:
	if i.Data[4] == 0x08 && i.Data[5] == 0x0e {
		if i.Data[7]&0b10000 > 0 {
			mods |= ModShift
		}
		return &Key{KeySpecial, mods, SpecialBackspace} // overlaps with C-h
	}
	if i.Data[4] == 0x0d && i.Data[5] == 0x1c {
		if i.Data[7]&0b10000 > 0 {
			mods |= ModShift
		}
		return &Key{KeySpecial, mods, SpecialEnter} // overlaps with C-m
	}
	if i.Data[4] == 0x09 && i.Data[5] == 0x0f {
		if i.Data[7]&0b10000 > 0 {
			mods |= ModShift
		}
		return &Key{KeySpecial, mods, SpecialTab} // overlaps with C-i
	}

	if i.Data[6] >= 0x01 && i.Data[6] <= 0x1a {
		// C-key
		var r rune = rune(i.Data[6]-0x01) + 'a'
		return &Key{KeyLetter, ModCtrl, r}
	}

	// here we use the virtual keyboard code which is at position 4:
	special, ok := vkCodes[i.Data[4]]
	if ok {
		if i.Data[7]&0b10000 > 0 {
			mods |= ModShift
		}
		return &Key{KeySpecial, mods, rune(special)}
	}

	// Because some things are inputted with C-A-*:
	if mods == ModCtrl|ModAlt {
		mods = 0
	}

	// here we use the unicode codepoint which is at position 6:
	var r rune = rune(i.Data[6])

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
