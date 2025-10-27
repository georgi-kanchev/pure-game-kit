package assets

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadTexture(id, b64 string) string {
	tryCreateWindow()

	var _, has = internal.Textures[id]
	if has {
		UnloadTextures(id)
	}

	var decompressed = storage.DecompressGZIP([]byte(text.FromBase64(b64)))
	var image = rl.LoadImageFromMemory(".png", decompressed, int32(len(decompressed)))
	var tex = rl.LoadTextureFromImage(image)

	internal.Textures[id] = &tex
	rl.UnloadImage(image)
	return id
}
func loadSound(id, b64 string) string {
	tryCreateWindow()
	tryInitAudio()

	var _, has = internal.Sounds[id]
	if has {
		UnloadSound(id)
	}

	var decompressed = storage.DecompressGZIP([]byte(text.FromBase64(b64)))
	var wave = rl.LoadWaveFromMemory(".ogg", decompressed, int32(len(decompressed)))
	var sound = rl.LoadSoundFromWave(wave)
	internal.Sounds[id] = &sound
	rl.UnloadWave(wave)

	return id
}

/*
func printSoundBase64(path string) {
	if path == "" {
		return
	}
	var raw = file.LoadBytes(path)
	var compressed = file.Compress(raw)
	var b64 = text.ToBase64(string(compressed))
	print(b64)
}
func printFontBase64(path string) {
	if path == "" {
		return
	}
	var bytes = file.LoadBytes(path)
	var compressed = file.Compress(bytes)
	var b64 = text.ToBase64(string(compressed))
	print(b64)
}
func printImageBase64(path string) {
	if path == "" {
		return
	}
	var img = rl.LoadImage(path)
	var bytes = rl.ExportImageToMemory(*img, ".png")
	var compressed = file.Compress(bytes)
	var b64 = text.ToBase64(string(compressed))
	print(b64)
}

func Main() {
	printSoundBase64("popup.ogg")
	// printFontBase64("")
	// printImageBase64("default-ui.png")
}
*/
