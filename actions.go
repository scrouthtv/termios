package termios

const (
	// ActionInit has to be sent to initialize the application.
	// For xterm-like terminals, this is the `smkx` code.
	ActionInit
	// ActionExit has to be sent to close the application.
	// For xterm-like terminals, this is the `rmkx` code.
	ActionExit
)