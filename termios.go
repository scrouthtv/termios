package termios

// termios.go defines the Terminal interface and some global constants.

const doDebug = false

// Terminal is an abstract terminal where the user can press arbitrary keys
// and the developer can write arbitrary strings as well as some actions
type Terminal interface {

	// Read reads a single keypress
	// and returns an array of keys in the order that they were typed
	// or in case of an error an empty list and the error.
	// If the terminal hasn't been set to raw mode, the user must first press enter for keys
	// to be sent to the application.
	Read() ([]Key, error)

	// readback reads a single byte sequence in raw mode.
	// It is used to send escape codes to the terminal and read the answer.
	readback([]byte) (int, error)

	// WriteString writes the specified string at the current position into the terminal
	// It returns the number of bytes (there may be multiple bytes in a character) written
	// or an error.
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

	// SetStyle sets the terminal style. Not all terminals support all styles (e. g. 24bit colors).
	SetStyle(Style) error

	// GetPosition returns the current cursor position. On some terminals, this takes some time.
	GetPosition() (*Position, error)

	// Move moves the cursor and point of writing to a new position specified by the Movement.
	// Implementations will not cross line borders if the provided horizontal movement exceeds line width.
	Move(*Movement) error

	// ClearScreen clears the screen depending on the ClearType.
	ClearScreen(ClearType) error

	// ClearLine clears this line depending on the ClearType.
	ClearLine(ClearType) error
}

type ClearType uint8
const (
	ClearToEnd ClearType = iota
	ClearToStart ClearType = iota
	ClearCompletely ClearType = iota
)

type Position struct {
	X int
	Y int
}

// TermSize groups the width and height of a terminal in characters.
type TermSize struct {
	Width  uint16
	Height uint16
}

// actor does many operations on the terminal.
// It is implemented by `vt` and `wincon`.
// All unix terminals should use the `vt` implementation,
// while for Windows terminals the correct implementation
// is determined at runtime
type actor interface {
	setStyle(Style) error
	move(*Movement) error
	getPosition() (*Position, error)
	clearScreen(ClearType) error
	clearLine(ClearType) error
}
