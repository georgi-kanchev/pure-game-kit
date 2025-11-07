package tilemap

import (
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	p "pure-game-kit/tiled/property"
	"pure-game-kit/tiled/tileset"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

func Property(mapId, property string) string {
	var data, has = internal.TiledMaps[mapId]
	if !has {
		return ""
	}

	switch property {
	case p.MapName:
		return data.Name
	case p.MapClass:
		return data.Class
	case p.MapTileWidth:
		return text.New(data.TileWidth)
	case p.MapTileHeight:
		return text.New(data.TileHeight)
	case p.MapColumns:
		return text.New(data.Width)
	case p.MapRows:
		return text.New(data.Height)
	case p.MapParallaxX:
		return text.New(data.ParallaxOriginX)
	case p.MapParallaxY:
		return text.New(data.ParallaxOriginY)
	case p.MapInfinite:
		return text.New(data.Infinite)
	case p.MapBackgroundColor:
		return text.New(col(data.BackgroundColor))
	}

	for _, v := range data.Properties {
		if v.Name == property {
			return v.Value
		}
	}
	return ""
}

func LayerProperty(mapId, layerNameOrId, property string) string {
	var mapData, _ = internal.TiledMaps[mapId]
	var _, objs, img, layer = findLayer(mapData, layerNameOrId)
	if mapData == nil || layer == nil {
		return ""
	}

	switch property {
	case p.LayerName:
		return layer.Name
	case p.LayerClass:
		return layer.Class
	case p.LayerVisible:
		return condition.If(layer.Visible == "", "true", "false")
	case p.LayerLocked:
		return text.New(layer.Locked)
	case p.LayerOpacity:
		return text.New(layer.Opacity)
	case p.LayerTint:
		return text.New(col(layer.Tint))
	case p.LayerOffsetX:
		return text.New(layer.OffsetX)
	case p.LayerOffsetY:
		return text.New(layer.OffsetY)
	case p.LayerParallaxX:
		return text.New(layer.ParallaxX)
	case p.LayerParallaxY:
		return text.New(layer.ParallaxY)
	//=================================================================
	case p.LayerColor:
		return text.New(col(objs.Color))
	case p.LayerDrawOrder:
		return text.New(objs.DrawOrder)
	//=================================================================
	case p.LayerImage:
		return text.New(img.Image.Source)
	case p.LayerTransparentColor:
		return text.New(img.Image.TransparentColor)
	case p.LayerRepeatX:
		return text.New(img.RepeatX)
	case p.LayerRepeatY:
		return text.New(img.RepeatY)
	case p.LayerWidth:
		return text.New(img.Image.Width)
	case p.LayerHeight:
		return text.New(img.Image.Height)
	}

	for _, prop := range layer.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}

func LayerTileId(mapId, layerNameOrId string, cellX, cellY int) uint32 {
	var mapData, _ = internal.TiledMaps[mapId]
	var wantedLayer, _, _, _ = findLayer(mapData, layerNameOrId)
	if mapData == nil || wantedLayer == nil {
		return 0
	}

	cellX -= int(mapData.WorldX) / mapData.TileWidth
	cellY -= int(mapData.WorldY) / mapData.TileHeight
	var cellIndex = number.Indexes2DToIndex1D(cellY, cellX, mapData.Width, mapData.Height)
	var tilesets = usedTilesets(mapData)
	var tilesIds = getTileIds(mapData, tilesets, wantedLayer)
	var tileId = tilesIds[cellIndex]
	var unoriented = flag.TurnOff(tileId, internal.Flips)
	var curTileset = currentTileset(tilesets, unoriented)
	return unoriented - curTileset.FirstTileId
}

func LayerObjectProperty(mapId, layerNameOrId, objectNameClassOrId, property string) string {
	var obj = getObj(mapId, layerNameOrId, objectNameClassOrId)
	if obj == nil {
		return ""
	}

	switch property {
	case p.ObjectName:
		return obj.Name
	case p.ObjectClass:
		return obj.Class
	case p.ObjectTemplate:
		return obj.Template
	case p.ObjectVisible:
		return condition.If(obj.Visible == "", "true", "false")
	case p.ObjectLocked:
		return text.New(obj.Locked)
	case p.ObjectX:
		return text.New(obj.X)
	case p.ObjectY:
		return text.New(obj.Y)
	case p.ObjectWidth:
		return text.New(obj.Width)
	case p.ObjectHeight:
		return text.New(obj.Height)
	case p.ObjectRotation:
		return text.New(obj.Rotation)
	case p.ObjectFlipX:
		return text.New(flag.IsOn(obj.Gid, internal.FlipX))
	case p.ObjectFlipY:
		return text.New(flag.IsOn(obj.Gid, internal.FlipY))
	case p.ObjectTileId:
		var id = flag.TurnOff(obj.Gid, internal.FlipX)
		id = flag.TurnOff(id, internal.FlipY)
		var mapData, _ = internal.TiledMaps[mapId]
		var current = currentTileset(usedTilesets(mapData), id)
		return text.New(id - current.FirstTileId)
	}

	for _, prop := range obj.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}
func LayerSprites(mapId, layerNameOrId, objectNameClassOrId string) []*graphics.Sprite {
	var mapData, _ = internal.TiledMaps[mapId]
	var tiles, objs, img, _ = findLayer(mapData, layerNameOrId)
	if mapData == nil {
		return []*graphics.Sprite{}
	}

	if img != nil {
		var assetId = path.RemoveExtension(path.New(mapData.Directory, img.Image.Source))
		var sprite = graphics.NewSprite(assetId, mapData.WorldX+img.OffsetX, mapData.WorldY+img.OffsetY)
		sprite.PivotX, sprite.PivotY = 0, 0
		return []*graphics.Sprite{sprite}
	}

	if objs != nil {
		var result = []*graphics.Sprite{}
		var usedTilesets = usedTilesets(mapData)
		for _, obj := range objs.Objects {
			if objectNameClassOrId != "" && obj.Name != objectNameClassOrId && obj.Class != objectNameClassOrId &&
				text.New(obj.Id) != objectNameClassOrId {
				continue
			}
			if obj.Gid == 0 {
				continue
			}

			var id = flag.TurnOff(obj.Gid, internal.FlipX)
			id = flag.TurnOff(id, internal.FlipY)
			var curTileset = currentTileset(usedTilesets, id)
			if curTileset == nil {
				continue
			}

			var assetId = text.New(curTileset.AtlasId, "/", id-curTileset.FirstTileId)
			var sprite = graphics.NewSprite(assetId, mapData.WorldX+obj.X, mapData.WorldY+obj.Y)
			sprite.X += float32(mapData.TileWidth)/2 + objs.OffsetX
			sprite.Y = sprite.Y - float32(mapData.TileHeight)/2 + objs.OffsetY
			sprite.Width, sprite.Height = float32(mapData.TileWidth), float32(mapData.TileHeight)
			sprite.ScaleX = condition.If(flag.IsOn(obj.Gid, internal.FlipX), float32(-1), 1)
			sprite.ScaleY = condition.If(flag.IsOn(obj.Gid, internal.FlipY), float32(-1), 1)
			sprite.Angle = obj.Rotation

			result = append(result, sprite)
		}
		return result
	}

	if tiles == nil {
		return []*graphics.Sprite{}
	}

	var result = make([]*graphics.Sprite, 0, mapData.Width*mapData.Height)
	var usedTilesets = usedTilesets(mapData)
	var tileIds = getTileIds(mapData, usedTilesets, tiles)

	for index, tile := range tileIds {
		var unoriented = flag.TurnOff(tile, internal.Flips)
		var curTileset = currentTileset(usedTilesets, unoriented)
		if curTileset == nil {
			continue
		}

		var width, height = float32(curTileset.TileWidth), float32(curTileset.TileHeight)
		width, height = tileRenderSize(width, height, mapData, curTileset)
		var ang, w, h = getTileOrientation(tile, width, height)
		var id = unoriented - curTileset.FirstTileId
		var tileId = text.New(curTileset.AtlasId, "/", id)
		var px, py float32 = 0.5, 0.5
		var j, i = number.Index1DToIndexes2D(index, mapData.Width, mapData.Height)
		var offX, offY = w / 2, h / 2

		if curTileset.Image.Source == "" {
			var tileObj = curTileset.MappedTiles[id]
			tileId = tileObj.TextureId
			w = float32(tileObj.Image.Width * condition.If(w < 0, -1, 1))
			h = float32(tileObj.Image.Height * condition.If(h < 0, -1, 1))
			w, h = tileRenderSize(w, h, mapData, curTileset)
			px, py = 0, 1
			offX, offY = 0, 0
			i++

			if curTileset.FillMode == "preserve-aspect-fit" {
				offX = float32(mapData.TileWidth)/2 - w/2
			}

			switch ang {
			case 90:
				offY -= w
			case 180:
				offX += w
				offY -= h
			case 270:
				offX += h
			}

			if ang == 90 && w < 0 {
				offX = w
				offY = 0
			} else if ang == 270 && w < 0 {
				offX = w + h
				offY = w
			}
		}

		if w < 0 {
			offX -= w
		}
		if h < 0 {
			offY -= h
		}

		var x = float32(j)*float32(mapData.TileWidth) + mapData.WorldX + tiles.OffsetX
		var y = float32(i)*float32(mapData.TileHeight) + mapData.WorldY + tiles.OffsetY
		var sprite = graphics.NewSprite(tileId, 0, 0)

		sprite.Angle = ang
		sprite.Width, sprite.Height = w, h
		sprite.PivotX, sprite.PivotY = px, py

		x += float32(curTileset.Offset.X)
		y += float32(curTileset.Offset.Y)
		sprite.X, sprite.Y = x+offX, y+offY

		result = append(result, sprite)
	}

	return result
}
func LayerTexts(mapId, layerNameOrId, objectNameClassOrId string) []*graphics.TextBox {
	var mapData, _ = internal.TiledMaps[mapId]
	var _, objs, _, _ = findLayer(mapData, layerNameOrId)
	var result = []*graphics.TextBox{}
	if mapData == nil {
		return result
	}

	if objs != nil {
		for _, obj := range objs.Objects {
			if objectNameClassOrId != "" && obj.Name != objectNameClassOrId && obj.Class != objectNameClassOrId &&
				text.New(obj.Id) != objectNameClassOrId {
				continue
			}
			if obj.Text.Value == "" {
				continue
			}

			var textbox = graphics.NewTextBox("", obj.X, obj.Y, obj.Text.Value)
			textbox.X += mapData.WorldX + objs.OffsetX
			textbox.Y += mapData.WorldY + objs.OffsetY
			textbox.Width, textbox.Height = float32(obj.Width), float32(obj.Height)
			textbox.Color = color.White
			textbox.LineHeight = float32(obj.Text.FontSize)
			textbox.PivotX, textbox.PivotY = 0, 0

			if obj.Text.Color != "" {
				textbox.Color = col(obj.Text.Color)
			}
			if obj.Text.FontSize == 0 {
				textbox.LineHeight = float32(mapData.TileHeight)
			}
			if obj.Text.Bold {
				textbox.Thickness = 0.75
			}

			switch obj.Text.AlignX {
			case "center":
				textbox.AlignmentX = 0.5
			case "right":
				textbox.AlignmentX = 1
			}

			switch obj.Text.AlignY {
			case "center":
				textbox.AlignmentY = 0.5
			case "bottom":
				textbox.AlignmentY = 1
			}

			result = append(result, &textbox)
		}
	}

	return result
}

// empty objectNameOrClass includes all objects within all tiles
//
// accepts only tile layers, ignores the rest
//
// layer offset snaps the shapes to the nearest cell (the grid does not move) -
// it is best to snap the layer offset to the tile size (in Tiled)
// so that the tile sprites and the tile shapes match visually
func LayerShapeGrid(mapId, tileLayerNameOrId, objectNameOrClass string) *geometry.ShapeGrid {
	var mapData, _ = internal.TiledMaps[mapId]
	if mapData == nil {
		return nil
	}

	var result = geometry.NewShapeGrid(mapData.TileWidth, mapData.TileHeight)
	var success = forEachTile(mapId, tileLayerNameOrId,
		func(x, y int, id uint32, layer *internal.LayerTiles, curTileset *internal.Tileset) {
			x += int(number.Round(layer.OffsetX/float32(curTileset.TileWidth), 0))
			y += int(number.Round(layer.OffsetY/float32(curTileset.TileHeight), 0))
			result.SetAtCell(x, y, tileset.TileObjectShapes(curTileset.AtlasId, id, objectNameOrClass)...)
		})
	return condition.If(success, result, nil)
}

// empty objectNameOrClass includes all objects within all tiles/layer
func LayerShapes(mapId, layerNameOrId, objectNameClassOrId string) []*geometry.Shape {
	var mapData, _ = internal.TiledMaps[mapId]
	var result = []*geometry.Shape{}
	if mapData == nil {
		return result
	}

	var tiles, objs, img, _ = findLayer(mapData, layerNameOrId)

	if img != nil {
		var shape = geometry.NewShapeRectangle(float32(img.Image.Width), float32(img.Image.Height), 0, 0)
		shape.X, shape.Y = mapData.WorldX+img.OffsetX, mapData.WorldY+img.OffsetY
		return []*geometry.Shape{shape}
	}
	if tiles != nil {
		return LayerShapeGrid(mapId, layerNameOrId, objectNameClassOrId).All()
	}
	if objs == nil {
		return result
	}

	for _, obj := range objs.Objects {
		if objectNameClassOrId != "" && obj.Name != objectNameClassOrId && obj.Class != objectNameClassOrId &&
			text.New(obj.Id) != objectNameClassOrId {
			continue
		}
		var ptsData = ""
		if obj.Polyline.Points != "" {
			ptsData = obj.Polyline.Points
		}
		if obj.Polygon.Points != "" {
			ptsData = obj.Polygon.Points
		}
		if ptsData == "" {
			var w, h = obj.Width, obj.Height
			ptsData = text.New(0, ",", 0, " ", w, ",", 0, " ", w, ",", h, " ", 0, ",", h)
		}
		var points = [][2]float32{}
		var pts = text.Split(ptsData, " ")
		for _, pt := range pts {
			var xy = text.Split(pt, ",")
			if len(xy) == 2 {
				var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
				x, y = point.RotateAroundPoint(x, y, 0, 0, obj.Rotation)
				points = append(points, [2]float32{x, y})
			}
		}

		var shape = geometry.NewShapeCorners(points...)
		shape.X = obj.X + mapData.WorldX + objs.OffsetX
		shape.Y = obj.Y + mapData.WorldY + objs.OffsetY

		if obj.Gid != 0 { // adjust tile object pivot
			shape.Y -= float32(mapData.TileHeight)
		}

		result = append(result, shape)
	}

	return result
}

// empty objectNameOrClass includes all objects within all tiles/layer
func LayerPoints(mapId, layerNameOrId, objectNameOrClass string) [][2]float32 {
	var result = [][2]float32{}
	var mapData, _ = internal.TiledMaps[mapId]
	if mapData == nil {
		return result
	}

	var success = forEachTile(mapId, layerNameOrId,
		func(x, y int, id uint32, layer *internal.LayerTiles, curTileset *internal.Tileset) {
			var pts = tileset.TileObjectPoints(curTileset.AtlasId, id, objectNameOrClass)
			for i := range pts {
				pts[i][0] += float32(x)*float32(curTileset.TileWidth) + layer.OffsetX
				pts[i][1] += float32(y)*float32(curTileset.TileHeight) + layer.OffsetY
			}
			result = append(result, pts...)
		})

	if success {
		return result
	}

	var _, objs, _, _ = findLayer(mapData, layerNameOrId)
	if objs == nil {
		return result
	}

	for _, obj := range objs.Objects {
		if obj.Width == 0 && obj.Height == 0 && obj.Polygon.Points == "" && obj.Polyline.Points == "" {
			result = append(result, [2]float32{mapData.WorldX + obj.X, mapData.WorldY + obj.Y})
		}
	}
	return result
}
