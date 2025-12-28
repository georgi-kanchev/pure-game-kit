package text

import (
	"pure-game-kit/utility/number"
	"strings"
	"unicode"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == '_' || r == '-' || r == '/' || r == '.'
}
