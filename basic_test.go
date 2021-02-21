package termios

import (
	"testing"
)

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
		for _, k := range buf {
			_, err = term.Write(k.String() + "\r\n")
			if err != nil {
				t.Error(err)
			}
		}
	}
	err = term.SetRaw(false)
	if err != nil {
		t.Error(err)
	}

	term.Close()
}

// FIXED in a3abd28:
/*func TestUmls(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Error(err)
	}

	vT := term.(*winTerm)

	term.SetRaw(true)

	var s string = "aäoöÜ"
	var rs []rune = []rune(s)
	var text []uint16 = make([]uint16, len(rs))

	for i, r := range rs {
		text[i] = uint16(r)
	}

	var written uint32 = 0

	err = windows.WriteConsole(vT.out, &text[0], 5, &written, nil)
	if err != nil {
		t.Error(err)
	}
}*/
