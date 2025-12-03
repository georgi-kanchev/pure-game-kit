package field

const (
	Class    = "class"    // [widget]
	Id       = "id"       // [widget] [container]
	X        = "x"        // [widget] [container]
	Y        = "y"        // [widget] [container]
	ThemeId  = "themeId"  // [widget] [container]
	Hidden   = "hidden"   // [widget] [container] "1" = true | "" = false
	Disabled = "disabled" // [widget] [container] "1" = true | "" = false

	Width    = "width"    // [widget] [theme] [container]
	Height   = "height"   // [widget] [theme] [container]
	AssetId  = "assetId"  // [widget] [theme] [container]
	GapX     = "gapX"     // [widget] [theme] [container]
	GapY     = "gapY"     // [widget] [theme] [container]
	TargetId = "targetId" // [container]

	Value         = "value"         // [widget]
	FillContainer = "fillContainer" // [widget]
	NewRow        = "newRow"        // [widget] example: "200", "" = GapY
	OffsetX       = "offsetX"       // [widget]
	OffsetY       = "offsetY"       // [widget]
	TooltipText   = "tooltipText"   // [widget]

	TooltipId     = "tooltipId"     // [widget] [theme]
	TooltipMargin = "tooltipMargin" // [widget] [theme]
	TooltipSound  = "tooltipSound"  // [widget] [theme] default: "~popup"
	Color         = "color"         // [widget] [theme] example: "255 0 0 255"
	FrameColor    = "frameColor"    // [widget] [theme] example: "255 0 0 255"
	FrameSize     = "frameSize"     // [widget] [theme] positive outward, negative inward
	BoxEdgeLeft   = "boxEdgeLeft"   // [widget] [theme]
	BoxEdgeRight  = "boxEdgeRight"  // [widget] [theme]
	BoxEdgeTop    = "boxEdgeTop"    // [widget] [theme]
	BoxEdgeBottom = "boxEdgeBottom" // [widget] [theme]

	Text                  = "text"                  // [widget] [theme]
	TextFontId            = "textFontId"            // [widget] [theme]
	TextColor             = "textColor"             // [widget] [theme] default: "127 127 127 255"
	TextLineHeight        = "textLineHeight"        // [widget] [theme] default: "30"
	TextLineGap           = "textLineGap"           // [widget] [theme]
	TextSymbolGap         = "textSymbolGap"         // [widget] [theme] default: "0.2"
	TextAlignmentX        = "textAlignmentX"        // [widget] [theme]
	TextAlignmentY        = "textAlignmentY"        // [widget] [theme]
	TextWordWrap          = "textWordWrap"          // [widget] [theme] example: "on" & "" = true, all else = false
	TextThickness         = "textThickness"         // [widget] [theme] default: "0.5"
	TextSmoothness        = "textSmoothness"        // [widget] [theme] default: "0.02"
	TextThicknessOutline  = "textThicknessOutline"  // [widget] [theme] default: "0.92"
	TextSmoothnessOutline = "textSmoothnessOutline" // [widget] [theme] default: "0.08"
	TextColorOutline      = "textColorOutline"      // [widget] [theme] example: "255 0 0 255"

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

	ButtonThemeIdHover   = "buttonThemeIdHover"   // [widget] [theme]
	ButtonThemeIdPress   = "buttonThemeIdPress"   // [widget] [theme]
	ButtonThemeIdDisable = "buttonThemeIdDisable" // [widget] [theme]
	ButtonHotkey         = "buttonHotkey"         // [widget] [theme]
	ButtonSoundPress     = "buttonSoundPress"     // [widget] [theme] default: "~press"
	ButtonSoundRelease   = "buttonSoundRelease"   // [widget] [theme] default: "~release"

	SliderStep          = "sliderStep"          // [widget] [theme] step <= 0 hides indicators, example: "-0.1", "0.2"
	SliderStepAssetId   = "sliderStepAssetId"   // [widget] [theme]
	SliderHandleAssetId = "sliderHandleAssetId" // [widget] [theme]
	SliderSound         = "sliderSound"         // [widget] [theme] default: "~slider"

	CheckboxThemeId  = "checkboxThemeId"  // [widget] [theme]
	CheckboxGroup    = "checkboxGroup"    // [widget] [theme]
	CheckboxSoundOn  = "checkboxSoundOn"  // [widget] [theme] default: "~on"
	CheckboxSoundOff = "checkboxSoundOff" // [widget] [theme] default: "~off"

	MenuContainerId = "menuContainerId" // [widget] [theme]
	MenuSound       = "menuSound"       // [widget] [theme] default: "~popup"

	InputFieldMargin      = "inputFieldMargin"      // [widget] [theme] default: "30"
	InputFieldPlaceholder = "inputFieldPlaceholder" // [widget] [theme] default: "Type..."
	InputFieldSoundType   = "inputFieldSoundType"   // [widget] [theme] default: "~write"
	InputFieldSoundErase  = "inputFieldSoundErase"  // [widget] [theme] default: "~erase"

	DraggableSpriteId    = "draggableSpriteId"    // [widget] [theme]
	DraggableSpriteColor = "draggableSpriteColor" // [widget] [theme] default: "255 255 255 255"
	DraggableSpriteScale = "draggableSpriteScale" // [widget] [theme] default: "1"
	DraggableSoundCancel = "draggableSoundCancel" // [widget] [theme] default: "~error"
)
