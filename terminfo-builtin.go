// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "os"
import "strings"

var linuxInfo info = *newEmptyTerminfo()
var screenInfo info = *newEmptyTerminfo()
var xtermInfo info = *newEmptyTerminfo()
var urxvtInfo info = *newEmptyTerminfo()

func init() {
	// TODO init the built-in tables
}

var extraKeys map[Key][]byte = make(map[Key][]byte)

func init() {
	extraKeys[Key{KeySpecial, ModCtrl, SpecialArrowUp}] = []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}
	extraKeys[Key{KeySpecial, ModCtrl, SpecialArrowDown}] = []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}
	extraKeys[Key{KeySpecial, ModCtrl, SpecialArrowRight}] = []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}
	extraKeys[Key{KeySpecial, ModCtrl, SpecialArrowLeft}] = []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}
}

// addKeys tries to add some keys that are equal for all terminals
// but not in the terminfo files:
//  - C-arrow
func addKeys(specialKeys *map[Key][]byte) {
	var eK Key
	var eSeq, iSeq []byte
	var ok bool

	for eK, eSeq = range extraKeys {
		iSeq, ok = (*specialKeys)[eK]
		if !ok || len(iSeq) == 0 {
			(*specialKeys)[eK] = eSeq
		}
	}
}

func loadBuiltinTerminfo() *info {
	var term string
	term = os.Getenv("TERM")
	term = strings.Split(term, "-")[0]

	switch term {
	case "linux":
		return &linuxInfo
	case "screen":
		return &screenInfo
	case "xterm", "termite":
		return &xtermInfo
	case "urxvt", "eterm":
		return &urxvtInfo
	}

	return nil
}
