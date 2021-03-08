package termios

import "strings"
import "fmt"

type Movement struct {
	flags uint8
	// x indicates movement in the x direction (along rows, crossing columns).
	// Positive values indicate movement to the right, negative to the left.
	// Negative values are not valid for absolute vertical movement and simply ignored.
	x int
	// y indicates movement in the y direction (along columns, crossing rows).
	// Positive values indicate movement downwards, negative upwards.
	// Negative values are not valid for absolute vertical movement and simply ignored.
	y int
}

const (
	horizAbs = 1 << iota
	vertAbs
)

func MoveTo(x, y int) *Movement {
	return &Movement{horizAbs | vertAbs, x, y}
}

func MoveBy(x, y int) *Movement {
	return &Movement{0, x, y}
}

func MoveLeft(x int) *Movement {
	return &Movement{0, -x, 0}
}

func MoveRight(x int) *Movement {
	return &Movement{0, x, 0}
}

func MoveUp(x int) *Movement {
	return &Movement{0, 0, x}
}

func MoveDown(x int) *Movement {
	return &Movement{0, 0, -x}
}

func (m *Movement) SetColumn(col int) *Movement {
	m.flags |= horizAbs
	m.x = col
	return m
}

func (m *Movement) SetDown(n int) *Movement {
	m.flags &^= vertAbs
	m.y = n
	return m
}

func (m *Movement) SetUp(n int) *Movement {
	m.flags &^= vertAbs
	m.y = -n
	return m
}

func (m *Movement) SetRight(n int) *Movement {
	m.flags &^= horizAbs
	m.x = n
	return m
}

func (m *Movement) SetLeft(n int) *Movement {
	m.flags &^= horizAbs
	m.x = -n
	return m
}

func (m *Movement) String() string {
	var out strings.Builder

	if m.flags | horizAbs != 0 {
		fmt.Fprintf(&out, "moves to column %d", m.x)
	} else {
		if m.x > 0 {
			fmt.Fprintf(&out, "moves by %d to the right", m.x)
		} else {
			fmt.Fprintf(&out, "moves by %d to the left", -m.x)
		}
	}

	if m.flags | horizAbs != 0 {
		fmt.Fprintf(&out, "and to row %d", m.y)
	} else {
		if m.y > 0 {
			fmt.Fprintf(&out, "and by %d down", m.y)
		} else {
			fmt.Fprintf(&out, "and by %d up", -m.y)
		}
	}

	return out.String()
}
