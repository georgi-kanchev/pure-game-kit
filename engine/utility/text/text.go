package text

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func OpenURL(url string) {
	rl.OpenURL(url)
}

func Calculate(mathExpression string) float64 {
	mathExpression = strings.ReplaceAll(mathExpression, " ", "")

	values := []float64{}
	operators := []rune{}
	bracketCountOpen := 0
	bracketCountClose := 0

	isOperator := func(c rune) bool {
		return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '%'
	}

	priority := func(op rune) int {
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

	applyOperator := func(val1, val2 float64, op rune) float64 {
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
				return math.Mod(val1, val2)
			}
		case '^':
			return math.Pow(val1, val2)
		}
		return math.NaN()
	}

	process := func() bool {
		if len(values) < 2 || len(operators) < 1 {
			return true
		}
		val2 := values[len(values)-1]
		values = values[:len(values)-1]
		val1 := values[len(values)-1]
		values = values[:len(values)-1]
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]
		values = append(values, applyOperator(val1, val2, op))
		return false
	}

	getNumber := func(expr string, i *int) float64 {
		start := *i
		for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
			(*i)++
		}
		numStr := expr[start:*i]
		(*i)--
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return math.NaN()
		}
		return val
	}

	for i := 0; i < len(mathExpression); i++ {
		c := rune(mathExpression[i])

		if unicode.IsDigit(c) || c == '.' {
			val := getNumber(mathExpression, &i)
			values = append(values, val)
		} else if c == '(' {
			operators = append(operators, c)
			bracketCountOpen++
		} else if c == ')' {
			bracketCountClose++
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if process() {
					return math.NaN()
				}
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else if isOperator(c) {
			for len(operators) > 0 && priority(operators[len(operators)-1]) >= priority(c) {
				if process() {
					return math.NaN()
				}
			}
			operators = append(operators, c)
		}

		if bracketCountClose > bracketCountOpen {
			return math.NaN()
		}
	}

	if bracketCountOpen != bracketCountClose {
		return math.NaN()
	}

	for len(operators) > 0 {
		if process() {
			return math.NaN()
		}
	}

	if len(values) == 0 {
		return math.NaN()
	}
	return values[len(values)-1]
}

func IsNumber(text string) bool {
	_, err := strconv.ParseFloat(text, 64)
	return err == nil
}

func PadLeftAndRight(text string, length int, padStr string) string {
	textLen := utf8.RuneCountInString(text)
	spaces := length - textLen
	if spaces <= 0 {
		return text
	}
	left := spaces / 2
	return PadRight(PadLeft(text, textLen+left, padStr), length, padStr)
}
func PadLeft(text string, totalWidth int, padStr string) string {
	textLen := utf8.RuneCountInString(text)
	padding := totalWidth - textLen
	if padding <= 0 || padStr == "" {
		return text
	}
	return repeatPad(padStr, padding) + text
}
func PadRight(text string, totalWidth int, padStr string) string {
	textLen := utf8.RuneCountInString(text)
	padding := totalWidth - textLen
	if padding <= 0 || padStr == "" {
		return text
	}
	return text + repeatPad(padStr, padding)
}

func Reveal(text string, progress float32) string {
	progress = float32(math.Min(1, math.Max(float64(progress), 0)))
	textLen := utf8.RuneCountInString(text)
	cutoff := int(math.Round(float64(progress) * float64(textLen)))

	return string([]rune(text)[cutoff:])
}

func Fit(text string, maxLength int) string {
	if maxLength == 0 {
		return ""
	}

	const indicator = "…"
	textRunes := []rune(text)
	indicatorLen := len([]rune(indicator))
	textLen := len(textRunes)
	absMax := int(math.Abs(float64(maxLength)))
	trimLen := absMax - indicatorLen

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

func ToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}
func FromBase64(textBase64 string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(textBase64)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

// region private

func repeatPad(padStr string, totalRunes int) string {
	if padStr == "" {
		return ""
	}
	var builder strings.Builder
	padRunes := []rune(padStr)
	for builder.Len() < totalRunes {
		for _, r := range padRunes {
			builder.WriteRune(r)
			if utf8.RuneCountInString(builder.String()) >= totalRunes {
				return truncateToRunes(builder.String(), totalRunes)
			}
		}
	}
	return truncateToRunes(builder.String(), totalRunes)
}
func truncateToRunes(s string, maxRunes int) string {
	var builder strings.Builder
	count := 0
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		builder.WriteRune(r)
		count++
	}
	return builder.String()
}

// endregion
