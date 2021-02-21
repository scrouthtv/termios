// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

type linuxParser struct {
	i *info
}

func newParser() (*linuxParser, error) {
	return &linuxParser{}, nil
}

// ParseUTF8 splits the inputted bytes into logical keypresses
func (p *linuxParser) asKey(in []byte) []Key {
	var keys []Key

	return keys
}
