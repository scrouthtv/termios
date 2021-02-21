package termios

// xterm supports advanced input through the altSendsEscape resource
// This library implements it like this:
// Upon initialization, enable all key modifiers:
// CSI > i m where i is 0, 1, 2, 4
// As well as altSendsEscape:
// CSI ? 1039 h
// Now Alt-<anything> will be sent as escape code:
// I don't know yet what to do about ctrl modifier.

// The documentation is very poor to say the least:
//  - Dickey has an article bashing Evans: https://invisible-island.net/xterm/modified-keys.html
//  - Evans has an article bashing Dickey: http://www.leonerd.org.uk/hacks/fixterms/

// There's a manual of 49 pages that mentions the related codes: https://invisible-island.net/xterm/ctlseqs/ctlseqs.pdf
// Key words: modifyKeyboard, modifyOtherKeys, altSendsEscape
