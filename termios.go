package termios

type Terminal interface {
	Read() ([]Key, error)
	Write(string) (int, error)
	IsOpen() bool
	IsRaw() bool

	Close()
	SetRaw(bool) error
}
