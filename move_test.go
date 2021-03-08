package termios

import "testing"

func TestMovement(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	term.Move(MoveTo(1, 1))
	term.WriteString("asdf")
	term.Move(MoveBy(0, 0).SetColumn(0))
	term.WriteString("ASDF")
	term.Move(MoveTo(0, 0).SetDown(3))
	term.WriteString("get the camera")
	term.Move(MoveTo(0, 0).SetUp(1))
	term.WriteString("moooooooooooom")
	term.Move(MoveTo(5, 0).SetDown(4))
	term.WriteString("I'm actually hacking")
	term.Move(MoveTo(6, 0).SetUp(1))
	term.WriteString("I'm hacking")

	term.Move(MoveTo(0, 10).SetRight(5))
	term.WriteString("Somewhere")
	term.Move(MoveTo(0, 11).SetRight(5))
	term.WriteString("over the rainbow")

	term.Move(MoveTo(8, 15))
	term.WriteString("cdef")
	term.Move(MoveBy(-6, 0))
	term.WriteString("ab")
	term.Move(MoveBy(4, 0))
	term.WriteString("gh")
	term.Move(MoveBy(-8, 1))
	term.WriteString("ijklmn")
	term.Move(MoveBy(-4, -2))
	term.WriteString("uvwxyz")

	term.WriteString("\n\n\n\n")

	term.Close()
}

func TestPosition(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	term.SetRaw(true)
	t.Log(term.GetPosition())

	term.Close()
}
