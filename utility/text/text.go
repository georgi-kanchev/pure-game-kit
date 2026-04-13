// Another important package, similar to the number one. It has many helper functions that operate on a string -
// transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
// to make them more digestible and clarify their API.
package text

import (
	"fmt"
	"reflect"
	"strings"

	"strconv"
)

func New(elements ...any) string {
	var builder strings.Builder
	for _, e := range elements {
		switch v := e.(type) {
		case string:
			builder.WriteString(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			builder.WriteString(fmt.Sprintf("%d", v))
		case float32:
			builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		case fmt.Stringer:
			builder.WriteString(v.String())
		default:
			var value = reflect.ValueOf(e)
			var valueType = value.Type()

			if valueType.Kind() == reflect.Struct {
				builder.WriteString(fmt.Sprintf("%+v", e)) // struct
				continue
			}

			if valueType.Kind() == reflect.Ptr && valueType.Elem().Kind() == reflect.Struct {
				builder.WriteString(fmt.Sprintf("%+v", value.Elem().Interface())) // pointer to struct
				continue
			}

			builder.WriteString(fmt.Sprint(e)) // fallback
		}
	}
	return builder.String()
}
