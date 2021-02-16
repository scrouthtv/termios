package keys

import (
	"testing"
	"fmt"
	"os"
	"runtime"

	"github.com/scrouthtv/termios"
)

func TestDumpKeys(t *testing.T) {

	var doDump string = os.Getenv("DUMP_KEYS")
	if runtime.GOOS != "windows" && doDump != "1" {
		t.Skip("Set DUMP_KEYS to 1 to dump the pressed keys (interactive test)")
	}

	var err error
	var term termios.Terminal

	term, err = termios.Open()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	err = term.SetRaw(true)
	if err != nil {
		panic(err)
	}

	var n, i int
	var buf []byte = make([]byte, 1024)
	var b byte
	var s string
	os.Stdout.Write([]byte(fmt.Sprintf(" #  hex binary  decimal string\r\n")))

	for {
		n, err = term.Read(buf)
		for i = 0; i < n; i++ {
			if buf[i] == 0x04 {
				term.Close()
				os.Exit(0)
			}
			b = buf[i]
			s = fmt.Sprintf("%03d %03X %08b %04d %q\r\n", i, b, b, b, b)
			os.Stdout.Write([]byte(s))
		}
	}
}
