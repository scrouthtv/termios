// +build windows

package termios

// wincon uses Windows' old Console API to do many things.
// Systems that support it should use the new VT emulation instead.
type wincon struct {
	parent *winTerm
}

func (a *wincon) setStyle(s Style) error {
	panic("todo")
	return nil
}

func (a *wincon) mapColorToWindows(c Color) uint8 {
	value := c.Downsample(Spectrum8).basic
	panic("todo")
	return 0
}

func (a *wincon) clear() error {
	panic("todo")
	return nil
}

func (a *wincon) getPosition() (*Position, error) {
	panic("todo")
	return nil, nil
}
