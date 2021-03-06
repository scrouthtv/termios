package termios

// color_util.go converts colors to VT escape codes which can be used on Unix and newer Windows versions.

import (
	"fmt"
	"strconv"
	"strings"
)

func (vt *vt) setStyle(s Style) error {
	var escape strings.Builder
	escape.WriteString("\x1b[0;")

	vt.writeExtras(s.Extras, &escape)

	escape.WriteString(vt.colorToEscapeCode(&s.Foreground, true))
	escape.WriteString(";")
	escape.WriteString(vt.colorToEscapeCode(&s.Background, false))
	escape.WriteString("m")

	_, err := vt.term.WriteString(escape.String())
	if err != nil {
		return &IOError{"writing color", err}
	}

	return nil
}

func (vt *vt) writeExtras(e TextAttribute, out *strings.Builder) {
	if e & TextBold != 0 {
		out.WriteString("1;")
	}
	if e & TextDim != 0 {
		out.WriteString("2;")
	}
	if e & TextUnderlined != 0 {
		out.WriteString("4;")
	}
	if e & TextBlink != 0 {
		out.WriteString("5;")
	}
	if e & TextReverse != 0 {
		out.WriteString("7;")
	}
	if e & TextHidden != 0 {
		out.WriteString("8;")
	}
	if e & TextCursive != 0 {
		out.WriteString("3;")
	}
}

func (vt *vt) colorToEscapeCode(c *Color, isFg bool) string {
	switch c.Spectrum() {
	case SpectrumDefault:
		return vt.defaultColorEscapeCode(isFg)
	case Spectrum8:
		return vt.color8ToEscapeCode(c.basic, isFg)
	case Spectrum16:
		return vt.color16ToEscapeCode(c.basic, isFg)
	case Spectrum256:
		return vt.color256ToEscapeCode(c.basic, isFg)
	case SpectrumRGB:
		return vt.colorRGBToEscapeCode(c.basic, c.green, c.blue, isFg)
	default:
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
		return "9" + strconv.FormatUint(uint64(value-brightOffset), 10)
	} else {
		return "10" + strconv.FormatUint(uint64(value-brightOffset), 10)
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
