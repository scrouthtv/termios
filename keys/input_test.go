package keys

import "testing"
import "os"

import "github.com/scrouthtv/termios"

func TestInput(t *testing.T) {

	var doInput string = os.Getenv("DO_INPUT")
	if doInput != "1" {
		t.Skip("Set $DO_INPUT to 1 to do the interactive input test")
	}

	term, err := termios.Open()
	term.SetRaw(true)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	p, err := newSpecialParser()

	os.Stdout.Write([]byte("q to exit\n"))

	var n int

	var buf []byte = make([]byte, 1024)
	for buf[0] != 'q' {
		n, _ = term.Read(buf)
		id, _ := p.ParseFirst(buf[:n])
		if id > 0 {
			os.Stdout.Write([]byte(specialNames[id] + "\n"))
		}
	}

	term.Close()
}
