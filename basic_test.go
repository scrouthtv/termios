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

	_, err = term.Read(buf)
	if err != nil {
		t.Error(err)
	}
	_, err = term.Write([]byte("Hello world!\n"))
	if err != nil {
		t.Error(err)
	}

	err = term.SetRaw(true)
	if err != nil {
		t.Error(err)
	}
	_, err = term.Read(buf)
	if err != nil {
		t.Error(err)
	}
	_, err = term.Write([]byte("Hello stray character!\n"))
	if err != nil {
		t.Error(err)
	}
	err = term.SetRaw(false)
	if err != nil {
		t.Error(err)
	}

	_, err = term.Read(buf)
	if err != nil {
		t.Error(err)
	}
	_, err = term.Write([]byte("Hello world!\n"))
	if err != nil {
		t.Error(err)
	}

	term.Close()
}
