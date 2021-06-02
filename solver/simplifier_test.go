package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

var firstTestEquation = equation.NewEquation([]symbol.Symbol{symbol.Const("A"),
	symbol.Const("B"), symbol.Const("A"), symbol.Var("x")}, []symbol.Symbol{symbol.Var("x"),
	symbol.Const("B"), symbol.Const("A"), symbol.Const("A")})

var firstTestResult = equation.NewSingleEquation(equation.NewEquation([]symbol.Symbol{symbol.Const("A"),
	symbol.Const("B"), symbol.Const("A"), symbol.Var("x")}, []symbol.Symbol{symbol.Var("x"),
	symbol.Const("B"), symbol.Const("A"), symbol.Const("A")}))

// A B A x = x B A A

func TestSimplifier_Simplify_First(t *testing.T) {
	var simplifier Simplifier
	err := simplifier.Init("{A, B}", "{x, y}", PrintOptions{}, SolveOptions{
		CycleRange:               10,
		FullGraph:                true,
		AlgorithmMode:            "Finite",
		SaveLettersSubstitutions: true,
		NeedsSimplification:      false,
	})
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_First error should be nil: %v", fmt.Sprintf("error initializing simplifier: %v \n", err))
		return
	}
	tree := NewTreeWEquation("0", firstTestEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_First error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(firstTestResult) {
		t.Errorf("TestSimplifier_Simplify_First error: result must be: %v", firstTestResult.String())
	}
}

var secondTestEquation = equation.NewEquation([]symbol.Symbol{symbol.Var("x"), symbol.Var("y"),
	symbol.Const("A"), symbol.Const("B")}, []symbol.Symbol{
	symbol.Const("A"), symbol.Const("B"), symbol.Var("x"), symbol.Var("y")})

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

// x y A B = A B x y

func TestSimplifier_Simplify_Second(t *testing.T) {
	var simplifier Simplifier
	err := simplifier.Init("{A, B}", "{x, y}", PrintOptions{}, SolveOptions{
		CycleRange:               10,
		FullGraph:                true,
		AlgorithmMode:            "Finite",
		SaveLettersSubstitutions: true,
		NeedsSimplification:      false,
	})
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Second error should be nil: %v", fmt.Sprintf("error initializing simplifier: %v \n", err))
		return
	}
	tree := NewTreeWEquation("0", secondTestEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Second error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(secondTestResult) {
		t.Errorf("TestSimplifier_Simplify_Second error: result must be: %v", secondTestResult.String())
	}
}

var thirdTestEquation = equation.NewEquation([]symbol.Symbol{
	symbol.Const("A"), symbol.Const("B"), symbol.Var("y")}, []symbol.Symbol{
	symbol.Var("y"), symbol.Const("A"), symbol.Const("B")})

var thirdTestResult = equation.NewSingleEquation(thirdTestEquation)

// A B y = y A B

func TestSimplifier_Simplify_Third(t *testing.T) {
	var simplifier Simplifier
	err := simplifier.Init("{A, B}", "{x, y}", PrintOptions{}, SolveOptions{
		CycleRange:               10,
		FullGraph:                true,
		AlgorithmMode:            "Finite",
		SaveLettersSubstitutions: true,
		NeedsSimplification:      false,
	})
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Third error should be nil: %v", fmt.Sprintf("error initializing simplifier: %v \n", err))
		return
	}
	tree := NewTreeWEquation("0", thirdTestEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Third error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(thirdTestResult) {
		t.Errorf("TestSimplifier_Simplify_Third error: result must be: %v", thirdTestResult.String())
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

// c A x y = x A y c  ==>  c y = y c & c A x = x A c

func TestSimplifier_Simplify_Symbols_1(t *testing.T) {
	var simplifier Simplifier
	err := simplifier.Init("{A, B}", "{x, y}", PrintOptions{}, SolveOptions{
		CycleRange:               10,
		FullGraph:                true,
		AlgorithmMode:            "Finite",
		SaveLettersSubstitutions: true,
		NeedsSimplification:      false,
	})
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error should be nil: %v", fmt.Sprintf("error initializing simplifier: %v \n", err))
		return
	}
	err = simplifier.solver.setLetterAlphabet("{c}")
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_1 error should be nil: %v", fmt.Sprintf("error setting letters alphabet: %v \n", err))
		return
	}
	tree := NewTreeWEquation("0", firstSymbolEquation)
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

var secondSymbolTestResult = equation.NewDisjunctionSystem([]equation.EquationsSystem{
	equation.NewConjunctionSystemFromEquations([]equation.Equation{
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c")},
			[]symbol.Symbol{symbol.Const("B")}),
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("r")},
			[]symbol.Symbol{symbol.Const("A")}),
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("t")},
			[]symbol.Symbol{symbol.Const("A")}),
		equation.NewEquation([]symbol.Symbol{symbol.Const("B"), symbol.Const("A"),
			symbol.Const("A"), symbol.Var("x")},
			[]symbol.Symbol{symbol.Var("x"), symbol.Const("B"), symbol.Const("A"),
				symbol.Const("A")}),
	}),

	equation.NewConjunctionSystemFromEquations([]equation.Equation{
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("c")},
			[]symbol.Symbol{symbol.LetterVar("t")}),
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("r")},
			[]symbol.Symbol{symbol.Const("B")}),
		equation.NewEquation([]symbol.Symbol{symbol.LetterVar("t"), symbol.Const("A"),
			symbol.Const("B"), symbol.Var("x")},
			[]symbol.Symbol{symbol.Var("x"), symbol.Const("B"), symbol.LetterVar("t"),
				symbol.Const("A")}),
	}),
})

// c A r x = x B t A  ==>  (c = B & t = A & r = A & B A A x = x B A A) V (c = t & r = B & t A B x = x B t A)

func TestSimplifier_Simplify_Symbols_2(t *testing.T) {
	var simplifier Simplifier
	err := simplifier.Init("{A, B}", "{x, y}", PrintOptions{}, SolveOptions{
		CycleRange:               10,
		FullGraph:                true,
		AlgorithmMode:            "Finite",
		SaveLettersSubstitutions: true,
		NeedsSimplification:      false,
	})
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error initializing simplifier: %v \n", err))
		return
	}
	err = simplifier.solver.setLetterAlphabet("{c, r, t}")
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error setting letters alphabet: %v \n", err))
		return
	}
	tree := NewTreeWEquation("0", secondSymbolEquation)
	err = simplifier.Simplify(&tree)
	if err != nil {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error should be nil: %v", fmt.Sprintf("error simplifing equation: %v \n", err))
		return
	}
	if !tree.simplified.Equals(secondSymbolTestResult) {
		t.Errorf("TestSimplifier_Simplify_Symbols_2 error: result must be: %v", secondSymbolTestResult.String())
	}
}
