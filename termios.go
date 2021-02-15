package termios

type Terminal interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	IsOpen() bool
	IsRaw() bool

	Close()
	SetRaw(bool) error
}
