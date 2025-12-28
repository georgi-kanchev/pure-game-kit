/*
Another important package, similar to the number one. It has many helper functions that operate on a string -
transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
to make them more digestible and clarify their API.
*/
package text

import (
	b64 "encoding/base64"
	"fmt"
	"math"
	"pure-game-kit/utility/number"
	"reflect"
	"regexp"

	"strconv"
	"strings"
	"unicode/utf8"
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

//=================================================================
// convert

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
func FormatByteSize(byteSize int) string {
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

// Split text into words (handles spaces, underscores and dashes).
func SplitWords(text string) []string {
	re := regexp.MustCompile(`[a-z0-9]+|[A-Z][a-z0-9]*|[A-Z]+(?=[A-Z][a-z0-9]|$)`)
	return re.FindAllString(strings.ToLower(text), -1)
}
func Split(text, divider string) []string {
	return strings.Split(text, divider)
}

//=================================================================
// edit

func ReplaceWith(text, part, with string) string {
	return strings.ReplaceAll(text, part, with)
}
func Remove(text string, parts ...string) string {
	for _, part := range parts {
		text = ReplaceWith(text, part, "")
	}
	return text
}

// Progress 0..1 for start-to-end and 0..-1 for end-to-start.
func Reveal(text string, progress float32) string {
	var textLen = float32(Length(text))
	if progress >= 0 {
		progress = number.Limit(progress, 0, 1)
		var cutoff = int(number.Round(progress * textLen))
		return string([]rune(text)[cutoff:])
	}
	progress = number.Limit(progress, -1, 0)
	var cutoff = int(number.Round(progress * textLen))
	return string([]rune(text)[cutoff:])
}
func RevealBy(text string, fromStart, fromEnd int) string {
	var runes = []rune(text)
	var length = len(runes)
	var start = number.Limit(fromStart, 0, length)
	var end = number.Limit(fromEnd, 0, length)
	return string(runes[start:end])
}

// Positive length trims from the end, negative length trims from the start. Default indicator if skipped: '…'
func Limit(text string, length int, indicator ...string) string {
	if length == 0 {
		return ""
	}

	var ind = "…"
	if len(indicator) > 0 {
		ind = indicator[0]
	}
	var textRunes = []rune(text)
	var indicatorLen = len([]rune(ind))
	var textLen = len(textRunes)
	var absMax = number.Unsign(length)
	var trimLen = absMax - indicatorLen

	if length > 0 && textLen > int(length) {
		if trimLen <= 0 {
			return ind
		}
		return string(textRunes[:trimLen]) + ind
	} else if length < 0 && textLen > absMax {
		if trimLen <= 0 {
			return ind
		}
		return ind + string(textRunes[textLen-trimLen:])
	}

	return text
}

// Returns the text before the part. If the part is not found, returns the original text.
func Before(text, part string) string {
	var pos = strings.Index(text, part)
	if pos == -1 {
		return text
	}
	return text[:pos]
}

// Returns the text after the part. If the part is not found, returns the original text.
func After(text, part string) string {
	pos := strings.Index(text, part)
	if pos == -1 {
		return text
	}
	return text[pos+len(part):]
}

// Returns the text between the first part and the second part. If one of the parts is not found,
// returns the original text.
func Between(text, firstPart, secondPart string) string {
	sPos := strings.Index(text, firstPart)
	if sPos == -1 {
		return text
	}

	// Start searching for the endAnchor AFTER the startAnchor
	remaining := text[sPos+len(firstPart):]
	ePos := strings.Index(remaining, secondPart)
	if ePos == -1 {
		return text
	}

	return remaining[:ePos]
}

// Returns a text consisting of the provided amount of copies of the text. If count is 0 or negative, returns "".
func Repeat(text string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(text, count)
}

// Breaks the text at the provided line length, regardless of word boundaries.
func Wrap(text string, lineLength int) string {
	if lineLength <= 0 || len(text) <= lineLength {
		return text
	}

	var sb strings.Builder
	var runes = []rune(text)

	for i, r := range runes {
		sb.WriteRune(r)
		var reachedEnd = (i+1)%lineLength == 0
		var notLastChar = i+1 < len(runes)
		if reachedEnd && notLastChar {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// Breaks the text into lines with a maximum length but only breaks at whitespace to keep words intact.
func WrapWords(text string, lineLength int) string {
	if lineLength <= 0 || len(text) <= lineLength {
		return text
	}

	var words = strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var sb strings.Builder
	var currentLineLength = 0

	for i, word := range words {
		var wordLen = len([]rune(word))
		if currentLineLength+wordLen > lineLength && currentLineLength > 0 { // longer than line length
			sb.WriteByte('\n')
			currentLineLength = 0
		}
		if currentLineLength > 0 {
			sb.WriteByte(' ') // add a space if it's not the start of a new line
			currentLineLength++
		}

		sb.WriteString(word)
		currentLineLength += wordLen

		if currentLineLength >= lineLength && i < len(words)-1 { // still longer than line length
			sb.WriteByte('\n') // the next word MUST start on a new line
			currentLineLength = 0
		}
	}
	return sb.String()
}

//=================================================================
// checks

func Length(text string) int {
	return utf8.RuneCountInString(text)
}
func IndexOf(text, part string) int {
	return strings.Index(text, part)
}
func Contains(text string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(text, part) {
			return false
		}
	}
	return true
}
func IsBlank(text string) bool {
	return Trim(text) == ""
}

//=================================================================
// start/end

func Pad(text string, length int, pad string) string {
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

func Trim(text string) string {
	return TrimEnd(TrimStart(text))
}
func TrimStart(text string) string {
	return strings.TrimLeft(text, " \r\n")
}
func TrimEnd(text string) string {
	return strings.TrimRight(text, " \r\n")
}

// Surrounds a text with the given start part and end part.
// If end part is empty, it uses the start part for both sides.
func SurroundWith(text, startPart string, endPart ...string) string {
	var end = startPart
	if len(endPart) > 0 {
		end = endPart[0]
	}
	return startPart + text + end
}

// Adds the part to the start of the text only if it doesn't already have it.
func EnsureStart(text, part string) string {
	if !strings.HasPrefix(text, part) {
		return part + text
	}
	return text
}

// Adds the part to the end of the text only if it doesn't already have it.
func EnsureEnd(text, part string) string {
	if !strings.HasSuffix(text, part) {
		return text + part
	}
	return text
}

// Removes the given start part and end part only if both are present.
// If end part is empty, it looks for the start part on both sides.
func Chop(text, startPart string, endPart ...string) string {
	var end = startPart
	if len(endPart) > 0 {
		end = endPart[0]
	}

	if strings.HasPrefix(text, startPart) && strings.HasSuffix(text, end) {
		return text[len(startPart) : len(text)-len(end)]
	}
	return text
}

// Removes the start part from the text if it exists.
func ChopStart(text, part string) string {
	return strings.TrimPrefix(text, part)
}

// Removes the end part from the text if it exists.
func ChopEnd(text, part string) string {
	return strings.TrimSuffix(text, part)
}
func StartsWith(text, value string) bool {
	return strings.HasPrefix(text, value)
}
func EndsWith(text, value string) bool {
	return strings.HasSuffix(text, value)
}
