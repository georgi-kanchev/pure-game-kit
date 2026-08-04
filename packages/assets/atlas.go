package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
	"pure-game-kit/packages/utility/text"
	"unsafe"
)

type AtlasId uint16

func LoadAtlas(imageId ImageId, xmlPath string) AtlasId {
	var data = internal.AtlasData{}
	storage.FromXML(file.LoadText(xmlPath), &data)
	if data.XMLName.Local == "" {
		return 0
	}

	data.Map = make(map[string][]int32, len(data.Groups))

	for a := range data.Groups {
		var frameCount = text.SplitCount(data.Groups[a].CropIndexes, " ")
		data.Map[data.Groups[a].Name] = make([]int32, frameCount)
		for i := range frameCount {
			var frameIndex = text.ToNumber[int](text.SplitAtIndex(data.Groups[a].CropIndexes, " ", i))
			var fr = data.Crops[frameIndex]
			var cropId = LoadImageCrop(imageId, float32(fr.X), float32(fr.Y), float32(fr.W), float32(fr.H))
			data.Map[data.Groups[a].Name][i] = int32(cropId)
		}
	}

	internal.NextAtlasId++
	internal.Atlases[internal.NextAtlasId] = data
	return AtlasId(internal.NextAtlasId)
}

func (t AtlasId) CropCount(groupName string) int {
	var atlases, has = internal.Atlases[uint16(t)]
	if !has {
		return 0
	}
	return len(atlases.Map[groupName])
}
func (t AtlasId) Crop(groupName string, index int) ImageId {
	var atlases, has = internal.Atlases[uint16(t)]
	if !has {
		return 0
	}
	var atlas, has2 = atlases.Map[groupName]
	if !has2 || index < 0 || index >= len(atlas) {
		return 0
	}
	return ImageId(atlas[index])
}
func (t AtlasId) Crops(groupName string) []ImageId {
	var atlases, has = internal.Atlases[uint16(t)]
	if !has {
		return nil
	}
	var atlas, has2 = atlases.Map[groupName]
	if !has2 {
		return nil
	} // no copy/no allocation cast of []int32 -> []ImageId, pointing at the same data
	return unsafe.Slice((*ImageId)(unsafe.SliceData(atlas)), len(atlas))
}

func (t AtlasId) Unload() {
	delete(internal.Atlases, uint16(t))
}
