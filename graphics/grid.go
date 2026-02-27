package graphics

type Grid struct {
	Node

	TileMapId string
	Effects   *Effects
}

func NewGrid(tileMapId string) *Grid {
	return &Grid{TileMapId: tileMapId}
}

//=================================================================
