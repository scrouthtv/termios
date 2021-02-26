package termios

import "testing"
import "fmt"

func Test8Color(t *testing.T) {
	names := []string{ "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white" }
	values := []Color{ ColorBlack, ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorMagenta, ColorCyan, ColorWhite }

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

func Test16Color(t *testing.T) {
	names := []string{ "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white" }
	values := []Color{ ColorDarkGray, ColorLightRed, ColorLightGreen, ColorLightYellow, ColorLightBlue,
		ColorLightMagenta, ColorLightCyan, ColorLightGray }

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
