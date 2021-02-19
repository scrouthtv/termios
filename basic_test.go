package termios

import "testing"

func TestBasicOpen(t *testing.T) {
	var term Terminal
	var err error
	term, err = Open()
	if err != nil {
		t.Error(err)
	}

	// raw mode is technically not needed on Windows
	err = term.SetRaw(true)
	if err != nil {
		t.Error(err)
	}

	//var buf []byte = make([]byte, 1024)
	var buf []Key

	for i := 0; i < 10; i++ {
		buf, err = term.Read()
		if err != nil {
			t.Error(err)
		}
		_, err = term.Write([]byte(buf[0].String() + "\r\n"))
		if err != nil {
			t.Error(err)
		}
	}
	err = term.SetRaw(false)
	if err != nil {
		t.Error(err)
	}

	term.Close()
}
