package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

var constAlph, _ = equation.NewAlphabet([]string{"A", "B", "C"}, 3, 1)
var varsAlph, _ = equation.NewAlphabet([]string{"x", "y", "z"}, 3, 1)

func Test_Length_1(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("x A A y y = y A x", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_1 failed: error should be nil")
		return
	}
	check, _, _ := checkLengthRules(&eq)
	if check {
		t.Errorf("Test_Length_1 failed: result should be false")
	}
}

func Test_Length_2(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("x y = y A x x", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_2 failed: error should be nil")
		return
	}
	check, _, _ := checkLengthRules(&eq)
	if check {
		t.Errorf("Test_Length_2 failed: result should be false")
	}
}

func Test_Length_3(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("x y y y y = B y B x", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_3 failed: error should be nil")
		return
	}
	check, _, _ := checkLengthRules(&eq)
	if check {
		t.Errorf("Test_Length_3 failed: result should be false")
	}
}

func Test_Length_4(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("B B x = x y y y", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_4 failed: error should be nil")
		return
	}
	check, _, _ := checkLengthRules(&eq)
	if check {
		t.Errorf("Test_Length_4 failed: result should be false")
	}
}

func Test_Length_Symbol_1(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("x y x y = x B y x A C", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_Symbol_1 failed: error should be nil")
		return
	}
	check, sym, n := checkLengthRules(&eq)
	if !check {
		t.Errorf("Test_Length_Symbol_1 failed: result should be true")
		return
	}
	if sym != symbol.Var("y") {
		t.Errorf("Test_Length_Symbol_1 failed: symbol should be %s, but got: %s", "y", sym.Value())
		return
	}
	if n != 3 {
		t.Errorf("Test_Length_Symbol_1 failed: number of symbols should be %d, but got: %d", 3, n)
	}
}

func Test_Length_Symbol_2(t *testing.T) {
	var eq equation.Equation
	err := eq.Init("x A A y z y = z x y y x", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Length_Symbol_2 failed: error should be nil")
		return
	}
	check, sym, n := checkLengthRules(&eq)
	if !check {
		t.Errorf("Test_Length_Symbol_2 failed: result should be true")
		return
	}
	if sym != symbol.Var("x") {
		t.Errorf("Test_Length_Symbol_2 failed: symbol should be %s, but got: %s", "x", sym.Value())
		return
	}
	if n != 2 {
		t.Errorf("Test_Length_Symbol_2 failed: number of symbols should be %d, but got: %d", 2, n)
	}
}

func Test_Multiplicity_1(t *testing.T) {
	var solver Solver
	err := solver.Init("{A, B, C}", "{x, y, z}", "x A A A = A z B x", PrintOptions{}, SolveOptions{AlgorithmMode: "Finite"})
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Multiplicity_1 failed: error should be nil")
		return
	}
	subst := equation.NewSubstitution(solver.equation.RightPart.Symbols[1], []symbol.Symbol{solver.getLetter()})
	resultEq := solver.equation.Substitute(subst)

	if resultEq.LeftPart.Length != 4 || resultEq.RightPart.Length != 4 {
		t.Errorf("Test_Multiplicity_1 failed: wrong substitute result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Var("x") ||
			resultEq.LeftPart.Symbols[1] != symbol.Const("A") || resultEq.LeftPart.Symbols[2] != symbol.Const("A") ||
			resultEq.LeftPart.Symbols[3] != symbol.Const("A") ||
			resultEq.RightPart.Symbols[0] != symbol.Const("A") ||
			!symbol.IsLetter(resultEq.RightPart.Symbols[1]) || resultEq.RightPart.Symbols[2] != symbol.Const("B") ||
			resultEq.RightPart.Symbols[3] != symbol.Var("x") {
			t.Errorf("Test_Multiplicity_1 failed: wrong reduce result: ")
			resultEq.Print()
			return
		}
	}

	check := analiseMultiplicity(&solver.equation)
	if check {
		t.Errorf("Test_Multiplicity_1 failed: result should be false")
	}
}

func Test_Multiplicity_2(t *testing.T) {
	var solver Solver
	err := solver.Init("{A, B, C, D}", "{x, y, z}", "B A A x y = y A z x", PrintOptions{}, SolveOptions{AlgorithmMode: "Finite"})
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Multiplicity_2 failed: error should be nil")
		return
	}
	subst := equation.NewSubstitution(solver.equation.RightPart.Symbols[2], []symbol.Symbol{solver.getLetter()})
	resultEq := solver.equation.Substitute(subst)

	if resultEq.LeftPart.Length != 5 || resultEq.RightPart.Length != 4 {
		t.Errorf("Test_Multiplicity_2 failed: wrong substitute result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Const("B") || resultEq.LeftPart.Symbols[1] != symbol.Const("A") ||
			resultEq.LeftPart.Symbols[2] != symbol.Const("A") || resultEq.LeftPart.Symbols[3] != symbol.Var("x") ||
			resultEq.LeftPart.Symbols[4] != symbol.Var("y") ||
			resultEq.RightPart.Symbols[0] != symbol.Var("y") || resultEq.RightPart.Symbols[1] != symbol.Const("A") ||
			!symbol.IsLetter(resultEq.RightPart.Symbols[2]) || resultEq.RightPart.Symbols[3] != symbol.Var("x") {
			t.Errorf("Test_Multiplicity_2 failed: wrong reduce result: ")
			resultEq.Print()
			return
		}
	}

	check := analiseMultiplicity(&solver.equation)
	if check {
		t.Errorf("Test_Multiplicity_2 failed: result should be false")
	}
}

func Test_Multiplicity_3(t *testing.T) {
	var solver Solver
	err := solver.Init("{A, B, C, D}", "{x, y, z}", "y x z A z B B =  B C A y A x C", PrintOptions{}, SolveOptions{AlgorithmMode: "Finite"})
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_Multiplicity_3 failed: error should be nil")
		return
	}
	subst := equation.NewSubstitution(solver.equation.LeftPart.Symbols[2], []symbol.Symbol{solver.getLetter()})
	resultEq := solver.equation.Substitute(subst)

	if resultEq.LeftPart.Length != 7 || resultEq.RightPart.Length != 7 {
		t.Errorf("Test_Multiplicity_3 failed: wrong substitute result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Var("y") || resultEq.LeftPart.Symbols[1] != symbol.Var("x") ||
			!symbol.IsLetter(resultEq.LeftPart.Symbols[2]) || resultEq.LeftPart.Symbols[3] != symbol.Const("A") ||
			!symbol.IsLetter(resultEq.LeftPart.Symbols[4]) || resultEq.LeftPart.Symbols[5] != symbol.Const("B") ||
			resultEq.LeftPart.Symbols[6] != symbol.Const("B") ||
			resultEq.RightPart.Symbols[0] != symbol.Const("B") || resultEq.RightPart.Symbols[1] != symbol.Const("C") ||
			resultEq.RightPart.Symbols[2] != symbol.Const("A") || resultEq.RightPart.Symbols[3] != symbol.Var("y") ||
			resultEq.RightPart.Symbols[4] != symbol.Const("A") || resultEq.RightPart.Symbols[5] != symbol.Var("x") ||
			resultEq.RightPart.Symbols[6] != symbol.Const("C") {
			t.Errorf("Test_Multiplicity_3 failed: wrong substitute result: ")
			resultEq.Print()
			return
		}
	}

	check := analiseMultiplicity(&solver.equation)
	if check {
		t.Errorf("Test_Multiplicity_3 failed: result should be false")
	}
}
