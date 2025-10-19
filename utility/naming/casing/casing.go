package casing

const (
	Random   = 1 << iota // RaNDomCasE
	Lower                // lowercase
	Upper                // UPPERCASE
	Camel                // camelCase
	Pascal               // PascalCase
	Sentence             // Sentence case
	Pingpong             // PiNgPoNg CaSe
	Pongping             // pOnGpInG cAsE
	Separated
)
