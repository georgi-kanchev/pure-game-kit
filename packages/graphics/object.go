package graphics

import (
	"pure-game-kit/packages/assets"
	geometry "pure-game-kit/packages/geometry2"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	txt "pure-game-kit/packages/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Object struct {
	geometry.Shape

	Mask    Area // In window space.
	Effects Effects

	// image ==========================================================

	ImageId   assets.ImageId
	ImageCrop Area // Zero value = original asset image (or asset crop)

	// text ===========================================================

	// Tags for embeded effects:
	//
	//	Underline:     ✅ // toggle
	//	Crossout:      ❎ // toggle
	//	Weight:        ⏬🔽🔁🔼⏫
	//	Size:          🔇🔈🔉🔊📢
	//	Color:         ⬜⬛🟥🟧🟨🟩🟦🟪🟫
	//	Outline Color: ⚪⚫🔴🟠🟡🟢🔵🟣🟤
	Text       string
	TextFontId assets.FontId

	// Caches the text visuals across frames. Call TextUpdateBatch when visual changes are needed.
	// Useful for a huge static text that changes rarely.
	TextBatch   bool
	textBatches []*internal.Batch

	// tilemap ========================================================

	TileAtlasId assets.TileAtlasId
	TileLayerId assets.TileLayerId
}

func NewShapePoint(x, y float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewPoint(x, y), Effects: eff}
}
func NewShapeCircle(x, y, radius float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewCircle(x, y, radius), Effects: eff}
}
func NewShapeRectangle(x, y, width, height, angle float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewRectangle(x, y, width, height, angle), Effects: eff}
}
func NewShapeRoundedRectangle(x, y, width, height, angle, roundness float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewRoundedRectangle(x, y, width, height, angle, roundness), Effects: eff}
}
func NewShapeCapsule(x1, y1, x2, y2, radius float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewCapsule(x1, y1, x2, y2, radius), Effects: eff}
}
func NewShapeLine(x1, y1, x2, y2, thickness float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewLine(x1, y1, x2, y2, thickness), Effects: eff}
}

func NewImage(x, y, scale float32, imageId assets.ImageId) Object {
	var _, _, w, h = imageId.CropArea()
	var eff = Effects(internal.DefaultEffects)
	return Object{Shape: geometry.NewRectangle(x, y, float32(w)*scale, float32(h)*scale, 0), ImageId: imageId, Effects: eff}
}
func NewTextbox(x, y, width, height float32, fontId assets.FontId, text ...any) Object {
	var rect = geometry.NewRectangle(x, y, width, height, 0)
	var eff = Effects(internal.DefaultEffects)
	eff.FillColor = palette.DarkGray
	return Object{Shape: rect, TextFontId: fontId, Text: txt.New(text...), Effects: eff}
}
func NewTilemap(scale float32, atlasId assets.TileAtlasId, layerId assets.TileLayerId) Object {
	var tilemap = Object{TileAtlasId: atlasId, TileLayerId: layerId, Effects: Effects(internal.DefaultEffects)}
	var atlas = internal.TileAtlases[uint8(atlasId)]
	var data = internal.TileLayers[uint8(layerId)]
	if atlas != nil && data != nil && data.Image != nil {
		tilemap.Width = float32(data.Image.Width*int32(atlas.TileSize)) * scale
		tilemap.Height = float32(data.Image.Height*int32(atlas.TileSize)) * scale
	}
	return tilemap
}

//=================================================================

func (o *Object) ViewFit(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = min(cw/o.Width, ch/o.Height)
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scale, y-(0.5)*o.Height*scale, 0
}
func (o *Object) ViewFill(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = max(cw/o.Width, ch/o.Height)
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scale, y-(0.5)*o.Height*scale, 0
}
func (o *Object) ViewStretch(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scaleX, scaleY = cw / o.Width, ch / o.Height
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scaleX, y-(0.5)*o.Height*scaleY, 0
}

//=================================================================

