package text

import (
	b64 "encoding/base64"
	"fmt"
	"math"
	"pure-game-kit/utility/number"
	"reflect"

	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Builder struct{ buffer *strings.Builder }

func NewBuilder() *Builder                       { return &Builder{buffer: &strings.Builder{}} }
func (builder *Builder) WriteText(text string)   { builder.buffer.WriteString(text) }
func (builder *Builder) WriteSymbol(symbol rune) { builder.buffer.WriteRune(symbol) }
func (builder *Builder) Clear()                  { builder.buffer.Reset() }
func (builder *Builder) ToText() string          { return builder.buffer.String() }

func New(elements ...any) string {
	var result = ""
	for _, e := range elements {
		switch v := e.(type) {
		case string:
			result += v
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case float32:
			result += strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			result += strconv.FormatFloat(v, 'f', -1, 64)
		case fmt.Stringer:
			result += v.String()
		default:
			var value = reflect.ValueOf(e)
			var valueType = value.Type()

			if valueType.Kind() == reflect.Struct {
				result += fmt.Sprintf("%+v", e) // struct
				continue
			}

			if valueType.Kind() == reflect.Ptr && valueType.Elem().Kind() == reflect.Struct {
				result += fmt.Sprintf("%+v", value.Elem().Interface()) // pointer to struct
				continue
			}

			result += fmt.Sprint(e) // fallback
		}
	}
	return result
}
func ToNumber[T number.Number](text string) T {
	var zero T

	switch any(zero).(type) {
	case float32:
		var result, err = strconv.ParseFloat(text, 32)
		if err != nil {
			return T(number.NaN())
		}
		return T(result)
	case float64:
		var result, err = strconv.ParseFloat(text, 64)
		if err != nil {
			return T(math.NaN())
		}
		return T(result)
	case int, int8, int16, int32, int64:
		result, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return zero
		}
		return T(result)
	case uint, uint8, uint16, uint32, uint64:
		result, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return zero
		}
		return T(result)

	default:
		return zero
	}
}

func PadLeftAndRight(text string, length int, pad string) string {
	var textLen = Length(text)
	var spaces = length - textLen
	if spaces <= 0 {
		return text
	}
	var left = spaces / 2
	return PadRight(PadLeft(text, textLen+left, pad), length, pad)
}
func PadLeft(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return repeatPad(pad, padding) + text
}
func PadRight(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return text + repeatPad(pad, padding)
}
func PadZeros(number float32, amountOfZeros int) string {
	if amountOfZeros == 0 {
		return New(number)
	}
	if amountOfZeros < 0 {
		var width = -amountOfZeros
		return fmt.Sprintf("%0*d", width, int(number))
	}
	return fmt.Sprintf("%.*f", amountOfZeros, number)
}

func Reveal(text string, progress float32) string {
	progress = number.Limit(progress, 0, 1)
	var textLen = float32(Length(text))
	var cutoff = int(number.Round(progress*textLen, -1))

	return string([]rune(text)[cutoff:])
}
func Fit(text string, maxLength int) string {
	if maxLength == 0 {
		return ""
	}

	const indicator = "…"
	var textRunes = []rune(text)
	var indicatorLen = len([]rune(indicator))
	var textLen = len(textRunes)
	var absMax = number.Unsign(maxLength)
	var trimLen = absMax - indicatorLen

	if maxLength > 0 && textLen > int(maxLength) {
		if trimLen <= 0 {
			return indicator
		}
		return string(textRunes[:trimLen]) + indicator
	} else if maxLength < 0 && textLen > absMax {
		if trimLen <= 0 {
			return indicator
		}
		return indicator + string(textRunes[textLen-trimLen:])
	}

	return text
}
func Replace(text, parts, with string) string {
	return strings.ReplaceAll(text, parts, with)
}
func Remove(text string, parts ...string) string {
	for _, part := range parts {
		text = Replace(text, part, "")
	}
	return text
}
func Trim(text string) string {
	return TrimRight(TrimLeft(text))
}
func TrimLeft(text string) string {
	return strings.TrimLeft(text, " \r\n")
}
func TrimRight(text string) string {
	return strings.TrimRight(text, " \r\n")
}

