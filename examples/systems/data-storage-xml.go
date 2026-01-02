package example

import "pure-game-kit/data/storage"

type XmlDog struct {
	Name        string `yaml:"name"`
	FavoriteToy string `yaml:"favoritetoy"`
	Breed       string `yaml:"breed,omitempty"` // hidden if empty
}

func StorageXML() {
	var dogs = []XmlDog{
		{Name: "Buddy", FavoriteToy: "Tennis Ball", Breed: "Beagle"},
		{Name: "Rex", FavoriteToy: "Stick"},
	}
	var xml = storage.ToXML(&dogs)
	var loadedDogs []XmlDog
	storage.FromXML(xml, &loadedDogs)
}