func (o *Object) PointToLocal(x, y float32) (localX, localY float32) {
	var dx, dy = x - o.X, y - o.Y
	var sinL, cosL = internal.SinCos(-o.Angle)
	return (dx*cosL - dy*sinL) + 0.5*o.Width, (dx*sinL + dy*cosL) + 0.5*o.Height
}
func (o *Object) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX, locY = localX - (0.5 * o.Width), localY - (0.5 * o.Height)
	var sinL, cosL = internal.SinCos(o.Angle)
	return (locX*cosL - locY*sinL) + o.X, (locX*sinL + locY*cosL) + o.Y
}
func (o *Object) ContainsPoint(x, y float32) bool {
	var lx, ly = o.PointToLocal(x, y)
	return lx >= 0 && ly >= 0 && lx < o.Width && ly < o.Height
}
func (o *Object) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	return o.PointToGlobal(o.Width*edgeX, o.Height*edgeY)
}

// text ===========================================================

// Works only when TextIsBatched is true. Needs to be called when the textbox visuals require a change.
// The update happens upon draw so calling this multiple times per frame is fine.
func (o *Object) TextUpdateBatch() {
	o.textBatches = nil
}

// private ========================================================

func (o *Object) measureLine(fromIndex int, lineHeight float32) (endIndex int, width float32, endLineHeight float32) {
	if fromIndex >= len(o.Text) {
		return fromIndex, 0, lineHeight
	}

	var originalLineHeight = o.Effects.TextLineHeight
	var gapX = o.Effects.TextSymbolGap * originalLineHeight / 255
	var font = internal.Fonts[uint8(o.TextFontId)]
	var x, totalWidth float32
	var prevGlyph internal.Glyph

	for i, r := range o.Text[fromIndex:] {
		if r == '\n' {
			return fromIndex + i, max(totalWidth, x), lineHeight
		}

		var sz = sizes[r]
		if sz != 0 {
			lineHeight = originalLineHeight * sz
			continue
		}

		x += prevGlyph.Kernings[r] * lineHeight
		var offsetX, _, w, _ = o.TextFontId.SymbolArea(r, lineHeight)
		var glyph = font.Chars[r]

		if o.Effects.TextWordWrap && r == ' ' {
			var wX, wTotal float32
			var wPrev internal.Glyph
			var wHeight = lineHeight
			for _, wr := range o.Text[fromIndex+i+1:] {
				if wr == ' ' || wr == '\n' {
					break
				}
				var wsz = sizes[wr]
				if wsz != 0 {
					wHeight = originalLineHeight * wsz
					continue
				}
				var wOffX, _, wW, _ = o.TextFontId.SymbolArea(wr, wHeight)
				var wGlyph = font.Chars[wr]
				wX += wPrev.Kernings[wr] * wHeight
				wPrev, wTotal = wGlyph, max(wX+wOffX+wW, wTotal)
				wX += wGlyph.Advance*wHeight + gapX
			}
			if x+glyph.Advance*lineHeight+gapX+max(wTotal, wX) > o.Width {
				return fromIndex + i, max(totalWidth, x), lineHeight
			}
		}
		x, prevGlyph, totalWidth = x+(glyph.Advance*lineHeight+gapX), glyph, max(x+offsetX+w, totalWidth)
	}
	return len(o.Text), max(totalWidth, x), lineHeight
}

//=================================================================

