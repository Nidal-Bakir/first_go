package adder

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("startup call to main")
	exitCode := m.Run()
	fmt.Println("cleanup call to main")
	os.Exit(exitCode)
}
func Test_addNum(t *testing.T) {
	result := addNum(2, 2)
	assert.Equal(t, 4, result)

	result = addNum(6, 6)
	assert.Equal(t, 12, result)

}

func TestTableTest(t *testing.T) {

	data := []struct {
		name           string
		n1             int
		n2             int
		op             rune
		expectedResult int
		expectedErrStr string
	}{
		{name: "add two numbers", n1: 2, n2: 2, op: Plus, expectedResult: 4, expectedErrStr: ""},
		{name: "sub two numbers", n1: 3, n2: 2, op: Sub, expectedResult: 1, expectedErrStr: ""},
		{name: "mul two numbers", n1: 3, n2: 2, op: Mul, expectedResult: 6, expectedErrStr: ""},
		{name: "Div two numbers", n1: 3, n2: 2, op: Div, expectedResult: 1, expectedErrStr: ""},
		{name: "Div two numbers when n2 equal zero", n1: 3, n2: 0, op: Div, expectedResult: 0, expectedErrStr: ErrDivisionByZero.Error()},
		{name: "should return error when unexpected operator is passed", n1: 3, n2: 5, op: 'G', expectedResult: 0, expectedErrStr: fmt.Sprintf("unexpected operation: %v", string('G'))},
	}

	for _, v := range data {
		t.Run(v.name, func(t *testing.T) {
			actual, err := MathOp(v.n1, v.n2, v.op)
			assert.Equal(t, v.expectedResult, actual)

			if err != nil {
				assert.Equal(t, v.expectedErrStr, err.Error())
			}
		})
	}
}
