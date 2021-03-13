// +build windows

package termios

import "github.com/scrouthtv/termios/bwin"

// wincon uses Windows' old Console API to do many things.
// Systems that support it should use the new VT emulation instead.
type wincon struct {
	parent *winTerm
}

func (a *wincon) setStyle(s Style) error {
	fg := a.mapColorToWindows(true, s.Foreground)
	bg := a.mapColorToWindows(false, s.Background)
	attr := a.mapAttrToWindows(s.Extras)
	return bwin.SetConsoleTextAttribute(a.parent.out, fg | bg | attr)
}

func (a *wincon) mapAttrToWindows(t TextAttribute) bwin.Attribute {
	var b bwin.Attribute

	if t & (TextUnderlined|TextBold) != 0 {
		b |= bwin.CommonLVBUnderscore
	}
	if t & TextReverse != 0 {
		b |= bwin.CommonLVBReverseVideo
	}

	return b
}

func (a *wincon) mapColorToWindows(isfg bool, c Color) bwin.Attribute {
	value := c.Downsample(Spectrum16).basic

	if value == ColorDefault.basic {
		return a.defaultcolor(isfg)
	} else if value >= ColorBlack.basic && value <= ColorWhite.basic {
		return a.map8ColorToWindows(isfg, value)
	} else if value >= ColorDarkGray.basic && value <= ColorLightGray.basic {
		if isfg {
			return a.map8ColorToWindows(true, value - brightOffset) | bwin.ForegroundIntensity
		} else {
			return a.map8ColorToWindows(false, value - brightOffset) | bwin.BackgroundIntensity
		}
	} else {
		return a.defaultcolor(isfg)
	}

}

func (a *wincon) defaultcolor(isfg bool) bwin.Attribute {
	if isfg {
		return bwin.ForegroundBlue | bwin.ForegroundGreen | bwin.ForegroundRed
	} else {
		return 0
	}
}


func (a *wincon) map8ColorToWindows(isfg bool, value uint8) bwin.Attribute {
	switch value {
	case ColorDefault.basic:
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
	case ColorCyan.basic:
		if isfg {
			return bwin.ForegroundGreen | bwin.ForegroundBlue
		} else {
			return bwin.BackgroundGreen | bwin.BackgroundBlue
			}
	case ColorWhite.basic:
		if isfg {
			return bwin.ForegroundRed | bwin.ForegroundGreen | bwin.ForegroundBlue
		} else {
			return bwin.BackgroundRed | bwin.BackgroundGreen | bwin.BackgroundBlue
		}
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
