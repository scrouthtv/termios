package bwin

import "testing"

import "golang.org/x/sys/windows"

func TestSetAttribute(t *testing.T) {
	in, err := windows.Open("CONIN$", windows.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	out, err := windows.Open("CONOUT$", windows.O_RDWR, 0)
	if err != nil {
		windows.Close(in)
		t.Fatal(err)
	}

	contents := uint16('q')
	cB := []byte{'p'}

	for i := 0; i < 8; i++ {
		err := SetConsoleTextAttribute(out, Attribute(i))
		if err != nil {
			t.Fatal(err)
		}
		windows.Write(out, cB)
		windows.WriteConsole(out, &contents, 1, nil, nil)
	}
	windows.Write(out, []byte{'\n'})
}
