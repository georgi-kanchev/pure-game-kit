package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	d "pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func GUIs() {
	var cam = graphics.NewCamera(1)
	var atlas, icons = assets.LoadDefaultAtlasIcons()
	var _, ids, box = assets.LoadDefaultAtlasUI()
	var hud = gui.NewFromXML(gui.NewElementsXML(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("label", f.Color, "0 0 0 0", f.Width, "300", f.Height, "100", f.GapX, "20", f.GapY, "20",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.TextAlignmentX, "0", f.TextAlignmentY, "0.5", f.TextColor, "0 0 0 255",
			f.TextLineHeight, "80", f.TooltipId, "tooltip"),
		// ======================================================
		gui.Theme("button", f.Color, "220 220 220 255", f.Width, "300", f.Height, "100", f.GapX, "20", f.GapY, "20",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[2], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.3", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "70", f.ButtonThemeIdHover, "button-hover", f.ButtonThemeIdPress, "button-press",
			f.TooltipId, "tooltip", f.SliderStep, "0.1", f.SliderHandleAssetId, box[14],
			f.SliderStepAssetId, ids[49], f.DraggableSpriteColor, "0 0 255", f.DraggableSpriteScale, "0.8"),
		gui.Theme("button-hover", f.Color, "255 255 255 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40",
			f.BoxEdgeBottom, "40", f.AssetId, box[5], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.3",
			f.TextColor, "127 127 127 255", f.TextLineHeight, "70", f.GapX, "20", f.GapY, "20",
			f.TooltipId, "tooltip"),
		gui.Theme("button-press", f.Color, "200 200 200 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[4], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.6", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "70", f.GapX, "20", f.GapY, "20", f.TooltipId, "tooltip"),
		// ======================================================
		gui.Theme("checkbox-on", f.Color, "220 220 220 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "70",
			f.ButtonThemeIdHover, "checkbox-on-hover", f.ButtonThemeIdPress, "checkbox-on-press",
			f.TooltipId, "tooltip", f.TooltipText, "Currently on!",
			f.TextEmbeddedAssetId1, icons[249], f.Text, "^^ ", f.TextAlignmentX, "0.6", f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-hover", f.Color, "255 255 255 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "70",
			f.TextEmbeddedAssetId1, icons[249], f.Text, "^^ ", f.TextAlignmentX, "0.6", f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-press", f.Color, "200 200 200 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "70",
			f.TextEmbeddedAssetId1, icons[249], f.Text, "^^ ", f.TextAlignmentX, "0.6", f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-off", f.Color, "220 220 220 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.ButtonThemeIdHover, "checkbox-off-hover", f.ButtonThemeIdPress, "checkbox-off-press",
			f.AssetId, box[9], f.TooltipId, "tooltip", f.TooltipText, "Currently off!"),
		gui.Theme("checkbox-off-hover", f.Color, "255 255 255 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[9]),
		gui.Theme("checkbox-off-press", f.Color, "200 200 200 255", f.Width, "100", f.Height, "100", f.GapX, "20",
			f.GapY, "20", f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[9]),
		// ======================================================
		gui.Container("panel", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "650",
			f.ThemeId, "button", f.GapX, "40", f.GapY, "20"),
		gui.Visual("background", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		// ======================================================
		gui.Visual("name-label", f.ThemeId, "label", f.Text, "Name", f.TooltipText, "Wow, tooltip for labels!"),
		gui.InputField("name", f.Width, "500", f.AssetId, box[9], f.Text, "Kenney",
			f.InputFieldPlaceholder, "Your name..."),
		gui.Visual("class-label", f.ThemeId, "label", f.Text, "Class", f.NewRow, ""),
		gui.InputField("class", f.Width, "500", f.AssetId, box[9], f.Text, "Cool"),
		gui.Visual("stepper-label", f.ThemeId, "label", f.Text, "Stepper", f.NewRow, ""),
		gui.Button("step-left", f.Width, "100", f.TextEmbeddedAssetId1, "arrow-left", f.Text, "^^ ",
			f.TooltipText, "Press this button to do absolutely nothing.", f.ButtonHotkey, "A"),
		gui.Visual("stepper", f.AssetId, box[9], f.Text, "10/10", f.TextAlignmentY, "0.5", f.GapX, "0",
			f.TextColor, "150 150 150 255"),
		gui.Button("step-right", f.Width, "100", f.TextEmbeddedAssetId1, icons[212], f.Text, "^^", f.GapX, "0"),
		gui.Visual("checkbox-label", f.ThemeId, "label", f.Text, "Checkbox", f.NewRow, ""),
		gui.Checkbox("checkbox", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on"),
		gui.Menu("dropdown", f.NewRow, "", f.AssetId, box[9], f.Text, "^^ List selection", f.Width, "820",
			f.TextAlignmentX, "0.05", f.TextAlignmentY, "0.5", f.TextColor, "150 150 150 255",
			f.TextEmbeddedAssetId1, "arrow-down", f.MenuContainerId, "menu"),
		gui.Visual("sliders-label", f.ThemeId, "label", f.Text, "Sliders", f.TextLineHeight, "60", f.NewRow, ""),
		gui.Slider("slider1", f.AssetId, box[10], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "820", f.Height, "40"),
		gui.Visual("slider1-0", f.ThemeId, "label", f.Text, "0", f.TextLineHeight, "50", f.NewRow, "50",
			f.Width, "50", f.Height, "50", f.TextAlignmentX, "0.5"),
		gui.Visual("slider1-100", f.ThemeId, "label", f.Text, "100", f.TextLineHeight, "50", f.GapX, "710",
			f.Width, "80", f.Height, "50", f.TextAlignmentX, "0.5"),
		gui.Slider("slider2", f.AssetId, box[10], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "820", f.Height, "40", f.SliderStep, "0"),
		gui.Visual("slider2-0", f.ThemeId, "label", f.Text, "0", f.TextLineHeight, "50", f.NewRow, "50",
			f.Width, "50", f.Height, "50", f.TextAlignmentX, "0.5"),
		gui.Visual("slider2-100", f.ThemeId, "label", f.Text, "100", f.TextLineHeight, "50", f.GapX, "710",
			f.Width, "80", f.Height, "50", f.TextAlignmentX, "0.5"),
		gui.Visual("divider", f.AssetId, box[12], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "820", f.Height, "40", f.TooltipText, "Tooltips for dividers?! WHAT?"),
		gui.Button("x", f.Width, "170", f.Height, "140", f.TextEmbeddedAssetId1, icons[250], f.Text, "^^ ",
			f.NewRow, "70", f.TextColor, "255 255 255 255", f.Color, "200 0 0 255",
			f.TooltipText, "This is a pretty squarish X button."),
		gui.Button("v", f.Width, "630", f.Height, "140", f.TextEmbeddedAssetId1, icons[249], f.Text, "^^ Accept ",
			f.TextColor, "255 255 255 255", f.Color, "0 200 0 255"),
		gui.Checkbox("outside1", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside2", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside3", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside4", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside5", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Visual("divider-2", f.AssetId, box[12], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "820", f.Height, "40"),
		gui.Draggable("slot-1", f.NewRow, "50", f.Width, "200", f.Height, "200", f.AssetId, box[0],
			f.DraggableSpriteId, icons[50]),
		gui.Draggable("slot-2", f.Width, "200", f.Height, "200", f.AssetId, box[0], f.DraggableSpriteId, icons[100]),
		gui.Draggable("slot-3", f.Width, "200", f.Height, "200", f.AssetId, box[0], f.TextColor, "test"),
		gui.Draggable("slot-4", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-5", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-6", f.NewRow, "", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-7", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-8", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-9", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		gui.Draggable("slot-10", f.Width, "200", f.Height, "200", f.AssetId, box[0]),
		// ======================================================
		gui.Container("menu", "", "", "820", "300", f.ThemeId, "button", f.Hidden, "+", f.GapX, "10", f.GapY, "10"),
		gui.Visual("menu-bgr", f.FillContainer, "", f.AssetId, box[0], f.Color, "200 200 200 255"),
		gui.Button("menu-1", f.GapX, "0", f.Width, "1000", f.Text, "Monday"),
		gui.Button("menu-2", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Tuesday"),
		gui.Button("menu-3", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Wednesday"),
		gui.Button("menu-4", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Thursday",
			f.TooltipText, "It's thursday, wohooo!"),
		gui.Button("menu-5", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Friday"),
		gui.Visual("weekend-label", f.ThemeId, "label", f.Text, "Weekend", f.NewRow, ""),
		gui.Button("menu-6", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Saturday"),
		gui.Button("menu-7", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Sunday"),
		// ======================================================
		gui.Container("tooltips", "", "", "", "", f.ThemeId, "button", f.Hidden, "+"),
		gui.Tooltip("tooltip", f.AssetId, box[7], f.Width, "700", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5"),
	))
	assets.LoadDefaultFont()
	assets.LoadDefaultSoundsUI()
	assets.SetTextureAtlasTile(atlas, "arrow-left", 14, 9, 1, 1, 0, true)
	assets.SetTextureAtlasTile(atlas, "arrow-down", 14, 9, 1, 1, 1, false)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.DarkGray)

		var grab = hud.DragOnGrab()
		if grab != "" {
			fmt.Printf("grab: %v\n", grab)
		}

		var from, to = hud.DragOnDrop()
		if from != "" && to != "" {
			var fromId = hud.Field(from, f.DraggableSpriteId)
			var toId = hud.Field(to, f.DraggableSpriteId)
			hud.SetField(from, f.DraggableSpriteId, toId)
			hud.SetField(to, f.DraggableSpriteId, fromId)
		}

		hud.UpdateAndDraw(cam)
	}
}
