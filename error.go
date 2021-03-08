package termios

import "strconv"
import "fmt"

type InvalidMovementError struct {
	value int
}

func (e *InvalidMovementError) Error() string {
	return "invalid movement " + strconv.Itoa(e.value)
}

type InvalidResponseError struct {
	id string
	resp string
}

func (e *InvalidResponseError) Error() string {
	return fmt.Sprintf("invalid reponse from terminal for %s: %q", e.id, e.resp)
}

type InvalidClearTypeError struct {
	ct ClearType
}

func (e *InvalidClearTypeError) Error() string {
	return "invalid cleartype " + strconv.FormatUint(uint64(e.ct), 10)
}
