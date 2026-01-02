package example

import "pure-game-kit/data/storage"

type BinaryAnimal interface {
	Name() string
}

type binaryPrivateField struct{ infinitePats bool } // not saved
type BinaryDog struct {
	binaryPrivateField
	FavoriteToy string
}
type BinaryCat struct {
	binaryPrivateField
	Meows int
}
type BinaryHorse struct {
	binaryPrivateField
	Speed float32
}

func (dog *BinaryDog) Name() string     { return "dogo" }
func (cat *BinaryCat) Name() string     { return "cato" }
func (horse *BinaryHorse) Name() string { return "juan" }

func StorageBinary() {
	var animals = []BinaryAnimal{
		&BinaryDog{FavoriteToy: "sticks", binaryPrivateField: binaryPrivateField{infinitePats: true}},
		&BinaryCat{Meows: 5, binaryPrivateField: binaryPrivateField{infinitePats: true}},
		&BinaryCat{Meows: 23, binaryPrivateField: binaryPrivateField{infinitePats: true}},
		&BinaryHorse{Speed: 12.3, binaryPrivateField: binaryPrivateField{infinitePats: true}}}
	var bytes = storage.ToBytes(&animals, &BinaryDog{}, &BinaryCat{}, &BinaryHorse{})
	var newAnimals []BinaryAnimal
	storage.FromBytes(bytes, &newAnimals, &BinaryDog{}, &BinaryCat{}, &BinaryHorse{})
}
