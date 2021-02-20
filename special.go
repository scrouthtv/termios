package termios

const (
	// SpecialBackspace is the key that deletes the character to the left of the cursor
	SpecialBackspace = iota
	// SpecialDelete is the key that deletes the character to the right of the cursor
	SpecialDelete
	SpecialEnter
	SpecialArrowLeft
	SpecialArrowRight
	SpecialArrowUp
	SpecialArrowDown
	SpecialHome
	SpecialEnd
	SpecialPgUp
	SpecialPgDown
	SpecialIns
	SpecialF1
	SpecialF2
	SpecialF3
	SpecialF4
	SpecialF5
	SpecialF6
	SpecialF7
	SpecialF8
	SpecialF9
	SpecialF10
	SpecialF11
	SpecialF12
	SpecialTab
	SpecialEscape
)

var specialNames []string = []string{
	"backspace", "delete", "enter", "cursor left", "cursor right", "cursor up", "cursor down",
	"home", "end", "page up", "page down", "insert",
	"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12",
	"tab", "escape",
}
