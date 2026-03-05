/*
Another important package, similar to the number one. It has many helper functions that operate on a string -
transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
to make them more digestible and clarify their API.
*/
package text

import (
	"fmt"
	"reflect"

	"strconv"
	"strings"
)

type Builder struct{ buffer *strings.Builder }

func NewBuilder() *Builder                       { return &Builder{buffer: &strings.Builder{}} }
func (builder *Builder) WriteText(text string)   { builder.buffer.WriteString(text) }
func (builder *Builder) WriteSymbol(symbol rune) { builder.buffer.WriteRune(symbol) }
func (builder *Builder) Clear()                  { builder.buffer.Reset() }
func (builder *Builder) ToText() string          { return builder.buffer.String() }

func New(elements ...any) string {
	var builder = NewBuilder()
	for _, e := range elements {
		switch v := e.(type) {
		case string:
			builder.WriteText(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			builder.WriteText(fmt.Sprintf("%d", v))
		case float32:
			builder.WriteText(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			builder.WriteText(strconv.FormatFloat(v, 'f', -1, 64))
		case fmt.Stringer:
			builder.WriteText(v.String())
		default:
			var value = reflect.ValueOf(e)
			var valueType = value.Type()

			if valueType.Kind() == reflect.Struct {
				builder.WriteText(fmt.Sprintf("%+v", e)) // struct
				continue
			}

			if valueType.Kind() == reflect.Ptr && valueType.Elem().Kind() == reflect.Struct {
				builder.WriteText(fmt.Sprintf("%+v", value.Elem().Interface())) // pointer to struct
				continue
			}

			builder.WriteText(fmt.Sprint(e)) // fallback
		}
	}
	return builder.ToText()
}
