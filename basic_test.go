package termios

// basic_test.go contains a test function that can be used to debug key sequences.

import (
	"fmt"
	"testing"
	"time"
)

func TestBasicOpen(t *testing.T) {
	var term Terminal
	var err error
	term, err = Open()
	if err != nil {
		t.Error(err)
	}

	var buf []Key

	for i := 0; i < 100; i++ {
		buf, err = term.Read()
		if err != nil {
			t.Error(err)
		}
		for _, k := range buf {
			_, err = term.WriteString(k.String() + "\r\n")
			if err != nil {
				t.Error(err)
			}
			if k.Type == KeyLetter && k.Mod == 0 && k.Value == 'q' {
				term.Close()
				t.SkipNow()
			}
		}
	}

	term.Close()
}

func TestSizeReading(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 20; i++ {
		term.WriteString(fmt.Sprintf("Current size: width %d, height %d\r\n",
			term.GetSize().Width, term.GetSize().Height))
		time.Sleep(1 * time.Second)
	}
}

// FIXED in a3abd28, partly reintroduced in e62fe12:
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
