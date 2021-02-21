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
		i.specialKeys[SpecialDelete] = caps["kclr"]
		i.specialKeys[SpecialEnter] = caps["kent"]

		i.actions[ActionInit] = caps["smkx"]
		i.actions[ActionClose] = caps["rmkx"]

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
