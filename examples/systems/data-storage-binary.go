package example

import "pure-game-kit/data/storage"

type Animal interface {
	Name() string
}

type privateField struct{ infinitePats bool } // not saved
type Dog struct {
	privateField
	FavoriteToy string
}
type Cat struct {
	privateField
	Meows int
}
type Horse struct {
	privateField
	Speed float32
}

func (dog *Dog) Name() string     { return "dogo" }
func (cat *Cat) Name() string     { return "cato" }
func (horse *Horse) Name() string { return "juan" }

func StorageBinary() {
	var animals = []Animal{
		&Dog{FavoriteToy: "sticks", privateField: privateField{infinitePats: true}},
		&Cat{Meows: 5, privateField: privateField{infinitePats: true}},
		&Cat{Meows: 23, privateField: privateField{infinitePats: true}},
		&Horse{Speed: 12.3, privateField: privateField{infinitePats: true}}}
	var bytes = storage.ToBytes(&animals, &Dog{}, &Cat{}, &Horse{})
	var newAnimals []Animal
	storage.FromBytes(bytes, &newAnimals, &Dog{}, &Cat{}, &Horse{})
}
