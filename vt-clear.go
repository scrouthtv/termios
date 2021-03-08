package termios

import (
	"strconv"
	"strings"
)

func (vt *vt) clearScreen(c ClearType) error {
	var err error

	switch c {
	case ClearToEnd:
		_, err = vt.term.WriteString("\x1b[0J")
	case ClearToStart:
		_, err = vt.term.WriteString("\x1b[1J")
	case ClearCompletely:
		_, err = vt.term.WriteString("\x1b[2J")
	default:
		err = &InvalidClearTypeError{c}
	}

	if err != nil {
		return &IOError{"clearing screen", err}
	}

	return nil
}

func (vt *vt) clearLine(c ClearType) error {
	var err error

	switch c {
	case ClearToEnd:
		_, err = vt.term.WriteString("\x1b[0K")
	case ClearToStart:
		_, err = vt.term.WriteString("\x1b[1K")
	case ClearCompletely:
		_, err = vt.term.WriteString("\x1b[2K")
	default:
		err = &InvalidClearTypeError{c}
	}

	if err != nil {
		return &IOError{"clearing screen", err}
	}

	return nil
}

func (vt *vt) getPosition() (*Position, error) {
	_, err := vt.term.Write([]byte{0x1b, 0x5b, 0x36, 0x6e, 0x0b})
	if err != nil {
		return nil, &IOError{"", err}
	}

	p := make([]byte, 32)
	n, err := vt.term.readback(p)
	if err != nil { //nolint:wsl // conflicts with gofumpt
		return nil, &IOError{"reading current position", err}
	}

	if n < 6 {
		return nil, &InvalidResponseError{"reading position", string(p[:n])}
	}

	if p[0] != 0x1b || p[1] != 0x5b || p[n-1] != 0x52 {
		return nil, &InvalidResponseError{"reading position", string(p[:n])}
	}

	pos := strings.Split(string(p[2:n-1]), ";")
	if len(pos) != 2 {
		return nil, &InvalidResponseError{"reading position", string(p[:n])}
	}

	x, err := strconv.Atoi(pos[0])
	if err != nil {
		return nil, &InvalidResponseError{"reading position", string(p[:n])}
	}

	y, err := strconv.Atoi(pos[1])
	if err != nil {
		return nil, &InvalidResponseError{"reading position", string(p[:n])}
	}

	return &Position{x, y}, nil
}
