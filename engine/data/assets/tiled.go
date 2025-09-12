package assets

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"strings"
)

func LoadTiledTilesets(tsxFilePaths ...string) []string {
	var resultIds = []string{}
	for _, filePath := range tsxFilePaths {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		defer file.Close()

		var tileset *internal.Tileset
		var err2 = xml.NewDecoder(file).Decode(&tileset)
		if err2 != nil {
			continue
		}

		var name = filepath.Base(filePath)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		var id = path.Join(path.Dir(filePath), name)
		resultIds = append(resultIds, id)
		internal.TiledTilesets[id] = tileset

		var texturePath = path.Join(path.Dir(filePath), tileset.Image.Source)
		var textureIds = LoadTextures(texturePath)
		if len(textureIds) == 0 {
			continue
		}

		var atlasId = SetTextureAtlas(textureIds[0], tileset.TileWidth, tileset.TileHeight, tileset.Spacing)
		var w, h = tileset.Columns, tileset.TileCount / tileset.Columns
		tileset.AtlasId = atlasId

		for i := range h {
			for j := range w {
				var index = number.Indexes2DToIndex1D(i, j, w, h)
				var rectId = text.New(atlasId, "[", index, "]")
				SetTextureAtlasTile(atlasId, rectId, float32(j), float32(i), 1, 1, 0, false)
			}
		}
	}
	return resultIds
}

func LoadTiledWorlds(worldFilePaths ...string) []string {
	var resultIds = []string{}
	for _, path := range worldFilePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		defer file.Close()

		var world *internal.World
		var err2 = json.NewDecoder(file).Decode(&world)
		if err2 != nil {
			continue
		}

		var name = filepath.Base(path)
		world.Directory = filepath.Dir(path)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		resultIds = append(resultIds, name)
		world.Name = name
		internal.TiledWorlds[name] = world
	}
	return resultIds
}

func LoadTiledMaps(tmxFilePaths ...string) []string {
	var resultIds = []string{}
	for _, path := range tmxFilePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		defer file.Close()

		var mapData *internal.Map
		var error = xml.NewDecoder(file).Decode(&mapData)
		if error != nil {
			continue
		}

		var name = filepath.Base(path)
		mapData.Directory = filepath.Dir(path)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		resultIds = append(resultIds, name)
		internal.TiledMaps[name] = mapData
	}
	return resultIds
}

func UnloadTiledMaps(tiledMapIds ...string) {
	for _, v := range tiledMapIds {
		delete(internal.TiledMaps, v)
	}
}
