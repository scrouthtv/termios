package termios

// termios.go defines the Terminal interface and some global constants.

const doDebug = false

// Terminal is an abstract terminal where the user can press arbitrary keys
// and the developer can write arbitrary strings as well as some actions
type Terminal interface {

	// Read reads a single keypress
	// Read returns an array of keys in the order that they were typed
	// or in case of an error an empty list and the error
	// The terminal is always openend in what one might consider "raw mode"
	Read() ([]Key, error)

	// WriteString writes the specified string at the current position into the terminal
	// It returns the number of bytes (there may be multiple bytes in a character) written
	// or an error
	WriteString(string) (int, error)

	// Write writes the specified data at the current position into the terminal.
	// It returns the number of bytes written or an error.
	Write([]byte) (int, error)

	// SetRaw enables or disables raw mode for this terminal.
	SetRaw(bool) error

	// IsOpen returns whether the developer can currently read from / write to
	// this terminal.
	IsOpen() bool
	Close()

	// GetSize returns the terminal's current size.
	GetSize() TermSize
}

// TermSize groups the width and height of a terminal in characters.
type TermSize struct {
	Width  uint16
	Height uint16
}
