package termios

// basic_test.go contains a test function that can be used to debug key sequences.

import (
	"fmt"
	"testing"
	"time"
)

func TestRawCooked(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()

	for i := 0; i < 3; i++ {
		err := term.SetRaw(false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = term.WriteString("Now it should be cooked")
		if err != nil {
			t.Fatal(err)
		}

		_, err = term.Read()
		if err != nil {
			t.Fatal(err)
		}

		err = term.SetRaw(true)
		if err != nil {
			t.Fatal(err)
		}

		_, err = term.WriteString("Now it should be raw")
		if err != nil {
			t.Fatal(err)
		}

		_, err = term.Read()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicOpen(t *testing.T) {
	term, err := Open()
	if err != nil {
		t.Error(err)
	}

	var buf []Key

	term.SetRaw(true)
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
		_, err = term.WriteString(fmt.Sprintf("Current size: width %d, height %d\r\n",
			term.GetSize().Width, term.GetSize().Height))
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(1 * time.Second)
	}
}
