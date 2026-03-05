package text

import "strings"

type Builder struct{ buffer *strings.Builder }

func NewBuilder() *Builder                       { return &Builder{buffer: &strings.Builder{}} }
func (builder *Builder) WriteText(text string)   { builder.buffer.WriteString(text) }
func (builder *Builder) WriteSymbol(symbol rune) { builder.buffer.WriteRune(symbol) }
func (builder *Builder) Clear()                  { builder.buffer.Reset() }
func (builder *Builder) ToText() string          { return builder.buffer.String() }
