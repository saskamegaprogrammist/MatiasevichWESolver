package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

var solveOptionsInvalid = SolveOptions{
	LengthAnalysis:             false,
	SplitByEquidecomposability: false,
	CycleRange:                 20,
	FullGraph:                  false,
	AlgorithmMode:              "Invalid",
	FullSystem:                 false,
}
var solveOptionsFinite = SolveOptions{
	LengthAnalysis:             false,
	SplitByEquidecomposability: false,
	CycleRange:                 20,
	FullGraph:                  false,
	AlgorithmMode:              "Finite",
	FullSystem:                 false,
}

var solveOptionsFiniteSplit = SolveOptions{
	LengthAnalysis:             false,
	SplitByEquidecomposability: true,
	CycleRange:                 20,
	FullGraph:                  false,
	AlgorithmMode:              "Finite",
	FullSystem:                 false,
}

var solveOptionsFiniteFullGraph = SolveOptions{
	LengthAnalysis:             false,
	SplitByEquidecomposability: false,
	CycleRange:                 20,
	FullGraph:                  true,
	AlgorithmMode:              "Finite",
	FullSystem:                 false,
}
var solveOptionsStandard = SolveOptions{
	LengthAnalysis:             false,
	SplitByEquidecomposability: false,
	CycleRange:                 20,
	FullGraph:                  false,
	AlgorithmMode:              "Standard",
	FullSystem:                 false,
}

var solveOptionsStandardWLength = SolveOptions{
	LengthAnalysis:             true,
	SplitByEquidecomposability: false,
	CycleRange:                 20,
	FullGraph:                  false,
	AlgorithmMode:              "Standard",
	FullSystem:                 false,
}
var printOptions = PrintOptions{
	Dot:       false,
	Png:       false,
	OutputDir: "../output_files",
}

var printOptionsPrint = PrintOptions{
	Dot:       true,
	Png:       true,
	OutputDir: "../output_files",
}

var test1InitErrorMessage = "error matching alphabet type: invalid algorithm type: Invalid"

func Test_Init_Error_1(t *testing.T) {
	var solver Solver
	err := solver.Init("", "", "", printOptions, solveOptionsInvalid)
	if err == nil {
		t.Errorf("Test_Init_Error_1 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test1InitErrorMessage {
			t.Errorf("Test_Init_Error_1 failed: wrong error message")
		}
	}
}

var test2InitErrorMessage = "error parsing constants: invalid constants alphabet: a,c"

func Test_Init_Error_2(t *testing.T) {
	var solver Solver
	err := solver.Init("a,c", "", "", printOptions, solveOptionsFinite)
	if err == nil {
		t.Errorf("Test_Init_Error_2: error shouldn\\'t be nil")
	} else {
		if err.Error() != test2InitErrorMessage {
			t.Errorf("Test_Init_Error_2 failed: wrong error message")
		}
	}
}

var test3InitErrorMessage = "error parsing constants: empty constant in alphabet: {a, c, , s}"

func Test_Init_Error_3(t *testing.T) {
	var solver Solver
	err := solver.Init("{a, c, , s}", "", "", printOptions, solveOptionsFinite)
	if err == nil {
		t.Errorf("Test_Init_Error_3: error shouldn\\'t be nil")
	} else {
		fmt.Println(err.Error())
		if err.Error() != test3InitErrorMessage {
			t.Errorf("Test_Init_Error_3 failed: wrong error message")
		}
	}
}

var test4InitErrorMessage = "error parsing vars: invalid constants alphabet: b"

func Test_Init_Error_4(t *testing.T) {
	var solver Solver
	err := solver.Init("{a, c, s}", "b", "", printOptions, solveOptionsFinite)
	if err == nil {
		t.Errorf("Test_Init_Error_4: error shouldn\\'t be nil")
	} else {
		if err.Error() != test4InitErrorMessage {
			t.Errorf("Test_Init_Error_4 failed: wrong error message")
		}
	}
}

var test5InitErrorMessage = "error parsing vars: empty constant in alphabet: {a, , s}"

func Test_Init_Error_5(t *testing.T) {
	var solver Solver
	err := solver.Init("{b, n}", "{a, , s}", "", printOptions, solveOptionsFinite)
	if err == nil {
		t.Errorf("Test_Init_Error_5: error shouldn\\'t be nil")
	} else {
		if err.Error() != test5InitErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_Init_Error_5 failed: wrong error message")
		}
	}
}

