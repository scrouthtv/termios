package termios

import "fmt"
import "strings"
import "strconv"

func (vt *vt) move(m *Movement) error {
	var err error

	if m.flags == horizAbs && m.x == 0 && m.y == 0 {
		_, err = vt.term.Write([]byte{ 0x0d }) // CR
	} else  if m.flags == horizAbs {
		if m.x == 0 {
			if m.y > 0 {
				_, err = fmt.Fprintf(vt.term, "\x1b[%dE", m.y) // move to beginning of m.y lines down
			} else if m.y < 0 {
				_, err = fmt.Fprintf(vt.term, "\x1b[%dF", -m.y) // move to beginning of -m.y lines up
			}
		} else {
			if m.x < 0 {
				return &InvalidMovementError{m.x}
			}
			if m.y > 0 {
				_, err = fmt.Fprintf(vt.term, "\x1b[%dB", m.y) // down m.y lines
				if err != nil {
					return err
				}
			} else if m.y < 0 {
				_, err = fmt.Fprintf(vt.term, "\x1b[%dA", -m.y) // up -m.y lines
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(vt.term, "\x1b[%dG", m.x) // move to column m.x
		}
	} else if m.flags == horizAbs | vertAbs {
			_, err = fmt.Fprintf(vt.term, "\x1b[%d;%dH", m.y, m.x) // move to position m.x / m.y
	} else if m.flags == vertAbs {
		pos, err := vt.term.GetPosition()
		if err != nil {
			return err
		}
		newx := pos.X + m.x
		if newx < 0 {
			newx = 0
		}
		_, err = fmt.Fprintf(vt.term, "\x1b[%d;%dH", m.y, newx) // move to position newx / m.y
	} else if m.flags == 0 {
		if m.x > 0 {
			_, err = fmt.Fprintf(vt.term, "\x1b[%dC", m.x) // move forward by m.x
			if err != nil {
				return err
			}
		} else if m.x < 0 {
			_, err = fmt.Fprintf(vt.term, "\x1b[%dD", -m.x) // move backwards by -m.x
			if err != nil {
				return err
			}
		}

		if m.y > 0 {
			_, err = fmt.Fprintf(vt.term, "\x1b[%dB", m.y) // down m.y lines
		} else if m.y < 0 {
			_, err = fmt.Fprintf(vt.term, "\x1b[%dA", -m.y) // up -m.y liens
		}
	}

	return err
}

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

	return err
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

	return err
}

func (vt *vt) getPosition() (*Position, error) {
	p := make([]byte, 32)
	n, err := vt.term.readback(p)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	y, err := strconv.Atoi(pos[1])
	if err != nil {
		return nil, err
	}

	return &Position{x, y}, nil
}