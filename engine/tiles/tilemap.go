package tiles

type Map struct {
	Tiles map[[2]int]string
}

func (tilemap *Map) SetTile(cellX, cellY int, tileId string) {
	if tilemap.Tiles == nil {
		tilemap.Tiles = make(map[[2]int]string)
	}

	tilemap.Tiles[[2]int{cellX, cellY}] = tileId
}
