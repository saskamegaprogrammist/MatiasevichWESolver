package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

var firstTestResult = equation.NewSingleEquation(equation.NewEquation([]symbol.Symbol{symbol.Const("A"),
	symbol.Const("B"), symbol.Const("A"), symbol.Var("x")}, []symbol.Symbol{symbol.Var("x"),
	symbol.Const("B"), symbol.Const("A"), symbol.Const("A")}))

func TestSimplifier_Simplify_First(t *testing.T) {
	var simplifier Simplifier
	var solver Solver
	err := solver.Init("{A, B}", "{x, y}", "A B A x = x B A A", printOptions, solveOptionsFiniteFullGraph)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_First error should be nil: %v", fmt.Sprintf("error initializing solver: %v \n", err))
		return
	}
	simplifier.InitWithSolver(&solver)
	tree := NewTree("0", solver.equation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_First error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(firstTestResult) {
		t.Errorf("TestSimplifier_Simplify_First error: result must be: %v", firstTestResult.String())
	}
}

var secondTestResult = equation.NewDisjunctionSystem([]equation.EquationsSystem{equation.NewConjunctionSystemFromEquations([]equation.Equation{
	equation.NewEquation([]symbol.Symbol{symbol.Const("A"), symbol.Const("B"), symbol.Var("x")},
		[]symbol.Symbol{symbol.Var("x"), symbol.Const("A"), symbol.Const("B")}),
	equation.NewEquation([]symbol.Symbol{symbol.Const("A"), symbol.Const("B"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"), symbol.Const("B")}),
}), equation.NewConjunctionSystemFromEquations([]equation.Equation{
	equation.NewEquation([]symbol.Symbol{symbol.Const("A"), symbol.Const("B"), symbol.Var("x")},
		[]symbol.Symbol{symbol.Var("x"), symbol.Const("B"), symbol.Const("A")}),
	equation.NewEquation([]symbol.Symbol{symbol.Const("B"), symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"), symbol.Const("B")}),
})})

func TestSimplifier_Simplify_Second(t *testing.T) {
	var simplifier Simplifier
	var solver Solver
	err := solver.Init("{A, B}", "{x, y}", "x y A B = A B x y", printOptions, solveOptionsFiniteFullGraph)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Second error should be nil: %v", fmt.Sprintf("error initializing solver: %v \n", err))
		return
	}
	simplifier.InitWithSolver(&solver)
	tree := NewTree("0", solver.equation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Second error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(secondTestResult) {
		t.Errorf("TestSimplifier_Simplify_Second error: result must be: %v", secondTestResult.String())
	}
}

var firstSymbolEquation = equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"),
	symbol.Const("A"), symbol.Var("x"), symbol.Var("y")}, []symbol.Symbol{symbol.Var("x"),
	symbol.Const("A"), symbol.Var("y"), symbol.LetterVar("c")})

var firstSymbolTestResult = equation.NewConjunctionSystemFromEquations([]equation.Equation{
	equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"), symbol.Const("A"), symbol.Var("x")},
		[]symbol.Symbol{symbol.Var("x"), symbol.Const("A"), symbol.LetterVar("c")}),
	equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.LetterVar("c")}),
})

func TestSimplifier_Simplify_Symbols_1(t *testing.T) {
	var simplifier Simplifier
	var solver Solver
	err := solver.InitWoEquation("{A, B}", "{x, y}", printOptions, solveOptionsFiniteFullGraph)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error should be nil: %v", fmt.Sprintf("error initializing solver: %v \n", err))
		return
	}
	err = solver.setLetterAlphabet("{c}")
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error should be nil: %v", fmt.Sprintf("error setting letters alphabet: %v \n", err))
		return
	}
	simplifier.InitWithSolver(&solver)
	tree := NewTree("0", firstSymbolEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(firstSymbolTestResult) {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error: result must be: %v", firstSymbolTestResult.String())
	}
}

var secondSymbolEquation = equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"),
	symbol.Const("A"), symbol.LetterVar("r"), symbol.Var("x")}, []symbol.Symbol{symbol.Var("x"),
	symbol.Const("B"), symbol.LetterVar("t"), symbol.Const("A")})

//
//var secondSymbolTestResult = equation.NewConjunctionSystemFromEquations([]equation.Equation{
//	equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"), symbol.Const("A"), symbol.Var("x")},
//		[]symbol.Symbol{symbol.Var("x"), symbol.Const("A"), symbol.LetterVar("c")}),
//	equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c"), symbol.Var("y")},
//		[]symbol.Symbol{symbol.Var("y"), symbol.LetterVar("c")}),
//})

func TestSimplifier_Simplify_Symbols_2(t *testing.T) {
	var simplifier Simplifier
	var solver Solver
	err := solver.InitWoEquation("{A, B}", "{x, y}", printOptions, solveOptionsFiniteFullGraph)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error initializing solver: %v \n", err))
		return
	}
	err = solver.setLetterAlphabet("{c, r, t}")
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error setting letters alphabet: %v \n", err))
		return
	}
	simplifier.InitWithSolver(&solver)
	tree := NewTree("0", secondSymbolEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	//if !tree.simplified.Equals(secondSymbolTestResult) {
	//	t.Errorf("TestSimplifier_Simplify_Symbols_2 error: result must be: %v", secondSymbolTestResult.String())
	//}
}