func ToBase64(text string) string {
	return b64.StdEncoding.EncodeToString([]byte(text))
}
func FromBase64(base64 string) string {
	var decodedBytes, err = b64.StdEncoding.DecodeString(base64)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

func ByteSize(byteSize int) string {
	const unit = 1024
	if byteSize < unit {
		return fmt.Sprintf("%d B", byteSize)
	}
	var div, exp = int(unit), 0
	for n := byteSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.3f %cB", float32(byteSize)/float32(div), "KMGTPE"[exp])
}
func Calculate(mathExpression string) float32 {
	mathExpression = Remove(mathExpression, " ")

	var values = []float32{}
	var operators = []rune{}
	var bracketCountOpen = 0
	var bracketCountClose = 0

	var isOperator = func(c rune) bool {
		return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '%'
	}

	var priority = func(op rune) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/', '%':
			return 2
		case '^':
			return 3
		default:
			return 0
		}
	}

	var applyOperator = func(val1, val2 float32, op rune) float32 {
		switch op {
		case '+':
			return val1 + val2
		case '-':
			return val1 - val2
		case '*':
			return val1 * val2
		case '/':
			if val2 != 0 {
				return val1 / val2
			}
		case '%':
			if val2 != 0 {
				return number.DivisionRemainder(val1, val2)
			}
		case '^':
			return number.Power(val1, val2)
		}
		return number.NaN()
	}

	var process = func() bool {
		if len(values) < 2 || len(operators) < 1 {
			return true
		}
		var val2 = values[len(values)-1]
		values = values[:len(values)-1]
		var val1 = values[len(values)-1]
		values = values[:len(values)-1]
		var op = operators[len(operators)-1]
		operators = operators[:len(operators)-1]
		values = append(values, applyOperator(val1, val2, op))
		return false
	}

	var getNumber = func(expr string, i *int) float32 {
		var start = *i
		for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
			(*i)++
		}
		var numStr = expr[start:*i]
		(*i)--
		return ToNumber[float32](numStr)
	}

	for i := 0; i < len(mathExpression); i++ {
		var c = rune(mathExpression[i])

		if unicode.IsDigit(c) || c == '.' {
			var val = getNumber(mathExpression, &i)
			values = append(values, val)
		} else if c == '(' {
			operators = append(operators, c)
			bracketCountOpen++
		} else if c == ')' {
			bracketCountClose++
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if process() {
					return number.NaN()
				}
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else if isOperator(c) {
			// Check for unary minus or plus
			if (i == 0) || mathExpression[i-1] == '(' || isOperator(rune(mathExpression[i-1])) {
				// It's a sign, parse the number after it
				i++
				if i >= len(mathExpression) {
					return number.NaN()
				}
				val := getNumber(mathExpression, &i)
				if c == '-' {
					val = -val
				}
				values = append(values, val)
			} else {
				// Normal binary operator
				for len(operators) > 0 && priority(operators[len(operators)-1]) >= priority(c) {
					if process() {
						return number.NaN()
					}
				}
				operators = append(operators, c)
			}
		}

		if bracketCountClose > bracketCountOpen {
			return number.NaN()
		}
	}

	if bracketCountOpen != bracketCountClose {
		return number.NaN()
	}

	for len(operators) > 0 {
		if process() {
			return number.NaN()
		}
	}

	if len(values) == 0 {
		return number.NaN()
	}
	return values[len(values)-1]
}
func OpenURL(url string) {
	rl.OpenURL(url)
}
func Split(text, divider string) []string {
	return strings.Split(text, divider)
}

func Length(text string) int {
	return utf8.RuneCountInString(text)
}
func IndexOf(text, part string) int {
	return strings.Index(text, part)
}
func Contains(text, part string) bool {
	return strings.Contains(text, part)
}
func StartsWith(text, value string) bool {
	return strings.HasPrefix(text, value)
}
func EndsWith(text, value string) bool {
	return strings.HasSuffix(text, value)
}

func LowerCase(text string) string {
	return strings.ToLower(text)
}
func UpperCase(text string) string {
	return strings.ToUpper(text)
}

//=================================================================
// private

func repeatPad(padStr string, totalRunes int) string {
	if padStr == "" {
		return ""
	}
	var builder strings.Builder
	var padRunes = []rune(padStr)
	for builder.Len() < totalRunes {
		for _, r := range padRunes {
			builder.WriteRune(r)
			if Length(builder.String()) >= totalRunes {
				return truncateToRunes(builder.String(), totalRunes)
			}
		}
	}
	return truncateToRunes(builder.String(), totalRunes)
}
func truncateToRunes(s string, maxRunes int) string {
	var builder strings.Builder
	var count = 0
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		builder.WriteRune(r)
		count++
	}
	return builder.String()
}
