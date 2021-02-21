// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "testing"

func TestSpecialParser(t *testing.T) {
	var i *info
	i, _ = loadTerminfo()

	t.Log(i)
}
