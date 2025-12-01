package blend

const (
	Alpha            = iota // Considering alpha (default)
	Additive                // Adding colors
	Multiplied              // Multiplying colors
	AddColors               // Adding colors (alternative)
	SubtractColors          // Subtracting colors (alternative)
	AlphaPremultiply        // Premultiplied considering alpha
)
