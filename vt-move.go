package termios

import (
	"fmt"
)

// move executes the movement by mapping it to one or more of the low-level movement
// escape codes.
func (vt *vt) move(m *Movement) error {
	var err error

	switch m.flags {
	case horizAbs | vertAbs:
		err = vt.moveTo(m.x, m.y)
	case horizAbs:
		if m.x == 0 && m.y == 0 {
			_, err = vt.term.Write([]byte{0x0d}) // CR
		} else {
			if m.x == 0 {
				err = vt.moveToStartBy(m.y)
			} else if m.x < 0 {
				return &InvalidMovementError{m.x}
			} else {
				err = vt.moveVert(m.y)
				if err != nil {
					return &IOError{"writing new cursor position", err}
				}
				err = vt.moveToColumn(m.x)
			}
		}
	case vertAbs:
		pos, perr := vt.term.GetPosition()
		if perr != nil {
			return &IOError{"reading current position", err}
		}

		newx := pos.X + m.x
		if newx < 0 {
			newx = 0
		}

		err = vt.moveTo(newx, m.y)
	case 0:
		err = vt.moveHoriz(m.x)
		if err != nil {
			return &IOError{"reading current position", err}
		}

		err = vt.moveVert(m.y)
	}

	if err != nil {
		return &IOError{"writing new position", err}
	}

	return nil
}

func (vt *vt) moveVert(y int) error {
	var err error
	if y > 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dB", y)
	} else if y < 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dA", -y)
	}

	return err //nolint:wrapcheck // internal routine
}

func (vt *vt) moveHoriz(x int) error {
	var err error
	if x > 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dC", x)
	} else if x < 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dD", -x)
	}

	return err //nolint:wrapcheck // internal routine
}

func (vt *vt) moveToColumn(x int) error {
	_, err := fmt.Fprintf(vt.term, "\x1b[%dG", x)
	return err //nolint:wrapcheck // internal routine
}

func (vt *vt) moveToStartBy(y int) error {
	var err error
	if y > 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dE", y)
	} else if y < 0 {
		_, err = fmt.Fprintf(vt.term, "\x1b[%dF", -y)
	}

	return err //nolint:wrapcheck // internal routine
}

func (vt *vt) moveTo(x, y int) error {
	_, err := fmt.Fprintf(vt.term, "\x1b[%d;%dH", y, x)
	return err //nolint:wrapcheck // internal routine
}
