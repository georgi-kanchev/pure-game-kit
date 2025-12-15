/*
Tiled is primarily a map-making software but also supports objects, collision shapes, custom properties,
texts and many more things.

The package constructs scenes out of Tiled files and has a somewhat deeply nested API (utmost 4 to 5 levels)
but relies on a simple design. It tries to group, organize and search data as out-of-the-box as possible and
aims to be easy to be used by other packages. It relies on quite a few package dependencies, such as:
utility, assets, geometry, graphics etc. The concept/depth map of the package looks like this:
  - Map (layered data).
  - Tileset (tile graphics), used by Map.
  - Project (custom data templates & re-usables), optionally used by Map/s.
  - - - -
  - Layer (tiles/objects/image), used by Map.
  - Tile (part of/whole image file), used by Tileset.
  - Object (collisions, paths, entities etc), used by Layer & Tile.
*/
package tiled