var test6InitErrorMessage = "error parsing equation: invalid equation: ab"

func Test_Init_Error_6(t *testing.T) {
	var solver Solver
	err := solver.Init("{b, n}", "{a, s}", "ab", printOptions, solveOptionsStandard)
	if err == nil {
		t.Errorf("Test_Init_Error_6 error shouldn\\'t be nil")
	} else {
		if err.Error() != test6InitErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_Init_Error_6 failed: wrong error message")
		}
	}
}

var test7InitErrorMessage = "error parsing constants: letters must be separated by space: {b,n}"

func Test_Init_Error_7(t *testing.T) {
	var solver Solver
	err := solver.Init("{b,n}", "{a, s}", "ab", printOptions, solveOptionsStandard)
	if err == nil {
		t.Errorf("Test_Init_Error_7 error shouldn\\'t be nil")
	} else {
		if err.Error() != test7InitErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_Init_Error_7 failed: wrong error message")
		}
	}
}

var trueStr = "TRUE"
var falseStr = "FALSE"
var cycledStr = "CYCLED"

func Test_Solve_1(t *testing.T) {
	var solver Solver
	err := solver.Init("{a}", "{u, v}", "u a v = v a u", printOptions, solveOptionsStandard)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_1 error should be nil")
	} else {
		result, _, _ := solver.Solve()
		if result != trueStr {
			t.Errorf("Test_Solve_1 result should be: %s, but got: %s", trueStr, result)
		}
	}
}

func Test_Solve_2(t *testing.T) {
	var solver Solver
	err := solver.Init("{a, b}", "{u}", "u u a = b u u", printOptions, solveOptionsStandard)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_2 error should be nil")
	} else {
		result, _, _ := solver.Solve()
		if result != cycledStr {
			t.Errorf("Test_Solve_2 result should be: %s, but got: %s", cycledStr, result)
		}
	}
}

func Test_Solve_3(t *testing.T) {
	var solver Solver
	err := solver.Init("{}", "{u, v, z}", "u u v v = z z", printOptions, solveOptionsStandard)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_3 error should be nil")
	} else {
		result, _, _ := solver.Solve()
		if result != trueStr {
			t.Errorf("Test_Solve_3 result should be: %s, but got: %s", trueStr, result)
		}
	}
}

func Test_Solve_4(t *testing.T) {
	var solver Solver
	err := solver.Init("{a}", "{u}", "a u = u", printOptions, solveOptionsStandard)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_4 error should be nil")
	} else {
		result, _, _ := solver.Solve()
		if result != falseStr {
			t.Errorf("Test_Solve_4 result should be: %s, but got: %s", falseStr, result)
		}
	}
}

func Test_Solve_5(t *testing.T) {
	var solver Solver
	err := solver.Init("{A, B}", "{x, y}", "x A A B y = y", printOptions, solveOptionsStandardWLength)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_5 error should be nil")
	} else {
		result, _, _ := solver.Solve()
		if result != falseStr {
			t.Errorf("Test_Solve_5 result should be: %s, but got: %s", falseStr, result)
		}
	}
}

// c d f e d b c f g a f b c = pl sy ie dh uf iu jm lf kw ic ex ih jb

// c d f = pl sy ie

var eq = equation.NewEquation([]symbol.Symbol{symbol.Const("c"), symbol.Const("d"), symbol.Const("f")},
	[]symbol.Symbol{symbol.LetterVar("c"), symbol.LetterVar("sy"), symbol.LetterVar("ie")})

func Test_Solve_6(t *testing.T) {
	var solver Solver
	err := solver.InitWoEquation("{c, e, f, d, b, g, a}", "{x, y}", printOptions, solveOptionsFiniteSplit)
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_6 error should be nil")
		return
	}
	err = solver.SetEquation(eq)
	if err != nil {
		fmt.Printf("error setting equation: %v \n", err)
		t.Errorf("Test_Solve_6 error should be nil")
		return
	}
	err = solver.setLetterAlphabet("{pl, sy, ie, dh, uf, iu, jm, lf, kw, ic, ex, ih, jb}")
	if err != nil {
		fmt.Printf("error setting letters alphabet: %v \n", err)
		t.Errorf("Test_Solve_6 error should be nil")
		return
	}
	result, _, _ := solver.Solve()
	if result != trueStr {
		t.Errorf("Test_Solve_6 result should be: %s, but got: %s", trueStr, result)
	}
}
