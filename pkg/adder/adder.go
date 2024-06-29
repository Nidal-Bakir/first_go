package adder

import (
	"errors"
	"fmt"
)

const (
	Plus rune = '+'
	Sub  rune = '-'
	Mul  rune = '*'
	Div  rune = '/'
)

var ErrDivisionByZero error = errors.New("error division by zero")

func addNum(x, y int) int {
	return x + y
}

func MathOp(n1, n2 int, op rune) (int, error) {
	switch op {
	case Plus:
		return n1 + n2, nil
	case Sub:
		return n1 - n2, nil
	case Mul:
		return n1 * n2, nil
	case Div:
		if n2 == 0 {
			return 0, ErrDivisionByZero
		}
		return n1 / n2, nil
	default:
		return 0, fmt.Errorf("unexpected operation: %v", string(op))
	}
}
