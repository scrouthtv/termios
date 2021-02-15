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

	term.SetRaw(true)
	term.Read(buf)
	term.Write([]byte("Hello stray character!\n"))
	term.SetRaw(false)

	term.Read(buf)
	term.Write([]byte("Hello world!\n"))

	term.Close()
}
