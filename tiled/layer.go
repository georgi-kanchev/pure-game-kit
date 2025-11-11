package tiled

import "pure-game-kit/internal"

type Layer struct {
	Properties map[string]any
	TileIds    []uint32  // used by Tile Layers only
	Objects    []*Object // used by Object Layers only
	TextureId  string    // used by Image Layers only
}

func newLayerTiles(data *internal.LayerTiles, project *Project) *Layer {
	return nil
}
func newLayerObjects(data *internal.LayerObjects, project *Project) *Layer {
	return nil
}
func newLayerImage(data *internal.LayerImage, project *Project) *Layer {
	return nil
}
