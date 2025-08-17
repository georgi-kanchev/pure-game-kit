package scene

import (
	"fmt"
	"path"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/internal"
	"strconv"
	"strings"
)

type Scene struct {
	width, height, tileWidth, tileHeight int
	parallaxX, parallaxY                 float32
	name, class                          string
	infinite                             bool
	backgroundColor                      uint
	atlases                              []string
	textures                             []string
}

func New(smoothTexture bool, tiledMapId string) Scene {
	var data, has = internal.TiledData[tiledMapId]
	var scene = Scene{}

	if !has {
		return scene
	}

	scene.name, scene.class = tiledMapId, data.Class
	scene.width, scene.height = data.Width, data.Height
	scene.tileWidth, scene.tileHeight = data.TileWidth, data.TileHeight
	scene.parallaxX, scene.parallaxY = data.ParallaxOriginX, data.ParallaxOriginY
	scene.infinite = data.Infinite
	scene.backgroundColor = color(data.BackgroundColor)
	scene.atlases = make([]string, len(data.Tilesets))

	for index, t := range data.Tilesets {
		var img = t.Image
		var textureIds = assets.LoadTextures(smoothTexture, path.Join(data.Directory, img.Source))
		if len(textureIds) == 0 {
			continue
		}

		// same name is a possibility so adding index
		var name = fmt.Sprintf("%v%v%v%v", tiledMapId, "#", index, t.Name)
		var atlasId = internal.LoadTextureAtlas(textureIds[0], name, t.TileWidth, t.TileHeight, t.Spacing)
		scene.atlases[index] = atlasId
		scene.textures = append(scene.textures, textureIds[0])

		for i := 0; i < t.TileCount/t.Columns; i++ {
			for j := 0; j < t.Columns; j++ {
				var rectId = fmt.Sprintf("%v%v%v%v%v%v", atlasId, "[", j, ",", i, "]")
				assets.SetTextureAtlasTile(atlasId, rectId, float32(j), float32(i), 1, 1, 0, false)
			}
		}
	}

	return scene
}

func (scene *Scene) Unload() {
	assets.UnloadTiledData(scene.name)
	assets.RemoveTextureAtlases(scene.atlases...)
	assets.UnloadTextures(scene.textures...)

	fmt.Printf("internal.TiledData: %v\n", internal.TiledData)
	fmt.Printf("internal.Atlases: %v\n", internal.Atlases)
	fmt.Printf("internal.AtlasRects: %v\n", internal.AtlasRects)
	fmt.Printf("internal.Textures: %v\n", internal.Textures)
}

func (scene *Scene) Size() (width, height int) {
	return scene.width, scene.height
}
func (scene *Scene) TileSize() (width, height int) {
	return scene.width, scene.height
}
func (scene *Scene) ParallaxOrigin() (x, y int) {
	return scene.width, scene.height
}
func (scene *Scene) Class() string {
	return scene.class
}
func (scene *Scene) Infinite() bool {
	return scene.infinite
}
func (scene *Scene) BackgroundColor() uint {
	return scene.backgroundColor
}

// #region private
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

// #endregion