func (o *Object) TilemapSet(column, row int, tile Tile) {
	o.TilemapSetArea(column, row, 1, 1, tile)
}
func (o *Object) TilemapSetArea(column, row, width, height int, tile Tile) {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	var atlas = internal.TileAtlases[uint8(o.TileAtlasId)]
	if layer == nil {
		return
	}

	var packed = uint32(tile.Id&0xFFFF) | uint32(tile.FrameSpeed&0x1F)<<16 | uint32(tile.FrameOffset&0x0F)<<21 |
		uint32(tile.FrameCount&0x0F)<<25 | uint32(tile.Rotations90&0x03)<<29

	if tile.Flip {
		packed |= (1 << 31)
	}

	var r, g = uint8((packed >> 24) & 0xFF), uint8((packed >> 16) & 0xFF)
	var b, a = uint8((packed >> 8) & 0xFF), uint8((packed >> 0) & 0xFF)
	var colr, rect = rl.NewColor(r, g, b, a), rl.NewRectangle(float32(column), float32(row), float32(width), float32(height))
	var w, h = o.TilemapSize()
	var _, cellHasPts = atlas.ShapesPerTile[tile.Id]

	for i := row; i < row+height; i++ {
		for j := column; j < column+width; j++ {
			var prevTile = o.TilemapAtCell(j, i)
			var _, prevCellHasPts = atlas.ShapesPerTile[prevTile.Id]
			if !prevCellHasPts && !cellHasPts {
				continue
			}

			var index1D = number.Indexes2DToIndex1D(j, i, w, h)
			layer.LastDirtyTime = internal.Runtime

			if cellHasPts {
				layer.CellsWithPoints[index1D] = struct{}{}
			} else {
				delete(layer.CellsWithPoints, index1D)
			}
		}
	}

	rl.ImageDrawRectangle(layer.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(layer.Texture, rect, collection.SameItems(width*height, colr))
}

//=================================================================

func (o *Object) TilemapAtCell(column, row int) Tile {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return Tile{}
	}

	var c = rl.GetImageColor(*layer.Image, int32(column), int32(row))
	var packed = uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)

	return Tile{
		Id:          uint16(packed & 0xFFFF),
		FrameSpeed:  byte((packed >> 16) & 0x1F),
		FrameOffset: byte((packed >> 21) & 0x0F),
		FrameCount:  byte((packed >> 25) & 0x0F),
		Rotations90: byte((packed >> 29) & 0x03),
		Flip:        (packed >> 31) == 1,
	}
}

func (o *Object) TilemapShapes() []geometry.Shape {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return nil
	}

	if layer.Image == nil || layer.Texture.Width == 0 { // is object layer
		var result = make([]geometry.Shape, len(layer.Objects))
		for i, s := range layer.Objects {
			result[i] = geometry.Shape{X: s[0], Y: s[1], Width: s[2], Height: s[3], Angle: s[4], Roundness: s[5]}
		}
		return result
	}

	var w, h = o.TilemapSize()
	var result []geometry.Shape
	for cellIndex1D := range layer.CellsWithPoints {
		var row, column = number.Index1DToIndexes2D(cellIndex1D, w, h)
		result = append(result, o.TilemapShapesAtCell(column, row)...)
	}
	return result
}
func (o *Object) TilemapShapesAtCell(column, row int) []geometry.Shape {
	var tile = o.TilemapAtCell(column, row)
	if tile.Id == 0 {
		return nil
	}
	return o.TilemapShapesFromTile(tile.Id)
}
func (o *Object) TilemapShapesFromTile(tileId uint16) []geometry.Shape {
	var atlas = internal.TileAtlases[uint8(o.TileAtlasId)]
	if atlas == nil {
		return nil
	}
	var shapes, has = atlas.ShapesPerTile[tileId]
	if !has {
		return nil
	}
	var result = make([]geometry.Shape, len(shapes))
	for i, s := range shapes {
		result[i] = geometry.Shape{X: s[0], Y: s[1], Width: s[2], Height: s[3], Angle: s[4], Roundness: s[5]}
	}
	return result
}

func (o *Object) TilemapSize() (columns, rows int) {
	return 1, 1
}
func (o *Object) TilemapSizeTile() (width, height float32) {
	var atlas = internal.TileAtlases[uint8(o.TileAtlasId)]
	if atlas == nil {
		return number.NaN(), number.NaN()
	}
	return float32(atlas.TileSize), float32(atlas.TileSize)
}
func (o *Object) TilemapSizeAtlas() (columns, rows int) {
	var atlas = internal.TileAtlases[uint8(o.TileAtlasId)]
	if atlas == nil {
		return 0, 0
	}
	var tw, th = 1, 1
	return tw / atlas.TileSize, th / atlas.TileSize
}

func (o *Object) TilemapTileCount() int {
	var w, h = o.TilemapSizeAtlas()
	return w * h
}
