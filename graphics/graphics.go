/*
Strictly tied to the window, drawing on it through a camera and converting between the two coordinate systems.
The camera's drawing consists of two categories: primitives and objects. While using the assets for drawing,
the graphical objects are still very lightweight and exist independently of them.
The concept map of the package looks like this:
  - Camera
  - Node, base drawable object
  - Sprite, extends Node
  - Box, extends Sprite
  - TextBox, extends Node
  - Grid, extends Node
*/
package graphics
