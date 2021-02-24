// +build windows

package termios

// oldActor acts on the terminal using the legacy SetConsoleTextAttribute
// console API functionality.
// Systems that support it should use the new VT emulation instead.
type oldActor struct {
	parent *winTerm
}

func (a *oldActor) setStyle(s Style) {

}

func (a *oldActor) mapColorToWindows(c Color) uint8 {
	value := c.Downsample(Spectrum8).basic

}
