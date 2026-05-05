// Used by the graphics camera to define how to draw something above what's already drawn.
// Useful for effects such as lighting or inverting colors etc.
package blend

const (
	Alpha            = iota // Considering alpha (default)
	Additive                // Adding colors
	Multiply                // Multiplying colors
	AddColors               // Adding colors (alternative)
	SubtractColors          // Subtracting colors (alternative)
	AlphaPremultiply        // Premultiplied considering alpha
)
