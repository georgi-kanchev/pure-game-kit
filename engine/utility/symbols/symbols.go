package symbols

import (
	"encoding/base64"
	"fmt"
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

	var values = []float64{}
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

	var applyOperator = func(val1, val2 float64, op rune) float64 {
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

	var getNumber = func(expr string, i *int) float64 {
		var start = *i
		for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
			(*i)++
		}
		var numStr = expr[start:*i]
		(*i)--
		var val, err = strconv.ParseFloat(numStr, 64)
		if err != nil {
			return math.NaN()
		}
		return val
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
					return math.NaN()
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
					return math.NaN()
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
						return math.NaN()
					}
				}
				operators = append(operators, c)
			}
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

func AreNumber(text string) bool {
	var _, err = strconv.ParseFloat(text, 64)
	return err == nil
}

func PadLeftAndRight(text string, length int, padStr string) string {
	var textLen = Count(text)
	var spaces = length - textLen
	if spaces <= 0 {
		return text
	}
	var left = spaces / 2
	return PadRight(PadLeft(text, textLen+left, padStr), length, padStr)
}
func PadLeft(text string, totalWidth int, padStr string) string {
	var textLen = Count(text)
	var padding = totalWidth - textLen
	if padding <= 0 || padStr == "" {
		return text
	}
	return repeatPad(padStr, padding) + text
}
func PadRight(text string, totalWidth int, padStr string) string {
	var textLen = Count(text)
	var padding = totalWidth - textLen
	if padding <= 0 || padStr == "" {
		return text
	}
	return text + repeatPad(padStr, padding)
}

func Count(text string) int {
	return utf8.RuneCountInString(text)
}

func Reveal(text string, progress float32) string {
	progress = float32(math.Min(1, math.Max(float64(progress), 0)))
	var textLen = Count(text)
	var cutoff = int(math.Round(float64(progress) * float64(textLen)))

	return string([]rune(text)[cutoff:])
}

func New(elements ...any) string {
	var result string
	for _, e := range elements {
		result += fmt.Sprint(e)
	}
	return result
}

func Fit(text string, maxLength int) string {
	if maxLength == 0 {
		return ""
	}

	const indicator = "…"
	var textRunes = []rune(text)
	var indicatorLen = len([]rune(indicator))
	var textLen = len(textRunes)
	var absMax = int(math.Abs(float64(maxLength)))
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

func ToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}
func FromBase64(textBase64 string) string {
	var decodedBytes, err = base64.StdEncoding.DecodeString(textBase64)
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
	var padRunes = []rune(padStr)
	for builder.Len() < totalRunes {
		for _, r := range padRunes {
			builder.WriteRune(r)
			if Count(builder.String()) >= totalRunes {
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

// endregion
