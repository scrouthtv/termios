// +build windows

package bwin

// Attribute are character attributes that can be set on a CharInfo.
// They specify foreground, background colors and additional styles.
type Attribute uint16

// see https://docs.microsoft.com/en-us/windows/console/char-info-str
// Color Attributes can be mixed to achieve different colors:
// 0x0 is black, 0xF is white,
// blue | green | intense is bright cyan.
const (

	// ForegroundBlue indicates that the text color should be blue.
	ForegroundBlue Attribute = 0x0001

	// ForegroundGreen indicates that the text color should be green.
	ForegroundGreen Attribute = 0x0002

	// ForegroundRed indicates that the text color should be red.
	ForegroundRed Attribute = 0x0004

	// ForegroundIntensity indicates that the character should be intensified.
	ForegroundIntensity Attribute = 0x0008

	// BackgroundBlue indicates that the background color should be blue.
	BackgroundBlue Attribute = 0x0010

	// BackgroundGreen indicates that the background color should be green.
	BackgroundGreen Attribute = 0x0020

	// BackgroundRed indicates that the background color should be red.
	BackgroundRed Attribute = 0x0040

	// BackgroundIntensity indicates that the background should be intensified.
	BackgroundIntensity Attribute = 0x0080

	// CommonLVBLeadingByte indicates the leading byte of a CJK character group.
	CommonLVBLeadingByte Attribute = 0x0100

	// CommonLVBTrailingByte indicates the trailing byte of a CJK character group.
	CommonLVBTrailingByte Attribute = 0x0200

	// CommonLVBGridHorizontal indcates a horizontal grid.
	CommonLVBGridHorizontal Attribute = 0x0400

	// CommonLVBGridLVertical indicates a left vertical grid.
	CommonLVBGridLVertical Attribute = 0x0800

	// CommonLVBGridRVertical indicates a right vertical grid
	CommonLVBGridRVertical Attribute = 0x1000

	// CommonLVBReverseVideo reverses foreground and background attributes.
	CommonLVBReverseVideo Attribute = 0x4000

	// CommonLVGUnderscore adds an underscore below the character.
	CommonLVBUnderscore Attribute = 0x8000
)
