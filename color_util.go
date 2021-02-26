package termios

// color_util.go converts colors to VT escape codes which can be used on Unix and newer Windows versions.

import (
	"strconv"
	"fmt"
)

func colorToEscapeCode(c *Color, isFg bool) string {
	if c.Spectrum() == SpectrumDefault {
		return defaultColorEscapeCode(isFg)
	} else if c.Spectrum() == Spectrum8 {
		return color8ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == Spectrum16 {
		return color16ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == Spectrum256 {
		return color256ToEscapeCode(c.basic, isFg)
	} else if c.Spectrum() == SpectrumRGB {
		return colorRGBToEscapeCode(c.basic, c.green, c.blue, isFg)
	} else {
		panic("not supported")
	}
}

func defaultColorEscapeCode(isFg bool) string {
	if isFg {
		return "39"
	} else {
		return "49"
	}
}

func color8ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "3" + strconv.FormatUint(uint64(value), 10)
	} else {
		return "4" + strconv.FormatUint(uint64(value), 10)
	}
}

func color16ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "9" + strconv.FormatUint(uint64(value - 8), 10)
	} else {
		return "10" + strconv.FormatUint(uint64(value - 8), 10)
	}
}

func color256ToEscapeCode(value uint8, isFg bool) string {
	if isFg {
		return "38;5;" + strconv.FormatUint(uint64(value), 10)
	} else {
		return "48;5;" + strconv.FormatUint(uint64(value), 10)
	}
}

func colorRGBToEscapeCode(r uint8, g uint8, b uint8, isFg bool) string {
	if isFg {
		return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
	} else {
		return fmt.Sprintf("48;2;%d;%d;%d", r, g, b)
	}
}
