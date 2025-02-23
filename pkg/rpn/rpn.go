package rpn

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/roadtoseniors/apicalc/pkg/stack"
)

const (
	emptyToken = iota
	wrongToken
	numberToken
	operatorToken
	leftBracketToken
	rightBracketToken
)

// преобразуем в обратную польскую запись
func NewRPN(input string) ([]string, error) {
	input = strings.ReplaceAll(input, "+", " + ")
	input = strings.ReplaceAll(input, "-", " - ")
	input = strings.ReplaceAll(input, "*", " * ")
	input = strings.ReplaceAll(input, "/", " / ")
	input = strings.ReplaceAll(input, "(", " ( ")
	input = strings.ReplaceAll(input, ")", " ) ")

	tokens := strings.Fields(input)

	rpnarr := make([]string, 0, len(tokens))

	stack := stack.NewStack[string]()

	predToken := emptyToken
	for _, token := range tokens {
		curToken := emptyToken

		if isOperator(token) {
			if isUnaryOperator(token, predToken) {
				rpnarr = append(rpnarr, "0")
				stack.Push(token)
				continue
			}

			for !stack.Empty() && isOperator(stack.Top()) {
				op := stack.Pop()
				if operatorPriority(op) <= operatorPriority(token) {
					rpnarr = append(rpnarr, op)
				} else {
					stack.Push(op)
					break
				}
			}

			stack.Push(token)
			curToken = operatorToken
		} else if token == "(" {
			stack.Push(token)
			curToken = leftBracketToken
		} else if token == ")" {
			for !stack.Empty() && stack.Top() != "(" {
				rpnarr = append(rpnarr, stack.Pop())
			}
			if stack.Empty() {
				return nil, fmt.Errorf("error: unpaired brackets")
			}
			stack.Pop()
			curToken = rightBracketToken
		} else {
			_, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, fmt.Errorf("incorrect token: '%s'", token)
			}
			rpnarr = append(rpnarr, token)
			curToken = numberToken
		}
		if !checkTokens(predToken, curToken) {
			return nil, fmt.Errorf("incorrect sequence near token: '%s'", token)
		}
		predToken = curToken
	}

	for !stack.Empty() {
		token := stack.Pop()
		if token == "(" {
			return nil, fmt.Errorf("error: unpaired brackets")
		}
		rpnarr = append(rpnarr, token)
	}

	if predToken != numberToken && predToken != rightBracketToken {
		return nil, fmt.Errorf("incorrect sequence near last token")
	}
	return rpnarr, nil
}

// приоритет оператора
func operatorPriority(op string) int {
	switch op {
	case "*", "/":
		return 1
	case "+", "-":
		return 2
	default:
		return -1
	}
}

// является ли оператором
func isOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

// проверка на унарность
func isUnaryOperator(op string, predToken int) bool {
	return (op == "-" || op == "+") &&
		(predToken == emptyToken || predToken == leftBracketToken || predToken == operatorToken)
}

// проверяем корректность последовательности токенов
func checkTokens(prev, cur int) bool {
	switch cur {
	case numberToken:
		return prev == emptyToken || prev == operatorToken || prev == leftBracketToken
	case leftBracketToken:
		return prev == emptyToken || prev == operatorToken || prev == leftBracketToken
	case rightBracketToken:
		return prev == numberToken || prev == rightBracketToken
	case operatorToken:
		return prev == numberToken || prev == rightBracketToken
	default:
		return false
	}
}
