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

	data.Crops = text.Replace(data.Crops, " ", "")
	data.Crops = text.Replace(data.Crops, "\n", "")

	var length = text.SplitCount(data.Crops, "|")
	data.CropsList = make([]int32, length)
	data.GroupsMap = make(map[string][]int32, length)

	for i := range length {
		var crop = text.SplitAtIndex(data.Crops, "|", i)
		var x = text.ToNumber[float32](text.SplitAtIndex(crop, ",", 0))
		var y = text.ToNumber[float32](text.SplitAtIndex(crop, ",", 1))
		var w = text.ToNumber[float32](text.SplitAtIndex(crop, ",", 2))
		var h = text.ToNumber[float32](text.SplitAtIndex(crop, ",", 3))
		data.CropsList[i] = int32(LoadImageCrop(imageId, x, y, w, h))
	}

	for a := range data.Groups {
		var frameCount = text.SplitCount(data.Groups[a].CropIndexes, " ")
		data.GroupsMap[data.Groups[a].Name] = make([]int32, frameCount)
		for i := range frameCount {
			var frameIndex = text.ToNumber[int](text.SplitAtIndex(data.Groups[a].CropIndexes, " ", i))
			data.GroupsMap[data.Groups[a].Name][i] = data.CropsList[frameIndex]
		}
	}

	internal.NextAtlasId++
	internal.Atlases[internal.NextAtlasId] = data
	return AtlasId(internal.NextAtlasId)
}

func (t AtlasId) AllCrops() []ImageId {
	var atlas, has = internal.Atlases[uint16(t)]
	if !has {
		return nil
	} // no copy/no allocation cast of []int32 -> []ImageId, pointing at the same data
	return unsafe.Slice((*ImageId)(unsafe.SliceData(atlas.CropsList)), len(atlas.CropsList))
}
func (t AtlasId) Crops(groupName string) []ImageId {
	var atlas1, has1 = internal.Atlases[uint16(t)]
	var atlas2, has2 = atlas1.GroupsMap[groupName]
	if !has1 || !has2 {
		return nil
	} // no copy/no allocation cast of []int32 -> []ImageId, pointing at the same data
	return unsafe.Slice((*ImageId)(unsafe.SliceData(atlas2)), len(atlas2))
}

func (t AtlasId) Unload() {
	delete(internal.Atlases, uint16(t))
}
