package example

import "pure-game-kit/packages/utility/storage"

type JsonDog struct {
	Name        string `json:"name"`
	FavoriteToy string `json:"favoritetoy"`
	Breed       string `json:"breed,omitempty"`
}

func StorageJSON() {
	var dogs = []JsonDog{
		{Name: "Buddy", FavoriteToy: "Tennis Ball", Breed: "Beagle"},
		{Name: "Rex", FavoriteToy: "Stick"},
	}
	var json = storage.ToJSON(&dogs)
	var loadedDogs []JsonDog
	storage.FromJSON(json, &loadedDogs)
}
