package tilemap

import (
	"pure-kit/engine/internal"
	"strconv"
	"strings"
)

type MapProperties struct {
	TileCountX, TileCountY,
	TileWidth, TileHeight int
	ParallaxX, ParallaxY float32
	Name, Class          string
	IsInfinite           bool
	BackgroundColor      uint
}

func ExtractMapProperties(mapId string) *MapProperties {
	var data, has = internal.TiledMaps[mapId]
	var mapProps *MapProperties = &MapProperties{}
	if !has {
		return mapProps
	}

	mapProps.Name, mapProps.Class = mapId, data.Class
	mapProps.TileCountX, mapProps.TileCountY = data.Width, data.Height
	mapProps.TileWidth, mapProps.TileHeight = data.TileWidth, data.TileHeight
	mapProps.ParallaxX, mapProps.ParallaxY = data.ParallaxOriginX, data.ParallaxOriginY
	mapProps.IsInfinite = data.Infinite
	mapProps.BackgroundColor = color(data.BackgroundColor)
	return mapProps
}

//=================================================================
// private

func color(hex string) uint {
	var trimmed = strings.TrimPrefix(hex, "#")

	if len(trimmed) == 6 {
		trimmed += "FF"
	} else if len(trimmed) != 8 {
		return 0
	}

	var value, err = strconv.ParseUint(trimmed, 16, 32)
	if err != nil {
		return 0
	}

	return uint(value)
}
