package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
	"pure-game-kit/packages/utility/text"
)

type AnimationsId uint16

func LoadAnimations(imageId ImageId, xmlPath string) AnimationsId {
	var data = internal.AnimationsData{}
	storage.FromXML(file.LoadText(xmlPath), &data)
	if data.XMLName.Local == "" {
		return 0
	}

	data.Map = make(map[string][]int32, len(data.Animations))

	for a := range data.Animations {
		var frameCount = text.SplitCount(data.Animations[a].Frames, " ")
		data.Map[data.Animations[a].Name] = make([]int32, frameCount)
		for i := range frameCount {
			var frameIndex = text.ToNumber[int](text.SplitAtIndex(data.Animations[a].Frames, " ", i))
			var fr = data.Frames[frameIndex]
			var cropId = LoadImageCrop(imageId, float32(fr.X), float32(fr.Y), float32(fr.W), float32(fr.H))
			data.Map[data.Animations[a].Name][i] = int32(cropId)
		}
	}

	internal.NextAnimationsId++
	internal.Animations[internal.NextAnimationsId] = data
	return AnimationsId(internal.NextAnimationsId)
}

func (t AnimationsId) FrameCount(animationName string) int {
	var anims, has = internal.Animations[uint16(t)]
	if !has {
		return 0
	}
	return len(anims.Map[animationName])
}
func (t AnimationsId) Frame(animationName string, index int) ImageId {
	var anims, has = internal.Animations[uint16(t)]
	if !has {
		return 0
	}
	var anim, has2 = anims.Map[animationName]
	if !has2 || index < 0 || index >= len(anim) {
		return 0
	}
	return ImageId(anim[index])
}

func (t AnimationsId) Unload() {
	delete(internal.Animations, uint16(t))
}
