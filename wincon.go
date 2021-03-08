// +build windows

package termios

import "github.com/scrouthtv/termios/bwin"

// wincon uses Windows' old Console API to do many things.
// Systems that support it should use the new VT emulation instead.
type wincon struct {
	parent *winTerm
}

func (a *wincon) setStyle(s Style) error {
	panic("todo") // TODO
	return nil
}

func (a *wincon) mapColorToWindows(isfg bool, c Color) bwin.Attribute {
	value := c.Downsample(Spectrum16).basic
	switch value {
	case ColorDefault.basic:
		if isfg {
			panic("todo") // TODO test is this the actual default color
			return bwin.ForegroundBlue | bwin.ForegroundGreen | bwin.ForegroundRed | bwin.ForegroundIntensity
		} else {
			panic("todo") // TODO
			return 0
		}
	case ColorBlack.basic:
		return 0
	case ColorRed.basic:
		if isfg {
			return bwin.ForegroundRed
		} else {
			return bwin.BackgroundRed
		}
	case ColorGreen.basic:
		if isfg {
			return bwin.ForegroundGreen
		} else {
			return bwin.BackgroundGreen
		}
	case ColorYellow.basic:
		if isfg {
			return bwin.ForegroundRed | bwin.ForegroundGreen
		} else {
			return bwin.BackgroundRed | bwin.BackgroundGreen
		}
	case ColorBlue.basic:
		if isfg {
			return bwin.ForegroundBlue
		} else {
			return bwin.BackgroundBlue
		}
	case ColorMagenta.basic:
		if isfg {
			return bwin.ForegroundRed | bwin.ForegroundBlue
		} else {
			return bwin.BackgroundRed | bwin.BackgroundBlue
		}
	default:
		panic("todo") // TODO
	}
	return 0
}

func (a *wincon) clear() error {
	panic("todo") // TODO
	return nil
}

func (a *wincon) getPosition() (*Position, error) {
	panic("todo") // TODO
	return nil, nil
}

func (a *wincon) clearLine(c ClearType) error {
	panic("todo") // TODO
	return nil
}

func (a *wincon) clearScreen(c ClearType) error {
	panic("todo") // TODO
	return nil
}

func (a *wincon) move(m *Movement) error {
	panic("todo") // TODO
	return nil
}
