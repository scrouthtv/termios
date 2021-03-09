package termios

const (
	brightOffset = 8
)

// Color is a color of a specified spectrum.
type Color struct {
	s Spectrum
	// basic holds the value for Spectrum8, 16, 256 or the red value for rgb
	basic uint8
	// green and blue are ignored if the spectrum is not rgb
	green uint8
	blue  uint8
}

var (
	// ColorDefault resets the terminal to the default color.
	ColorDefault = Color{SpectrumDefault, 0, 0, 0}

	// ColorBlack is the black color.
	ColorBlack = Color{Spectrum8, 0, 0, 0}

	// ColorRed is the red color.
	ColorRed = Color{Spectrum8, 1, 0, 0}

	// ColorGreen is the green color.
	ColorGreen = Color{Spectrum8, 2, 0, 0}

	// ColorYellow is the yellow color, mixed from red and green.
	ColorYellow = Color{Spectrum8, 3, 0, 0}

	// ColorBlue is the blue color.
	ColorBlue = Color{Spectrum8, 4, 0, 0}

	// ColorMagenta is the magenta color, mixed from red and blue.
	ColorMagenta = Color{Spectrum8, 5, 0, 0}

	// ColorCyan is the cyan color, mixed from green and blue.
	ColorCyan = Color{Spectrum8, 6, 0, 0}

	// ColorWhite is the white color, mixed from red, green and blue.
	ColorWhite = Color{Spectrum8, 7, 0, 0}

	// ColorDarkGray is the black color.
	ColorDarkGray = Color{Spectrum16, 8, 0, 0}

	// ColorLightRed is the red color.
	ColorLightRed = Color{Spectrum16, 9, 0, 0}

	// ColorLightGreen is the green color.
	ColorLightGreen = Color{Spectrum16, 10, 0, 0}

	// ColorLightYellow is the yellow color, mixed from red and green.
	ColorLightYellow = Color{Spectrum16, 11, 0, 0}

	// ColorLightBlue is the blue color.
	ColorLightBlue = Color{Spectrum16, 12, 0, 0}

	// ColorLightMagenta is the magenta color, mixed from red and blue.
	ColorLightMagenta = Color{Spectrum16, 13, 0, 0}

	// ColorLightCyan is the cyan color, mixed from green and blue.
	ColorLightCyan = Color{Spectrum16, 14, 0, 0}

	// ColorLightGray is the white color, mixed from red, green and blue.
	ColorLightGray = Color{Spectrum16, 15, 0, 0}
)

// Spectrum returns the spectrum of a color.
func (c *Color) Spectrum() Spectrum {
	return c.s
}

// Downsample returns the closest color in the specified target spectrum.
// If target has more colors than the original color's spectrum, the old color is returned.
func (c *Color) Downsample(target Spectrum) *Color {
	if target >= c.s {
		return c
	}

	// TODO
	return nil
}
