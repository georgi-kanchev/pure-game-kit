package assets

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadTexture(id, b64 string, smooth bool) string {
	tryCreateWindow()

	var _, has = internal.Textures[id]
	if has {
		UnloadTextures(id)
	}

	var raw, _ = base64.StdEncoding.DecodeString(b64)
	var decompressed = decompress(raw)
	var image = rl.LoadImageFromMemory(".png", decompressed, int32(len(decompressed)))
	var tex = rl.LoadTextureFromImage(image)

	if smooth {
		rl.SetTextureFilter(tex, rl.FilterBilinear)
	}

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
	var decompressed = decompress(raw)
	var wave = rl.LoadWaveFromMemory(".mp3", decompressed, int32(len(decompressed)))
	var sound = rl.LoadSoundFromWave(wave)
	internal.Sounds[id] = &sound
	rl.UnloadWave(wave)

	return id
}

func decompress(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = gzip.NewReader(buf)
	if err != nil {
		return data
	}
	defer gr.Close()

	var result, err2 = io.ReadAll(gr)
	if err2 != nil {
		return data
	}

	return result
}
func compress(data []byte) []byte {
	var buf bytes.Buffer
	var gw = gzip.NewWriter(&buf)
	var _, err = gw.Write(data)
	if err != nil {
		return data
	}
	if err := gw.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}

func printSoundBase64(path string) {
	if path == "" {
		return
	}
	var raw, _ = os.ReadFile(path)
	var compressed = compress(raw)
	var b64 = base64.StdEncoding.EncodeToString(compressed)
	print(b64)
}
func printFontBase64(path string) {
	if path == "" {
		return
	}
	var bytes = file.LoadBytes(path)
	var compressed = compress(bytes)
	var b64 = base64.StdEncoding.EncodeToString(compressed)
	print(b64)
}
func printImageBase64(path string) {
	if path == "" {
		return
	}
	var img = rl.LoadImage(path)
	var bytes = rl.ExportImageToMemory(*img, ".png")
	var compressed = compress(bytes)
	var b64 = base64.StdEncoding.EncodeToString(compressed)
	print(b64)
}

func Main() {
	printSoundBase64("")
	printFontBase64("")
	printImageBase64("default-ui.png")
}
