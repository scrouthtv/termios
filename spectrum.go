package termios

// RGB is an rgb color with 8 bpc.
type RGB struct {
	r uint8
	g uint8
	b uint8
}

type coldef map[*Color]RGB

var coldef8 = make(coldef)
var coldef16 = make(coldef)
var coldef256 = make(coldef)

func init() {
	coldef8[&ColorBlack] = RGB{ 0, 0, 0 }
	coldef8[&ColorRed] = RGB{ 255, 0, 0 }
	coldef8[&ColorGreen] = RGB{ 0, 255, 0 }
	coldef8[&ColorYellow] = RGB{ 255, 255, 0 }
	coldef8[&ColorBlue] = RGB{ 0, 0, 255 }
	coldef8[&ColorMagenta] = RGB{ 255, 0, 255 }
	coldef8[&ColorCyan] = RGB{ 0, 255, 255 }
	coldef8[&ColorWhite] = RGB{ 255, 255, 255 }

	coldef16[&ColorBlack] = RGB{ 0, 0, 0 }
	coldef16[&ColorRed] = RGB{ 127, 0, 0 }
	coldef16[&ColorGreen] = RGB{ 0, 127, 0 }
	coldef16[&ColorYellow] = RGB{ 127, 127, 0 }
	coldef16[&ColorBlue] = RGB{ 0, 0, 127 }
	coldef16[&ColorMagenta] = RGB{ 127, 0, 127 }
	coldef16[&ColorCyan] = RGB{ 0, 127, 127 }
	coldef16[&ColorWhite] = RGB{ 170, 170, 170 }
	coldef16[&ColorDarkGray] = RGB{ 85, 85, 85 }
	coldef16[&ColorLightRed] = RGB{ 255, 0, 0 }
	coldef16[&ColorLightGreen] = RGB{ 0, 255, 0 }
	coldef16[&ColorLightYellow] = RGB{ 255, 255, 0 }
	coldef16[&ColorLightBlue] = RGB{ 0, 0, 255 }
	coldef16[&ColorLightMagenta] = RGB{ 255, 0, 255 }
	coldef16[&ColorLightCyan] = RGB{ 0, 255, 255 }
	coldef16[&ColorLightGray] = RGB{ 255, 255, 255 }


}
