package property

const (
	Class   = "class"   // [widget]
	Id      = "id"      // [widget] [container]
	X       = "x"       // [widget] [container]
	Y       = "y"       // [widget] [container]
	ThemeId = "themeId" // [widget] [container]

	Width   = "width"   // [widget] [theme] [container]
	Height  = "height"  // [widget] [theme] [container]
	AssetId = "assetId" // [widget] [theme] [container]
	GapX    = "gapX"    // [widget] [theme] [container]
	GapY    = "gapY"    // [widget] [theme] [container]

	FillContainer = "fillContainer" // [widget]
	NewRow        = "newRow"        // [widget] example: "200", "" = GapY
	OffsetX       = "offsetX"       // [widget]
	OffsetY       = "offsetY"       // [widget]

	Color                 = "color"                 // [widget] [theme] separated with space, example: "255 0 0 255"
	BoxEdgeLeft           = "boxEdgeLeft"           // [widget] [theme]
	BoxEdgeRight          = "boxEdgeRight"          // [widget] [theme]
	BoxEdgeTop            = "boxEdgeTop"            // [widget] [theme]
	BoxEdgeBottom         = "boxEdgeBottom"         // [widget] [theme]
	Text                  = "text"                  // [widget] [theme]
	TextFontId            = "textFontId"            // [widget] [theme]
	TextColor             = "textColor"             // [widget] [theme] separated with space, example: "255 0 0 255"
	TextLineHeight        = "textLineHeight"        // [widget] [theme] default: "60"
	TextLineGap           = "textLineGap"           // [widget] [theme]
	TextAlignmentX        = "textAlignmentX"        // [widget] [theme]
	TextAlignmentY        = "textAlignmentY"        // [widget] [theme]
	TextWordWrap          = "textWordWrap"          // [widget] [theme] example: "on" & "" = true, all else = false
	TextThickness         = "textThickness"         // [widget] [theme] default: "0.5"
	TextSmoothness        = "textSmoothness"        // [widget] [theme] default: "0.02"
	TextThicknessOutline  = "textThicknessOutline"  // [widget] [theme] default: "0.92"
	TextSmoothnessOutline = "textSmoothnessOutline" // [widget] [theme] default: "0.08"
	TextColorOutline      = "textColorOutline"      // [widget] [theme] separated with space, example: "255 0 0 255"

	ButtonAssetIdHover   = "buttonAssetIdHover"   // [widget] [theme]
	ButtonAssetIdPress   = "buttonAssetIdPress"   // [widget] [theme]
	ButtonAssetIdDisable = "buttonAssetIdDisable" // [widget] [theme]
)
