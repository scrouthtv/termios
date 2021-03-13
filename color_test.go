package termios

import (
	"fmt"
	"testing"
)

func TestAttr(t *testing.T) {
	names := []string{"bold", "dim", "underlined", "blink", "reverse", "hidden", "cursive", "underlined blink", "reverse cursive", "bold reverse underlined blink"}
	values := []TextAttribute{ TextBold, TextDim, TextUnderlined, TextBlink, TextReverse, TextHidden, TextCursive,
			TextUnderlined | TextBlink, TextReverse | TextCursive, TextBold | TextReverse | TextUnderlined | TextBlink }

	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()

	for i, name := range names {
		term.SetStyle(Style{ColorDefault, ColorDefault, values[i]})
		term.WriteString(name)
		term.Write([]byte{ '\n' })
	}
	term.SetStyle(Style{ColorDefault, ColorDefault, 0})
}

func Test8Color(t *testing.T) {
	names := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	values := []Color{
		ColorBlack, ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorMagenta, ColorCyan, ColorWhite,
	}

	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	// HEADER:
	for _, name := range names {
		if len(name) > 10 {
			name = name[:10]
		}

		_, err = term.WriteString(fmt.Sprintf("%10s ", name))
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = term.WriteString("\r\n")
	if err != nil {
		t.Fatal(err)
	}

	for fg, fgname := range names {
		for bg, bgname := range names {
			err = term.SetStyle(Style{values[fg], values[bg], 0})
			if err != nil {
				t.Fatal(err)
			}

			if len(fgname) > 3 {
				fgname = fgname[:3]
			}

			if len(bgname) > 3 {
				bgname = bgname[:3]
			}

			_, err = term.WriteString(fmt.Sprintf("%3s on %3s", fgname, bgname))
			if err != nil {
				t.Fatal(err)
			}

			err = term.SetStyle(Style{ColorDefault, ColorDefault, 0})
			if err != nil {
				t.Fatal(err)
			}

			_, err = term.WriteString(" ")
			if err != nil {
				t.Fatal(err)
			}
		}

		_, err = term.WriteString("\r\n")
		if err != nil {
			t.Fatal(err)
		}
	}

	term.Close()
}

//nolint:errcheck // this one should only be tested after 16Color worked successfully
func Test16Color(t *testing.T) {
	names := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	values := []Color{
		ColorDarkGray, ColorLightRed, ColorLightGreen, ColorLightYellow, ColorLightBlue,
		ColorLightMagenta, ColorLightCyan, ColorLightGray,
	}

	term, err := Open()
	if err != nil {
		t.Fatal(err)
	}

	// HEADER:
	for _, name := range names {
		if len(name) > 10 {
			name = name[:10]
		}

		term.WriteString(fmt.Sprintf("%10s ", name))
	}

	term.WriteString("\r\n")

	for fg, fgname := range names {
		for bg, bgname := range names {
			term.SetStyle(Style{values[fg], values[bg], 0})

			if len(fgname) > 3 {
				fgname = fgname[:3]
			}

			if len(bgname) > 3 {
				bgname = bgname[:3]
			}

			term.WriteString(fmt.Sprintf("%3s on %3s", fgname, bgname))
			term.SetStyle(Style{ColorDefault, ColorDefault, 0})
			term.WriteString(" ")
		}

		term.WriteString("\r\n")
	}

	term.Close()
}
