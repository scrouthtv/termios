// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "github.com/xo/terminfo"

type info struct {
	specialKeys map[Key][]byte
	actions     map[int][]byte
}

func newEmptyTerminfo() *info {
	return &info{make(map[Key][]byte), make(map[int][]byte)}
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

		i.specialKeys[Key{KeySpecial, 0, SpecialBackspace}] = caps["kbs"]
		i.specialKeys[Key{KeySpecial, 0, SpecialDelete}] = caps["kdch1"]
		i.specialKeys[Key{KeySpecial, 0, SpecialEnter}] = caps["kent"]

		i.specialKeys[Key{KeySpecial, 0, SpecialArrowLeft}] = caps["kcub1"]
		i.specialKeys[Key{KeySpecial, 0, SpecialArrowRight}] = caps["kcuf1"]
		i.specialKeys[Key{KeySpecial, 0, SpecialArrowUp}] = caps["kcuu1"]
		i.specialKeys[Key{KeySpecial, 0, SpecialArrowDown}] = caps["kcud1"]

		i.specialKeys[Key{KeySpecial, 0, SpecialHome}] = caps["khome"]
		i.specialKeys[Key{KeySpecial, 0, SpecialEnd}] = caps["kend"]
		i.specialKeys[Key{KeySpecial, 0, SpecialPgUp}] = caps["kpp"]
		i.specialKeys[Key{KeySpecial, 0, SpecialPgDown}] = caps["knp"]
		i.specialKeys[Key{KeySpecial, 0, SpecialIns}] = caps["kich1"]

		addKeys(&i.specialKeys)

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
			// If not (e. g. there is no built-in for this terminal), return the first error
			return nil, err
		}
	}
}

func (info *info) readSpecialKey(in []byte) (Key, int) {
	var i int
	var k Key
	var c []byte
	var b byte

nextSpecial:
	for k, c = range info.specialKeys {
		if len(in) >= len(c) {
			for i, b = range c {
				if in[i] != b {
					continue nextSpecial
				}
			}
			return k, len(c)
		}
	}
	return InvalidKey, 1
}
