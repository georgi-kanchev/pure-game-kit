package property

const (
	Class   = "class"   // [widget]
	Id      = "id"      // [widget] [container]
	X       = "x"       // [widget] [container]
	Y       = "y"       // [widget] [container]
	ThemeId = "themeId" // [widget] [container]
	Hidden  = "hidden"  // [widget] [container]

	Width   = "width"   // [widget] [theme] [container]
	Height  = "height"  // [widget] [theme] [container]
	AssetId = "assetId" // [widget] [theme] [container]
	GapX    = "gapX"    // [widget] [theme] [container]
	GapY    = "gapY"    // [widget] [theme] [container]

	Disabled      = "disabled"      // [widget]
	FillContainer = "fillContainer" // [widget]
	NewRow        = "newRow"        // [widget] example: "200", "" = GapY
	OffsetX       = "offsetX"       // [widget]
	OffsetY       = "offsetY"       // [widget]

	Color         = "color"         // [widget] [theme] separated with space, example: "255 0 0 255"
	BoxEdgeLeft   = "boxEdgeLeft"   // [widget] [theme]
	BoxEdgeRight  = "boxEdgeRight"  // [widget] [theme]
	BoxEdgeTop    = "boxEdgeTop"    // [widget] [theme]
	BoxEdgeBottom = "boxEdgeBottom" // [widget] [theme]

	Text                  = "text"                  // [widget] [theme]
	TextFontId            = "textFontId"            // [widget] [theme]
	TextColor             = "textColor"             // [widget] [theme] separated with space, example: "255 0 0 255"
	TextLineHeight        = "textLineHeight"        // [widget] [theme] default: "60"
	TextLineGap           = "textLineGap"           // [widget] [theme]
	TextSymbolGap         = "textSymbolGap"         // [widget] [theme] default: "0.2"
	TextAlignmentX        = "textAlignmentX"        // [widget] [theme]
	TextAlignmentY        = "textAlignmentY"        // [widget] [theme]
	TextWordWrap          = "textWordWrap"          // [widget] [theme] example: "on" & "" = true, all else = false
	TextThickness         = "textThickness"         // [widget] [theme] default: "0.5"
	TextSmoothness        = "textSmoothness"        // [widget] [theme] default: "0.02"
	TextThicknessOutline  = "textThicknessOutline"  // [widget] [theme] default: "0.92"
	TextSmoothnessOutline = "textSmoothnessOutline" // [widget] [theme] default: "0.08"
	TextColorOutline      = "textColorOutline"      // [widget] [theme] separated with space, example: "255 0 0 255"

	TextEmbeddedAssetsTag      = "textEmbeddedAssetsTag"      // [widget] [theme] default: "^"
	TextEmbeddedAssetId1       = "textEmbeddedAssetId1"       // [widget] [theme]
	TextEmbeddedAssetId2       = "textEmbeddedAssetId2"       // [widget] [theme]
	TextEmbeddedAssetId3       = "textEmbeddedAssetId3"       // [widget] [theme]
	TextEmbeddedAssetId4       = "textEmbeddedAssetId4"       // [widget] [theme]
	TextEmbeddedAssetId5       = "textEmbeddedAssetId5"       // [widget] [theme]
	TextEmbeddedColorsTag      = "textEmbeddedColorsTag"      // [widget] [theme] default: "`"
	TextEmbeddedColor1         = "textEmbeddedColor1"         // [widget] [theme]
	TextEmbeddedColor2         = "textEmbeddedColor2"         // [widget] [theme]
	TextEmbeddedColor3         = "textEmbeddedColor3"         // [widget] [theme]
	TextEmbeddedColor4         = "textEmbeddedColor4"         // [widget] [theme]
	TextEmbeddedColor5         = "textEmbeddedColor5"         // [widget] [theme]
	TextEmbeddedThicknessesTag = "textEmbeddedThicknessesTag" // [widget] [theme] default: "*"
	TextEmbeddedThickness1     = "textEmbeddedThickness1"     // [widget] [theme] default: "0.5"
	TextEmbeddedThickness2     = "textEmbeddedThickness2"     // [widget] [theme] default: "0.5"
	TextEmbeddedThickness3     = "textEmbeddedThickness3"     // [widget] [theme] default: "0.5"
	TextEmbeddedThickness4     = "textEmbeddedThickness4"     // [widget] [theme] default: "0.5"
	TextEmbeddedThickness5     = "textEmbeddedThickness5"     // [widget] [theme] default: "0.5"

	ButtonHoverThemeId   = "buttonHoverThemeId"   // [widget] [theme]
	ButtonPressThemeId   = "buttonPressThemeId"   // [widget] [theme]
	ButtonDisableThemeId = "buttonDisableThemeId" // [widget] [theme]
)
