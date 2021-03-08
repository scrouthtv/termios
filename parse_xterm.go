// +build linux freebsd openbsd netbsd dragonfly darwin

package termios

import (
	"fmt"
	"os"
	"strings"
)

// xterm supports advanced input through the altSendsEscape resource
// This library implements it like this:
// Upon initialization, enable all modify*:
// CSI > i m where i is 0, 1, 2, 4
// As well as altSendsEscape:
// CSI ? 1039 h
// Now Alt-<anything> will be sent as escape code:
// I don't know yet what to do about ctrl modifier.

// The documentation is very poor to say the least:
//  - Dickey has an article bashing Evans: https://invisible-island.net/xterm/modified-keys.html
//  - Evans has an article bashing Dickey: http://www.leonerd.org.uk/hacks/fixterms/

// Then there's a manual of 49 pages that has some technical details:
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.pdf
// Key words: modifyKeyboard, modifyOtherKeys, altSendsEscape

type xtermParser struct {
	parent *nixTerm
}

func (p *xtermParser) open() {
	var s strings.Builder // buffer the opening sequence

	for _, i := range []byte{'0', '1', '2', '4'} {
		s.Write([]byte{0x1b, '[', '>', i, ';', '1', 'm'})
	}

	s.Write([]byte{0x1b, '[', '?', '1', '0', '3', '9', 'h'})

	p.parent.WriteString(s.String())
}

func (p *xtermParser) exit() {
	var s strings.Builder

	for _, i := range []byte{'0', '1', '2', '4'} {
		s.Write([]byte{0x1b, '[', '>', i, ';', '0', 'm'})
	}

	s.Write([]byte{0x1b, '[', '?', '1', '0', '3', '9', 'l'})

	p.parent.WriteString(s.String())
}

func (p *xtermParser) asKey(in []byte) []Key {
	var keys []Key

	if doDebug {
		os.Stdout.WriteString("Have to parse [ ")

		for _, b := range in {
			os.Stdout.WriteString(fmt.Sprintf("0x%x ", b))
		}

		os.Stdout.WriteString("]\r\n")
	}

	panic("xterm parser not yet implemented")

	return keys
}
