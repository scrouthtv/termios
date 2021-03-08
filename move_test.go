package termios

import "testing"

//nolint:errcheck // it's a test function
func TestClearLine(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()

	term.Move(MoveTo(0, 0).SetUp(0))
	size := term.GetSize()

	var i uint16
	for i = 0; i < size.Width; i++ {
		term.Write([]byte{'o'})
	}

	term.Move(MoveBy(-50, 0))
	term.WriteString("press enter to clear to end")
	term.Read()
	term.Move(MoveBy(30, -1))
	term.ClearLine(ClearToEnd)
	term.Read()
}

//nolint:errcheck // it's a test function
func TestClearScreen(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	term.ClearScreen(ClearCompletely)
	term.Move(MoveTo(1, 1))

	size := term.GetSize()

	var i uint16
	for i = 0; i < size.Height*size.Width; i++ {
		term.Write([]byte{'x'})
	}

	term.Move(MoveTo(int(size.Width/3), int(size.Height/2)))
	term.WriteString("Press any key to clear from here down.")
	term.Read()
	term.ClearScreen(ClearToEnd)
	term.Read()

	for i = 0; i < size.Height*size.Width; i++ {
		term.Write([]byte{'x'})
	}
	term.Move(MoveTo(int(size.Width/3), int(size.Height/2)))
	term.WriteString("Press any key to clear from here down.")
	term.Read()
	term.ClearScreen(ClearToStart)
	term.Read()

	term.Close()
}

// Output of TestMovement should look like this:
/*
ASDF

moooooooooooom
get the camera

     I'm hacking
    I'm actually hacking


          Somewhere
              over the rainbow


       uvwxyz
     abcdefgh
     ijklmn

--- PASS: TestMovement (0.03s)
*/

//nolint:errcheck // it's a test function
func TestMovement(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	term.ClearScreen(ClearCompletely)
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

	term.Move(MoveTo(0, 10).SetRight(5)) // move to line 10 and 5 to the right
	term.WriteString("Somewhere")
	term.Move(MoveTo(0, 11).SetRight(5)) // move to line 11 and 5 to the right
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

	t.Log(term.GetPosition())

	term.Close()
}
