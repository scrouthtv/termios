// +build windows

package keys

type winParser struct {}

func newSpecialParser() (*winParser, error) {
	return &winParser{}, nil
}

func (w *winParser) ParseFirst([]byte) (int, int) {
	// TODO
	return 0, 0
}
