package math_solver

import (
	"context"
	"io"

	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var r io.Reader = strings.NewReader(`1+2
8+7*9
(7+2+
55-8/7`)

type mockMathSolver struct {
	mock.Mock
}

func (ms *mockMathSolver) Resolve(ctx context.Context, exp string) (result float64, err error) {
	args := ms.Called(ctx, exp)
	return args.Get(0).(float64), args.Error(1)
}

func TestMathSolver(t *testing.T) {
	solverMock := new(mockMathSolver)
	p := NewProcessor(solverMock)

	data := []struct {
		q   string
		ans float64
		err error
	}{
		{q: "1+2", ans: 3, err: nil},
		{q: "8+7*9", ans: 71, err: nil},
		{q: "(7+2+", ans: 0, err: ErrSyntaxError},
		{q: "55-8/7", ans: 53.857142857, err: nil},
	}

	for _, v := range data {
		solverMock.On(
			"Resolve",
			context.Background(),
			v.q,
		).Return(v.ans, v.err)

		result, err := p.Solve(context.Background(), r)
		require.Equal(t, v.ans, result)
		if err != nil {
			require.ErrorIs(t, err, v.err)
		}
	}

	solverMock.AssertExpectations(t)

	// when the io reader is consumed, the code should return ErrNoMathToSolve
	// and should NOT call Resolve()
	result, err := p.Solve(context.Background(), r)
	require.Equal(t, 0.0, result)
	require.ErrorIs(t, err, ErrNoMathToSolve)
	solverMock.AssertNotCalled(t, "Resolve", context.Background(), "")
}

func Test_readToNewLine(t *testing.T) {
	l, err := readToNewLine(r)
	require.Equal(t, "123", l)
	require.Nil(t, err)

	l, err = readToNewLine(r)
	require.Equal(t, "456", l)
	require.Nil(t, err)

	l, err = readToNewLine(r)
	require.Equal(t, "789", l)
	require.Nil(t, err)

	l, err = readToNewLine(r)
	require.Equal(t, "", l)
	require.Nil(t, err)
}
