package termios

// Style collects all attributes that can be given to a character in the terminal.
// It includes foreground, background color and more text attributes.
type Style struct {
	Foreground Color
	Background Color
	Extras     TextAttribute
}

// TextAttribute sets more styling options on text.
type TextAttribute uint8

const (
	// TextDefault unsets all text attributes.
	TextDefault TextAttribute = 0

	// TextBold makes the text appear bold.
	// On some terminals it will create bright text instead.
	TextBold TextAttribute = 0x1

	// TextDim makes the text appear dim.
	TextDim TextAttribute = 0x2

	// TextUnderlined underlines the text.
	TextUnderlined TextAttribute = 0x4

	// TextBlink blinks the text. Does not work on some terminals.
	TextBlink TextAttribute = 0x8

	// TextReverse reverses foreground and background color.
	TextReverse TextAttribute = 0x10

	// TextHidden hides the text.
	TextHidden TextAttribute = 0x20

	// TextCursive prints the text in cursive. Does only work on few terminals.
	TextCursive TextAttribute = 0x40
)

// Spectrum specifies the colorspace a color is given in.
type Spectrum uint8

const (
	// SpectrumDefault only has one color: the terminal's default color.
	SpectrumDefault Spectrum = 0

	// Spectrum8 indicates that the spectrum has 8 colors.
	// It usually includes red, blue, green and binary combinations of these.
	// It needs 3 bits of storage.
	Spectrum8 Spectrum = 1

	// Spectrum16 indicates that the spectrum has 16 colors.
	// It usually adds a bright modifier to the 8 color spectrum.
	// It needs 4 bits of storage.
	Spectrum16 Spectrum = 2

	// Spectrum256 indicates that the spectrum has 256 colors.
	// It is usually a terminal-specific gradient.
	// It takes 8 bits of storage.
	Spectrum256 Spectrum = 3

	// SpectrumRGB indicates that the spectrum has all RGB colors.
	// Each component can be specified freely between 0 and 255.
	// It takes 24 bits of storage.
	SpectrumRGB Spectrum = 4
)

// MoreThan tests whether this spectrum has *more* colors than the other.
func (s *Spectrum) MoreThan(other *Spectrum) bool {
	return uint8(*s) > uint8(*other)
}

func (s *Spectrum) Equal(other *Spectrum) bool {
	return *s == *other
}
