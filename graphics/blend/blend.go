// Used by the graphics camera to define how to draw something. Useful for effects such as lighting.
package blend

const (
	Alpha            = iota // Considering alpha (default)
	Additive                // Adding colors
	Multiplied              // Multiplying colors
	AddColors               // Adding colors (alternative)
	SubtractColors          // Subtracting colors (alternative)
	AlphaPremultiply        // Premultiplied considering alpha
)
