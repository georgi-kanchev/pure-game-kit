package tiled

import (
	"pure-game-kit/execution/condition"
	geo "pure-game-kit/geometry"
	gfx "pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/text"
)

func defaultValueText(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func parseProperty(prop *internal.Property, project *Project) any {
	switch prop.Type {
	case "bool":
		return text.LowerCase(prop.Value) == "true"
	case "int", "object":
		return text.ToNumber[int](prop.Value)
	case "float":
		return text.ToNumber[float32](prop.Value)
	case "color":
		return color.Hex(prop.Value)
	case "class":
		if project == nil {
			return prop.Value
		}

		var class, hasClass = project.Classes[prop.CustomType]
		if !hasClass {
			return prop.Value
		}

		var classMembers = class.(map[string]any)
		var result = make(map[string]any, len(classMembers))
		for n, v := range classMembers {
			result[n] = v

			for _, prop := range prop.Properties {
				if prop.Name == n {
					result[prop.Name] = parseProperty(&prop, project)
				}
			}
		}
		return result
	}
	return prop.Value
}
func currentTileset(Map *Map, tile uint32) (tileset *Tileset, firstId uint32) {
	var result *Tileset
	var bestFirstId uint32 = 0
	for i, tileset := range Map.Tilesets {
		var firstId = Map.TilesetsFirstTileIds[i]
		if firstId <= tile && firstId >= bestFirstId {
			bestFirstId = firstId
			result = tileset
		}
	}
	return result, bestFirstId
}
func tileOrientation(tileId uint32, w, h, th float32, image bool) (ang, newW, newH, offX, offY float32) {
	var flipH = flag.IsOn(tileId, internal.FlipX)
	var flipV = flag.IsOn(tileId, internal.FlipY)
	var flipDiag = flag.IsOn(tileId, internal.FlipDiag)

	ang = 0.0
	newW, newH = w, h
	offX, offY = 0, condition.If(image, th-h, 0)

	if flipH && !flipV && flipDiag { // rotation 90
		ang = 90
		offX = h
		offY = condition.If(image, th-w, 0)
	} else if flipH && flipV && !flipDiag { // rotation 180
		ang = 180
		offX = w
		offY = condition.If(image, th, h)
	} else if !flipH && flipV && flipDiag { // rotation 270
		ang = 270
		offY = condition.If(image, th, w)
	} else if flipH && !flipV && !flipDiag { // flip x only
		newW = -w
		offX = w
	} else if flipH && flipV && flipDiag { // flip x + rotation 90
		ang = 90
		newW = -w
		offX = h
		offY = condition.If(image, th, w)
	} else if !flipH && flipV && !flipDiag { // flip x + rotation 180
		newH = -h
		offY = condition.If(image, th, h)
	} else if !flipH && !flipV && flipDiag { // flip x + rotation 270
		ang = 270
		newW = -w
		offY = condition.If(image, th-w, 0)
	}
	return ang, newW, newH, offX, offY
}
func draw(c *gfx.Camera, spr []*gfx.Sprite, txt []*gfx.TextBox, sh []*geo.Shape, pt, ln [][2]float32, col uint) {
	c.DrawSprites(spr...)
	c.DrawTextBoxes(txt...)

	for _, shape := range sh {
		var pts = shape.CornerPoints()
		c.DrawShapes(color.FadeOut(col, 0.5), pts...)
		c.DrawLinesPath(0.5, col, pts...)
	}

	c.DrawLinesPath(0.5, col, ln...)
	c.DrawPoints(0.5, col, pt...)
}
