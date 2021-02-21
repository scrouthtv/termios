// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "os"
import "strings"

var linuxInfo info = *newEmptyTerminfo()
var screenInfo info = *newEmptyTerminfo()
var xtermInfo info = *newEmptyTerminfo()
var urxvtInfo info = *newEmptyTerminfo()

func init() {

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
