package assets

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"pure-kit/engine/internal"
	"strings"
)

func LoadTiledData(tmxFilePaths ...string) []string {
	var resultIds = []string{}
	for _, path := range tmxFilePaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		var scene internal.TiledMap
		var error = xml.NewDecoder(file).Decode(&scene)
		if error != nil {
			continue
		}

		var name = filepath.Base(path)
		scene.Directory = filepath.Dir(path)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		resultIds = append(resultIds, name)
		internal.TiledData[name] = scene
	}
	return resultIds
}

func UnloadTiledData(tiledMapIds ...string) {
	for _, v := range tiledMapIds {
		delete(internal.TiledData, v)
	}
}
