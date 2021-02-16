// +build !windows,!linux,!freebsd,!netbsd,!openbsd,!dragonfly,!darwin

package keys

import "unicode/utf8"

type stubParser struct {}

func newSpecialParser() (*stubParser, error) {
	return &stubParser{}, nil
}

func (s *stubParser) ParseFirst([]byte) (Key, int) {
	return InvalidKey, 1
}
