package termios

import "testing"

func TestBasicOpen(t *testing.T) {
	var term Terminal
	var err error
	term, err = Open()
	if err != nil {
		t.Error(err)
	}

	var buf []byte = make([]byte, 1024)
	term.Read(buf)
	term.Write([]byte("Hello world!\n"))

	term.MakeRaw()
	term.Read(buf)
	term.Write([]byte("Hello stray character!\n"))
	term.MakeCooked()

	term.Close()
}
