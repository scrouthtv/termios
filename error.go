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
