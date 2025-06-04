package assets

import (
	"path/filepath"
	"pure-kit/engine/data/file"
	"pure-kit/engine/data/folder"
	"pure-kit/engine/window"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TextureIsLoaded(id string) bool {
	var _, contains = textures[id]
	return contains
}
func TexturesLoaded() []string {
	var keys []string
	for k := range textures {
		keys = append(keys, k)
	}
	return keys
}

func LoadTextures(filePaths ...string) {
	if !rl.IsWindowReady() {
		window.Recreate()
	}

	for _, path := range filePaths {
		if !file.Exists(path) {
			continue
		}

		var absolutePath = filepath.Join(folder.PathOfExecutable(), path)
		var texture = rl.LoadTexture(absolutePath)
		path = strings.ReplaceAll(path, file.Extension(path), "")

		if texture.Width != 0 {
			textures[path] = texture
		}

	}
}

// region private

var textures = make(map[string]rl.Texture2D)

// endregion
