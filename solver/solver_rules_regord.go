package solver

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

func checkFirstSymbols(eq *equation.Equation) bool {
	if eq.LeftPart.Length != 0 && eq.RightPart.Length != 0 {
		if symbol.IsLetterOrVar(eq.LeftPart.Symbols[0]) && symbol.IsLetterOrVar(eq.RightPart.Symbols[0]) {
			return eq.LeftPart.Symbols[0].Value() == eq.RightPart.Symbols[0].Value()
		}
		return true
	}
	return true
}
