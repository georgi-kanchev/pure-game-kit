package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	d "pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func GUIs() {
	var cam = graphics.NewCamera(1)
	var atlas, icons = assets.LoadDefaultAtlasIcons()
	var _, ids, box = assets.LoadDefaultAtlasUI()
	var hud = gui.NewFromXMLs(cam, gui.NewElementsXML(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("label", f.Color, "0 0 0 0", f.Width, "150", f.Height, "50", f.GapX, "10", f.GapY, "10",
			f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.TextAlignmentX, "0", f.TextAlignmentY, "0.5", f.TextColor, "0 0 0 255",
			f.TextLineHeight, "40", f.TooltipId, "tooltip"),
		gui.Theme("button", f.Color, "220 220 220 255", f.Width, "150", f.Height, "50", f.GapX, "10", f.GapY, "10",
			f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[2], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.3", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "35", f.ButtonThemeIdHover, "button-hover", f.ButtonThemeIdPress, "button-press",
			f.TooltipId, "tooltip", f.SliderStep, "0.1", f.SliderHandleAssetId, box[14],
			f.SliderStepAssetId, ids[49], f.DraggableAssetColor, "0 0 255"),
		gui.Theme("button-hover", f.Color, "255 255 255 255", f.Width, "150", f.Height, "50",
			f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20",
			f.BoxEdgeBottom, "20", f.AssetId, box[5], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.3",
			f.TextColor, "127 127 127 255", f.TextLineHeight, "35", f.GapX, "10", f.GapY, "10",
			f.TooltipId, "tooltip"),
		gui.Theme("button-press", f.Color, "200 200 200 255", f.Width, "150", f.Height, "50",
			f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[4], f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.6", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "35", f.GapX, "10", f.GapY, "10", f.TooltipId, "tooltip"),
		gui.Theme("checkbox-on", f.Color, "220 220 220 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "35",
			f.ButtonThemeIdHover, "checkbox-on-hover", f.ButtonThemeIdPress, "checkbox-on-press",
			f.TooltipId, "tooltip", f.TooltipText, "Currently on!", f.Text, "X", f.TextAlignmentX, "0.6",
			f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-hover", f.Color, "255 255 255 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "35", f.Text, "X",
			f.TextAlignmentX, "0.6", f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-press", f.Color, "200 200 200 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[9], f.TextColor, "80 80 80 255", f.TextLineHeight, "35", f.Text, "X",
			f.TextAlignmentX, "0.6", f.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-off", f.Color, "220 220 220 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.ButtonThemeIdHover, "checkbox-off-hover", f.ButtonThemeIdPress, "checkbox-off-press",
			f.AssetId, box[9], f.TooltipId, "tooltip", f.TooltipText, "Currently off!"),
		gui.Theme("checkbox-off-hover", f.Color, "255 255 255 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[9]),
		gui.Theme("checkbox-off-press", f.Color, "200 200 200 255", f.Width, "50", f.Height, "50", f.GapX, "10",
			f.GapY, "10", f.BoxEdgeLeft, "20", f.BoxEdgeRight, "20", f.BoxEdgeTop, "20", f.BoxEdgeBottom, "20",
			f.AssetId, box[9]),
		gui.Container("panel", d.CameraLeftX+"+25", d.CameraCenterY+"-275", d.CameraWidth+"-50", "550",
			f.ThemeId, "button", f.GapX, "20", f.GapY, "10"),
		gui.Visual("background", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		gui.Visual("name-label", f.ThemeId, "label", f.Text, "Name", f.TooltipText, "Wow, tooltip for labels!"),
		gui.InputField("name", f.Width, "250", f.AssetId, box[9], f.Text, "Kenney", f.InputFieldMargin, "10",
			f.InputFieldPlaceholder, "Your name..."),
		gui.Visual("class-label", f.ThemeId, "label", f.Text, "Class", f.NewRow, ""),
		gui.InputField("class", f.Width, "250", f.AssetId, box[9], f.Text, "Cool", f.InputFieldMargin, "10"),
		gui.Visual("stepper-label", f.ThemeId, "label", f.Text, "Stepper", f.NewRow, ""),
		gui.Button("step-left", f.Width, "50", f.Text, "L",
			f.TooltipText, "Press this button to do absolutely nothing.", f.ButtonHotkey, "A"),
		gui.Visual("stepper", f.AssetId, box[9], f.Text, "10/10", f.TextAlignmentY, "0.5", f.GapX, "0",
			f.TextColor, "150 150 150 255"),
		gui.Button("step-right", f.Width, "50", f.Text, "R", f.GapX, "0",
			f.ButtonHotkey, "D"),
		gui.Visual("hotkey-info", f.ThemeId, "label", f.Text, "(press A or D)", f.Width, "225",
			f.TextColor, "127 127 127"),
		gui.Visual("checkbox-label", f.ThemeId, "label", f.Text, "Checkbox", f.NewRow, ""),
		gui.Checkbox("checkbox", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on"),
		gui.Menu("dropdown", f.NewRow, "", f.AssetId, box[9], f.Text, "V List selection", f.Width, "410",
			f.TextAlignmentX, "0.05", f.TextAlignmentY, "0.5", f.TextColor, "150 150 150 255",
			f.MenuContainerId, "menu"),
		gui.Visual("sliders-label", f.ThemeId, "label", f.Text, "Sliders", f.TextLineHeight, "30", f.NewRow, ""),
		gui.Slider("slider1", f.AssetId, box[10], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20"),
		gui.Visual("slider1-0", f.ThemeId, "label", f.Text, "0", f.TextLineHeight, "25", f.NewRow, "25",
			f.Width, "25", f.Height, "25", f.TextAlignmentX, "0.5"),
		gui.Visual("slider1-100", f.ThemeId, "label", f.Text, "100", f.TextLineHeight, "25", f.GapX, "355",
			f.Width, "40", f.Height, "25", f.TextAlignmentX, "0.5"),
		gui.Slider("slider2", f.AssetId, box[10], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20", f.SliderStep, "0"),
		gui.Visual("slider2-0", f.ThemeId, "label", f.Text, "0", f.TextLineHeight, "25", f.NewRow, "25",
			f.Width, "25", f.Height, "25", f.TextAlignmentX, "0.5"),
		gui.Visual("slider2-100", f.ThemeId, "label", f.Text, "100", f.TextLineHeight, "25", f.GapX, "355",
			f.Width, "40", f.Height, "25", f.TextAlignmentX, "0.5"),
		gui.Visual("divider", f.AssetId, box[12], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20", f.TooltipText, "Tooltips for dividers?! WHAT?"),
		gui.Button("x", f.Width, "85", f.Height, "70", f.Text, "X",
			f.NewRow, "35", f.TextColor, "255 255 255 255", f.Color, "200 0 0 255",
			f.TooltipText, "This is a pretty squarish X button."),
		gui.Button("v", f.Width, "315", f.Height, "70", f.Text, "V Accept ",
			f.TextColor, "255 255 255 255", f.Color, "0 200 0 255"),
		gui.Checkbox("outside1", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on", f.NewRow, "",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside2", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside3", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside4", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on",
			f.CheckboxGroup, "radio"),
		gui.Checkbox("outside5", f.ThemeId, "checkbox-off", f.CheckboxThemeId, "checkbox-on",
			f.CheckboxGroup, "radio"),
		gui.Visual("tab1", f.ThemeId, "label", f.Text, "First Tab", f.ToggleButtonId, "outside1", f.Hidden, "1",
			f.NewRow, "", f.Width, "250"),
		gui.Slider("tab1-slider", f.AssetId, box[10], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20", f.ToggleButtonId, "outside1", f.Hidden, "1"),
		gui.Visual("tab2", f.ThemeId, "label", f.Text, "Second Tab", f.ToggleButtonId, "outside2", f.Hidden, "1",
			f.NewRow, "", f.Width, "250"),
		gui.InputField("tab2-inputbox", f.Width, "250", f.AssetId, box[9], f.Text, "Cool Tab!",
			f.InputFieldPlaceholder, "Your input...", f.Hidden, "1", f.NewRow, "", f.ToggleButtonId, "outside2"),
		gui.Visual("tab3", f.ThemeId, "label", f.Text, "Third Tab", f.ToggleButtonId, "outside3", f.Hidden, "1",
			f.NewRow, "", f.Width, "250"),
		gui.Visual("tab4", f.ThemeId, "label", f.Text, "Fourth Tab", f.ToggleButtonId, "outside4", f.Hidden, "1",
			f.NewRow, "", f.Width, "250"),
		gui.Visual("tab5", f.ThemeId, "label", f.Text, "Fifth Tab", f.ToggleButtonId, "outside5", f.Hidden, "1",
			f.NewRow, "", f.Width, "250"),
		gui.Visual("divider-2", f.AssetId, "", f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20"),
		gui.Button("btn1", f.Text, "Parent1", f.NewRow, ""),
		gui.Button("btn2", f.Text, "Child1", f.ToggleButtonId, "btn1", f.OffsetX, "25", f.NewRow, ""),
		gui.Button("btn3", f.Text, "Child2", f.ToggleButtonId, "btn1", f.OffsetX, "25", f.NewRow, ""),
		gui.Button("btn4", f.Text, "Child3", f.ToggleButtonId, "btn1", f.OffsetX, "25", f.NewRow, ""),
		gui.Button("btn5", f.Text, "Parent2", f.NewRow, ""),
		gui.Button("btn6", f.Text, "Child", f.ToggleButtonId, "btn5", f.OffsetX, "25", f.NewRow, ""),
		gui.Button("btn7", f.Text, "Nested", f.ToggleButtonId, "btn5", f.OffsetX, "25", f.NewRow, ""),
		gui.Button("btn8", f.Text, "Deep", f.ToggleButtonId, "btn7", f.OffsetX, "50", f.NewRow, ""),
		gui.Visual("divider-3", f.AssetId, box[12], f.BoxEdgeTop, "0", f.BoxEdgeBottom, "0", f.NewRow, "",
			f.Width, "410", f.Height, "20"),
		gui.Draggable("slot-1", f.NewRow, "25", f.Width, "100", f.Height, "100", f.AssetId, box[0],
			f.DraggableAssetId, icons[50]),
		gui.Draggable("slot-2", f.Width, "100", f.Height, "100", f.AssetId, box[0], f.DraggableAssetId, icons[100]),
		gui.Draggable("slot-3", f.Width, "100", f.Height, "100", f.AssetId, box[0], f.TextColor, "test"),
		gui.Draggable("slot-4", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-5", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-6", f.NewRow, "", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-7", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-8", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-9", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Draggable("slot-10", f.Width, "100", f.Height, "100", f.AssetId, box[0]),
		gui.Container("menu", "", "", "410", "150", f.ThemeId, "button", f.Hidden, "1", f.GapX, "5", f.GapY, "5"),
		gui.Visual("menu-bgr", f.FillContainer, "", f.AssetId, box[0], f.Color, "200 200 200 255"),
		gui.Button("menu-1", f.GapX, "0", f.Width, "500", f.Text, "Monday"),
		gui.Button("menu-2", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Tuesday"),
		gui.Button("menu-3", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Wednesday"),
		gui.Button("menu-4", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Thursday",
			f.TooltipText, "It's thursday, wohooo!"),
		gui.Button("menu-5", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Friday"),
		gui.Visual("weekend-label", f.ThemeId, "label", f.Text, "Weekend", f.NewRow, ""),
		gui.Button("menu-6", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Saturday"),
		gui.Button("menu-7", f.GapX, "0", f.NewRow, "", f.Width, d.OwnerWidth, f.Text, "Sunday"),
		gui.Container("tooltips", "", "", "", "", f.ThemeId, "button", f.Hidden, "1"),
		gui.Tooltip("tooltip", f.AssetId, box[7], f.Width, "350", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5"),
	))
	assets.LoadDefaultSoundsUI()
	assets.SetTextureAtlasTile(atlas, "arrow-left", 14, 9, 1, 1, 0, true)
	assets.SetTextureAtlasTile(atlas, "arrow-down", 14, 9, 1, 1, 1, false)

	hud.Scale = 2.0

	for window.KeepOpen() {
		cam.DrawGrid(2, 100, 100, palette.DarkGray)

		var grab = hud.DragJustGrabbed()
		if grab != "" {
			fmt.Printf("grab: %v\n", grab)
		}

		var from, to = hud.DragJustDropped()
		if from != "" && to != "" {
			var fromId = hud.Field(from, f.DraggableAssetId)
			var toId = hud.Field(to, f.DraggableAssetId)
			hud.SetField(from, f.DraggableAssetId, toId)
			hud.SetField(to, f.DraggableAssetId, fromId)
		}

		hud.UpdateAndDraw()
		cam.DrawTextDebug(true, true, true, true)
	}
}
