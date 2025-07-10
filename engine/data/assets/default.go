package assets

import (
	"encoding/base64"
	"pure-kit/engine/internal"
	"pure-kit/engine/storage"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadTexture(id, b64 string) string {
	tryCreateWindow()

	var _, has = internal.Textures[id]
	if has {
		UnloadTextures(id)
	}

	var imgData, _ = base64.StdEncoding.DecodeString(b64)
	var decompressed = storage.Decompress(imgData)
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
		UnloadSounds(id)
	}

	var raw, _ = base64.StdEncoding.DecodeString(b64)
	var decompressed = storage.Decompress(raw)

	var wave = rl.LoadWaveFromMemory(".mp3", decompressed, int32(len(decompressed)))
	var sound = rl.LoadSoundFromWave(wave)
	internal.Sounds[id] = &sound
	rl.UnloadWave(wave)

	return id
}

//	func printSoundBase64(soundPath string) {
//		var raw, _ = os.ReadFile(soundPath)
//		var compressed = storage.Compress(raw)
//		var b64 = base64.StdEncoding.EncodeToString(compressed)
//		print(b64)
//	}
// func printImageBase64(imgPath string) {
// 	var img = rl.LoadImage(imgPath)
// 	var bytes = rl.ExportImageToMemory(*img, ".png")
// 	var compressed = storage.Compress(bytes)
// 	var b64 = base64.StdEncoding.EncodeToString(compressed)

// 	print(b64)
// }
// func Main() {
// 	printImageBase64("default-ui.png")
// }
