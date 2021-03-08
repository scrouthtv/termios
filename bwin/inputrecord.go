// +build windows

package bwin

import "fmt"
import "strings"

type InputRecord struct {
	Type uint16
	Data [8]uint16
}

func (r *InputRecord) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Type: %d\n\r", r.Type))
	sb.WriteString("Data: ")
	for _, b := range r.Data {
		sb.WriteString(fmt.Sprintf("0x%x ", b))
	}
	sb.WriteString("\n\r")

	return sb.String()
}
