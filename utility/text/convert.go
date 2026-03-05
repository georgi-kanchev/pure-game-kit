package text

import (
	b64 "encoding/base64"
	"fmt"
	"math"
	"pure-game-kit/utility/number"
	"regexp"
	"strconv"
	"strings"
)

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
	if text == "" {
		return nil
	}
	return strings.Split(text, divider)
}
func SplitLines(text string) []string {
	return Split(text, "\n")
}
