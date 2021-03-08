// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

// terminfo.go loads all required information from a terminfo file on the drive.

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

		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialBackspace}] = caps["kbs"]
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

		i.specialKeys[Key{KeySpecial, 0, SpecialF1}] = caps["kf1"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF2}] = caps["kf2"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF3}] = caps["kf3"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF4}] = caps["kf4"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF5}] = caps["kf5"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF6}] = caps["kf6"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF7}] = caps["kf7"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF8}] = caps["kf8"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF9}] = caps["kf9"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF10}] = caps["kf10"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF11}] = caps["kf11"]
		i.specialKeys[Key{KeySpecial, 0, SpecialF12}] = caps["kf12"]

		i.specialKeys[Key{KeySpecial, ModShift, SpecialF1}] = caps["kf13"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF2}] = caps["kf14"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF3}] = caps["kf15"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF4}] = caps["kf16"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF5}] = caps["kf17"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF6}] = caps["kf18"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF7}] = caps["kf19"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF8}] = caps["kf20"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF9}] = caps["kf21"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF10}] = caps["kf22"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF11}] = caps["kf23"]
		i.specialKeys[Key{KeySpecial, ModShift, SpecialF12}] = caps["kf24"]

		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF1}] = caps["kf25"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF2}] = caps["kf26"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF3}] = caps["kf27"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF4}] = caps["kf28"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF5}] = caps["kf29"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF6}] = caps["kf30"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF7}] = caps["kf31"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF8}] = caps["kf32"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF9}] = caps["kf33"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF10}] = caps["kf34"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF11}] = caps["kf35"]
		i.specialKeys[Key{KeySpecial, ModCtrl, SpecialF12}] = caps["kf36"]

		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF1}] = caps["kf37"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF2}] = caps["kf38"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF3}] = caps["kf39"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF4}] = caps["kf40"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF5}] = caps["kf41"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF6}] = caps["kf42"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF7}] = caps["kf43"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF8}] = caps["kf44"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF9}] = caps["kf45"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF10}] = caps["kf46"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF11}] = caps["kf47"]
		i.specialKeys[Key{KeySpecial, ModCtrl | ModShift, SpecialF12}] = caps["kf48"]

		addKeys(&i.specialKeys)

		i.actions[ActionInit] = caps["smkx"]
		i.actions[ActionExit] = caps["rmkx"]

		return i, nil
	}

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
