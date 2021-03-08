// +build windows

package bwin

// Coord is a coordinate of a character cell in the console.
// The origin (0, 0) is at the top left.
type Coord struct {

	// X is the horizontal coordinate, column value or width.
	X int16

	// Y is the vertical coordinate, row value or height.
	Y int16
}
