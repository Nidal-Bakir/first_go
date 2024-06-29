package math_solver

import (
	"context"
	"errors"
	"io"
)

var ErrNoMathToSolve = errors.New("error not math to solve")
var ErrSyntaxError = errors.New("error syntax error")

type MathSolver interface {
	Resolve(ctx context.Context, exp string) (float64, error)
}

type Processor struct {
	solver MathSolver
}

func NewProcessor(solver MathSolver) *Processor {
	return &Processor{solver: solver}
}

func (p Processor) Solve(ctx context.Context, r io.Reader) (float64, error) {
	exp, err := readToNewLine(r)
	if err != nil {
		return 0, err
	}

	if exp == "" {
		return 0, ErrNoMathToSolve
	}

	return p.solver.Resolve(ctx, exp)
}

func readToNewLine(r io.Reader) (string, error) {
	var out []byte
	buff := make([]byte, 1)

	for {
		n, err := r.Read(buff)

		if n != 0 {
			if buff[0] == '\n' {
				break
			}
			out = append(out, buff[0])
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	return string(out), nil
}
