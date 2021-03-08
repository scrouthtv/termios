package termios

// color_util.go converts colors to VT escape codes which can be used on Unix and newer Windows versions.

import (
	"strings"
	"strconv"
	"fmt"
)

func (vt *vt) setStyle(s Style) error {
	if s.Extras != 0 {
		panic("styles not impl")
	}

	var escape strings.Builder
	escape.WriteString("\x1b[")
	escape.WriteString(vt.colorToEscapeCode(&s.Foreground, true))
	escape.WriteString(";")
	escape.WriteString(vt.colorToEscapeCode(&s.Background, false))
	escape.WriteString("m")

	_, err := vt.term.WriteString(escape.String())
	return err
}

func (vt *vt) colorToEscapeCode(c *Color, isFg bool) string {
	if c.Spectrum() == SpectrumDefault {
		return vt.defaultColorEscapeCode(isFg)
	} else if c.Spectrum() == Spectrum8 {
		return vt.color8ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == Spectrum16 {
		return vt.color16ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == Spectrum256 {
		return vt.color256ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == SpectrumRGB {
		return vt.colorRGBToEscapeCode(c.basic, c.green, c.blue, isFg)
	} else {
		panic("not supported")
	}
}

func (vt *vt) defaultColorEscapeCode(isFg bool) string {
	if isFg {
		return "39"
	} else {
		return "49"
	}
}

func (vt *vt) color8ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "3" + strconv.FormatUint(uint64(value), 10)
	} else {
		return "4" + strconv.FormatUint(uint64(value), 10)
	}
}

func (vt *vt) color16ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "9" + strconv.FormatUint(uint64(value - 8), 10)
	} else {
		return "10" + strconv.FormatUint(uint64(value - 8), 10)
	}
}

func (vt *vt) color256ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "38;5;" + strconv.FormatUint(uint64(value), 10)
	} else {
		return "48;5;" + strconv.FormatUint(uint64(value), 10)
	}
}

func (vt *vt) colorRGBToEscapeCode(r uint8, g uint8, b uint8, isFg bool) string {
	if isFg {
		return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
	} else {
		return fmt.Sprintf("48;2;%d;%d;%d", r, g, b)
	}
}
