// Another important package, similar to the number one. It has many helper functions that operate on a string -
// transformations, checks, formatting or executing a result on it. Also wraps some standard string functions
// to make them more digestible and clarify their API.
package text

import (
	"fmt"
	"math"
	"pure-game-kit/packages/utility/number"
	"reflect"
	"strings"
	"unicode"

	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func New(elements ...any) string {
	builder.Reset()
	for _, e := range elements {
		switch v := e.(type) {
		case string:
			builder.WriteString(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			fmt.Fprintf(&builder, "%d", v)
		case float32:
			builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			if v {
				builder.WriteString("true")
			} else {
				builder.WriteString("false")
			}
		case fmt.Stringer:
			builder.WriteString(v.String())
		default:
			var value = reflect.ValueOf(e)
			var valueType = value.Type()
			if valueType.Kind() == reflect.Struct {
				fmt.Fprintf(&builder, "%+v", e) // struct
				continue
			}
			if valueType.Kind() == reflect.Pointer && valueType.Elem().Kind() == reflect.Struct {
				fmt.Fprintf(&builder, "%+v", value.Elem().Interface()) // pointer to struct
				continue
			}
			fmt.Fprint(&builder, e) // fallback
		}
	}
	return builder.String()
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

func Calculate(mathExpression string, variables func(string) float32) float32 {
	mathExpression = Remove(mathExpression, " ")
	calcValues = calcValues[:0]
	calcOperators = calcOperators[:0]
	var bracketCountOpen, bracketCountClose int

	for i := 0; i < len(mathExpression); i++ {
		var c = rune(mathExpression[i])

		if unicode.IsDigit(c) || c == '.' {
			calcValues = append(calcValues, calcGetNumber(mathExpression, &i))
		} else if c == '(' {
			calcOperators = append(calcOperators, c)
			bracketCountOpen++
		} else if c == ')' {
			bracketCountClose++
			for len(calcOperators) > 0 && calcOperators[len(calcOperators)-1] != '(' {
				if calcProcess() {
					return number.NaN()
				}
			}
			if len(calcOperators) > 0 {
				calcOperators = calcOperators[:len(calcOperators)-1]
			}
		} else if calcIsOperator(c) {
			// Check for unary minus or plus
			if (i == 0) || mathExpression[i-1] == '(' || calcIsOperator(rune(mathExpression[i-1])) {
				i++
				if i >= len(mathExpression) {
					return number.NaN()
				}
				val := calcGetNumber(mathExpression, &i)
				if c == '-' {
					val = -val
				}
				calcValues = append(calcValues, val)
			} else {
				// Normal binary operator
				for len(calcOperators) > 0 && calcPriority(calcOperators[len(calcOperators)-1]) >= calcPriority(c) {
					if calcProcess() {
						return number.NaN()
					}
				}
				calcOperators = append(calcOperators, c)
			}
		} else if unicode.IsLetter(c) {
			var start = i
			for i < len(mathExpression) && (unicode.IsLetter(rune(mathExpression[i])) || unicode.IsDigit(rune(mathExpression[i]))) {
				i++
			}
			var name = mathExpression[start:i]
			i--
			if variables == nil {
				return number.NaN()
			}
			var v = variables(name)
			if number.IsNaN(v) {
				return number.NaN()
			}
			calcValues = append(calcValues, v)
		}

		if bracketCountClose > bracketCountOpen {
			return number.NaN()
		}
	}

	if bracketCountOpen != bracketCountClose {
		return number.NaN()
	}

	for len(calcOperators) > 0 {
		if calcProcess() {
			return number.NaN()
		}
	}

	if len(calcValues) == 0 {
		return number.NaN()
	}
	return calcValues[len(calcValues)-1]
}
func OpenURL(url string) {
	rl.OpenURL(url)
}

// private ========================================================

var builder strings.Builder
var calcValues = make([]float32, 0, 8)
var calcOperators = make([]rune, 0, 8)

func repeatPad(padStr string, totalRunes int) string {
	if padStr == "" {
		return ""
	}
	var padRunes = []rune(padStr)
	builder.Reset()
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
	var count = 0
	builder.Reset()
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		builder.WriteRune(r)
		count++
	}
	return builder.String()
}
func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == '_' || r == '-' || r == '/' || r == '.'
}

func calcIsOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '%'
}
func calcPriority(op rune) int {
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
func calcApplyOp(val1, val2 float32, op rune) float32 {
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
func calcProcess() bool {
	if len(calcValues) < 2 || len(calcOperators) < 1 {
		return true
	}
	var val2 = calcValues[len(calcValues)-1]
	calcValues = calcValues[:len(calcValues)-1]
	var val1 = calcValues[len(calcValues)-1]
	calcValues = calcValues[:len(calcValues)-1]
	var op = calcOperators[len(calcOperators)-1]
	calcOperators = calcOperators[:len(calcOperators)-1]
	calcValues = append(calcValues, calcApplyOp(val1, val2, op))
	return false
}
func calcGetNumber(expr string, i *int) float32 {
	var start = *i
	for *i < len(expr) && (unicode.IsDigit(rune(expr[*i])) || expr[*i] == '.') {
		(*i)++
	}
	var numStr = expr[start:*i]
	(*i)--
	return ToNumber[float32](numStr)
}
