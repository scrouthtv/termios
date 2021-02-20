// +build windows

package termios

type InputRecord struct {
	Type uint16
	Data [10]uint16
}

/*type KeyEventRecord struct {
	Type        uint16
	IsDown      int
	RepeatCount uint16
	VKeyCode    uint16
	VScanCode   uint16
	Unicode     uint16
	ControlKeys uint32
}

func (r *KeyEventRecord) String() string {
	var out strings.Builder

	fmt.Fprintf(&out, "Type: %d\r\n", r.Type)

	if r.IsDown > 0 {
		out.WriteString("Is down\r\n")
	} else {
		out.WriteString("Is released \r\n")
	}

	fmt.Fprintf(&out, "RepeatCount: %d\r\n", r.RepeatCount)
	fmt.Fprintf(&out, "key code: 0x%x\r\n", r.VKeyCode)
	fmt.Fprintf(&out, "unicode: %U %q\r\n", r.Unicode, r.Unicode)
	fmt.Fprintf(&out, "control keys: 0x%x", r.ControlKeys)
	fmt.Fprint(&out, "\r\n\r\n")

	return out.String()
}*/
