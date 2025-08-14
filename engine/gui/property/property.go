package property

const (
	Class   = "class"   // [widget]
	Id      = "id"      // [container] [widget]
	X       = "x"       // [container] [widget]
	Y       = "y"       // [container] [widget]
	ThemeId = "themeId" // [container] [widget]

	NewRow  = "newRow"  // [widget] example: "200", "" = auto (max height row)
	OffsetX = "offsetX" // [widget]
	OffsetY = "offsetY" // [widget]

	Width  = "width"  // [theme] [container] [widget]
	Height = "height" // [theme] [container] [widget]
	Color  = "color"  // [theme] [container] [widget] separated with space, example: "255 0 0 255"

	Text           = "text"           // [theme] [container] [widget]
	TextFontId     = "textFontId"     // [theme] [container] [widget]
	TextColor      = "textColor"      // [theme] [container] [widget] separated with space, example: "255 0 0 255"
	TextLineHeight = "textLineHeight" // [theme] [container] [widget]
	TextLineGap    = "textLineGap"    // [theme] [container] [widget]
)
