/*
Tiled is primarily a map-making software but also supports objects, collision shapes, custom properties,
texts and many more things.

This package constructs Scenes out of Tiled assets. The workflow looks like this:

	Tiled                           // Software
	Files                           // OS
	Raw Assets                      // Engine (may reload from files)
	Projects/Scenes/Tilesets        // Engine (may reload from assets)
	Graphics/Geometry/Gameplay Data // Game (not cached in Engine, always recalculated)

This intermediate level between raw assets and usable contents acts as a data dump of anything that
Tiled exports. It helps extract whatever the game engine needs to render & calculate collisions, as well as
whatever the game needs as custom gameplay properties.

Since it has all the data, the package has a somewhat deeply nested API but relies on a simple design and
tries to group, organize & search data as out-of-the-box as possible.
It relies on quite a few package dependencies, such as: utility, assets, geometry, graphics etc.
The concept layout of the package looks like this:

	Project
	|
	+-> Reusable properties & templates
	|
	+-> Tilesets (unique)
	    |
	    +-> Tiles
	        |
	        +-> Objects
	            |
	            +-> Points

	Scenes (may exist without Projects)
	|
	+-> Used Tilesets (may be unique or point to Project ones)
	|
	+-> Layers
	    |
	    +-> Tile IDs (points to Tileset Tiles)
	    |
	    +-> Objects
	    |    |
	    |    +-> Points
	    |
	    +-> Image
*/
package tiled
