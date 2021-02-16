// +build windows

package keys

import "unicode/utf8"

type winParser struct {}

func newSpecialParser() (*winParser, error) {
	return &winParser{}, nil
}

func (w *winParser) ParseFirst([]byte) (Key, int) {
	// TODO
	return Key{KeyInvalid, 0, utf8.RuneError}, 1
}
