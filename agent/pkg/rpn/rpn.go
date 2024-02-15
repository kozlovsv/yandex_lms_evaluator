package rpn

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kozlovsv/evaluator/agent/pkg/models"
)

var priorities = map[string]uint8{
	"+": 3,
	"-": 3,
	"*": 4,
	"/": 4,
}

var numbers = "0123456789"

func Evaluate(s string, settings models.Settinsg) (float64, error) {
	rpn, err := convertToRpn(s)
	if err != nil {
		return 0, err
	}
	return evaluateRPN(rpn, settings)
}

// New creates a new Reverse Polish Notation with a string pattern.
func convertToRpn(s string) ([]string, error) {
	s = strings.ReplaceAll(s, " ", "")
	nop, n := 0, len(s)
	operators := make([]string, n)
	notation := []string{}

	i := 0
	popPushOp := func(op string) {
		priority := priorities[op]
		for nop > 0 && priorities[operators[nop-1]] >= priority {
			nop--
			notation = append(notation, operators[nop])
		}
		operators[nop] = op
		nop++
		i++
	}

	for i < n {
		c := s[i]
		switch c {
		case ')':
			for nop > 0 && operators[nop-1] != "(" {
				nop--
				notation = append(notation, operators[nop])
			}
			if nop == 0 || operators[nop-1] != "(" {
				return nil, fmt.Errorf("'%v' has no '(' found for ')' at %v", s, i)
			}
			nop--
			if nop > 0 && operators[nop-1] == "!" {
				notation = append(notation, "!")
				nop--
			}
			i++
		case '(':
			operators[nop] = "("
			nop++
			i++
		case '*', '/', '+', '-':
			popPushOp(string(c))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			hasDot := false
			num := []byte{}
			for i < n {
				c = s[i]
				if c == '.' {
					if hasDot || len(num) == 0 {
						return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i, c)
					}
					hasDot = true
				} else if !strings.Contains(numbers, string(c)) {
					break
				}
				num = append(num, c)
				i++
			}
			notation = append(notation, string(num))
		default:
			return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i, c)
		}
	}

	for nop > 0 {
		nop--
		if op := operators[nop]; op != "(" {
			notation = append(notation, op)
		}
	}

	return notation, nil
}

func evaluateRPN(expression []string, settings models.Settinsg) (float64, error) {
	stack := make([]float64, 0)

	for _, token := range expression {
		if num, err := strconv.Atoi(token); err == nil {
			stack = append(stack, (float64)(num))
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid RPN expression")
			}

			// Pop the last two operands
			op2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			op1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Perform the operation based on the operator
			switch token {
			case "+":
				stack = append(stack, op1+op2)
				timer := time.NewTimer(time.Duration(settings.OpPlusTime) * time.Millisecond)
				<-timer.C
			case "-":
				stack = append(stack, op1-op2)
				timer := time.NewTimer(time.Duration(settings.OpMinusTime) * time.Millisecond)
				<-timer.C
			case "*":
				stack = append(stack, op1*op2)
				timer := time.NewTimer(time.Duration(settings.OpMultTime) * time.Millisecond)
				<-timer.C
			case "/":
				if op2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				stack = append(stack, op1/op2)
				timer := time.NewTimer(time.Duration(settings.OpDivTime) * time.Millisecond)
				<-timer.C
			default:
				return 0, fmt.Errorf("unknown operator: %s", token)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid RPN expression")
	}

	return stack[0], nil
}
