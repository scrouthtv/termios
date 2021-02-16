// +build linux freebsd netbsd openbsd dragonfly darwin

package keys

import "unicode/utf8"

import "github.com/xo/terminfo"

type linuxParser struct {
	shortCaps map[string][]byte
}

func newSpecialParser() (*linuxParser, error) {
	var info *terminfo.Terminfo
	var err error
	info, err = terminfo.LoadFromEnv()
	if err != nil {
		return nil, err
	}
	return &linuxParser{info.StringCapsShort()}, nil
}

func specialKeyFromSpecial(special byte) Key {
	return Key{KeySpecial, special, utf8.RuneError}
}

func (l *linuxParser) ParseFirst(in []byte) (Key, int) {

	// some basic constants:
	if in[0] == 0x7F {
		return specialKeyFromSpecial(SpecialBackspace), 1
	} else if in[0] == 0x0D {
		return specialKeyFromSpecial(SpecialEnter), 1
	}

	var shortCap string
	var length int
	shortCap, length = l.getCap(in)
	switch shortCap {
	case "kbs":
		return specialKeyFromSpecial(SpecialBackspace), length
	case "kdch1":
		return specialKeyFromSpecial(SpecialDelete), length
	case "kcub1":
		return specialKeyFromSpecial(SpecialArrowLeft), length
	case "kcuf1":
		return specialKeyFromSpecial(SpecialArrowRight), length
	case "kcuu1":
		return specialKeyFromSpecial(SpecialArrowUp), length
	case "kcud1":
		return specialKeyFromSpecial(SpecialArrowDown), length
	case "khome":
		return specialKeyFromSpecial(SpecialHome), length
	case "kend":
		return specialKeyFromSpecial(SpecialEnd), length
	default:
		return InvalidKey, 1
	}
}

func (l *linuxParser) getCap(in []byte) (string, int) {
	var shortCap string
	var content []byte
	var i int

	nextCap:
	for shortCap, content = range l.shortCaps {
		if len(in) < len(content) {
			continue
		} else {
			for i = 0; i < len(content); i++ {
				if in[i] != content[i] {
					continue nextCap
				}
			}
			return shortCap, len(content)
		}
	}
	return "", 0
}
