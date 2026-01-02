package example

import "pure-game-kit/data/storage"

type YamlDog struct {
	Name        string `yaml:"name"`
	FavoriteToy string `yaml:"favoritetoy"`
	Breed       string `yaml:"breed,omitempty"` // hidden if empty
}

func StorageYAML() {
	var dogs = []YamlDog{
		{Name: "Buddy", FavoriteToy: "Tennis Ball", Breed: "Beagle"},
		{Name: "Rex", FavoriteToy: "Stick"},
	}
	var yaml = storage.ToYAML(&dogs)
	var loadedDogs []YamlDog
	storage.FromYAML(yaml, &loadedDogs)
}
