// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "github.com/xo/terminfo"

type info struct {
	specialKeys map[int][]byte
	actions     map[int][]byte
}

func newEmptyTerminfo() *info {
	return &info{make(map[int][]byte), make(map[int][]byte)}
}

func loadTerminfo() (*info, error) {
	var ti *terminfo.Terminfo
	var i *info
	var err error

	// First try loading the terminfo file:
	ti, err = terminfo.LoadFromEnv()

	// If that works, parse & return it:
	if err == nil {
		var caps map[string][]byte = ti.StringCapsShort()
		i = newEmptyTerminfo()

		i.specialKeys[SpecialBackspace] = caps["kbs"]
		i.specialKeys[SpecialDelete] = caps["kdch1"]
		i.specialKeys[SpecialEnter] = caps["kent"]
		i.specialKeys[SpecialArrowLeft] = caps["kcub1"]
		i.specialKeys[SpecialArrowRight] = caps["kcuf1"]
		i.specialKeys[SpecialArrowUp] = caps["kcuu1"]
		i.specialKeys[SpecialArrowDown] = caps["kcud1"]

		i.actions[ActionInit] = caps["smkx"]
		i.actions[ActionExit] = caps["rmkx"]

		return i, nil
	} else {
		// If that does not work, try to load a builtin info:
		i = loadBuiltinTerminfo()
		if i != nil {
			// If that works, return it:
			return i, nil
		} else {
			// If not, return the earlier error
			return nil, err
		}
	}
}

func (info *info) readSpecialKey(in []byte) (Key, int) {
	var s, i int
	var c []byte
	var b byte

nextSpecial:
	for s, c = range info.specialKeys {
		if len(in) >= len(c) {
			for i, b = range c {
				if in[i] != b {
					continue nextSpecial
				}
			}
			return Key{KeySpecial, 0, rune(s)}, len(c)
		}
	}
	return InvalidKey, 1
}
