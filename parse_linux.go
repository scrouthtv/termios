// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import "os"
import "fmt"

type linuxParser struct {
	i *info
}

func newParser() (*linuxParser, error) {
	i, err := loadTerminfo()
	if err != nil {
		return nil, err
	}
	return &linuxParser{i}, nil
}

// ParseUTF8 splits the inputted bytes into logical keypresses
func (p *linuxParser) asKey(in []byte) []Key {
	var keys []Key

	var position int = 0
	var l int
	var k Key

	os.Stdout.WriteString("Have to parse [ ")
	for _, b := range in {
		os.Stdout.WriteString(fmt.Sprintf("0x%x ", b))
	}
	os.Stdout.WriteString("]\r\n")

	for position < len(in) {
		os.Stdout.WriteString(fmt.Sprintf("p_l#35: Reading @%d: %x\r\n", position, in[position]))
		if in[position] >= 0x01 && in[position] <= 0x1a {
			// C-key
			var r rune = rune(in[position]-0x01) + 'a'
			keys = append(keys, Key{KeyLetter, ModCtrl, r})
			position++
		} else if in[position] == 0x7f {
			keys = append(keys, Key{KeySpecial, 0, SpecialBackspace})
			position++
		} else if in[position] == 0x1b {
			// parse escape code:
			k, l = p.i.readSpecialKey(in[position:])
			os.Stdout.WriteString("It's a special key: ")
			os.Stdout.WriteString(k.String())
			os.Stdout.WriteString(fmt.Sprintf("\r\nIt's %d long\r\n", l))
			keys = append(keys, k)
			position += l + 1
		} else {
			position++
		}
	}

	return keys
}

func (p *linuxParser) formatSimpleAction(a int) []byte {
	switch a {
	case ActionInit:
		return p.i.actions[ActionInit]
	case ActionExit:
		return p.i.actions[ActionExit]
	default:
		return nil
	}
}
