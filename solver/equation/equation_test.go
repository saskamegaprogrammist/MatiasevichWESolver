package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

var constAlph = Alphabet{
	words:         []string{"a", "b", "c"},
	size:          3,
	maxWordLength: 1,
}
var varsAlph = Alphabet{
	words:         []string{"ux", "yi", "p"},
	size:          3,
	maxWordLength: 2,
}

var test1InitEqErrorMessage = "invalid equation: aa"

func Test_InitEq_Error_1(t *testing.T) {
	var eq Equation
	err := eq.Init("aa", &constAlph, &varsAlph)
	if err == nil {
		t.Errorf("Test_InitEq_Error_1 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test1InitEqErrorMessage {
			t.Errorf("Test_InitEq_Error_1 failed: wrong error message")
		}
	}
}

var test2InitEqErrorMessage = "error matching alphabet: error matching word: no match found with word: o"

func Test_InitEq_Error_2(t *testing.T) {
	var eq Equation
	err := eq.Init("a = o", &constAlph, &varsAlph)
	if err == nil {
		t.Errorf("Test_InitEq_Error_2 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test2InitEqErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_InitEq_Error_2 failed: wrong error message")
		}
	}
}

var test3InitEqErrorMessage = "invalid equation: a=o"

func Test_InitEq_Error_3(t *testing.T) {
	var eq Equation
	err := eq.Init("a=o", &constAlph, &varsAlph)
	if err == nil {
		t.Errorf("Test_InitEq_Error_3 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test3InitEqErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_InitEq_Error_3 failed: wrong error message")
		}
	}
}

var constAlphMisleading = Alphabet{
	words:         []string{"aaaa", "aab", "aabx"},
	size:          3,
	maxWordLength: 4,
}
var varsAlphMisleading = Alphabet{
	words:         []string{"aaaaa", "o", "bxp"},
	size:          3,
	maxWordLength: 5,
}

func Test_InitEq_1(t *testing.T) {
	var eq Equation
	err := eq.Init("aaaaa aabx = aaaa bxp", &constAlphMisleading, &varsAlphMisleading)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_InitEq_1 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Var("aaaaa") || eq.LeftPart.Symbols[1] != symbol.Const("aabx") ||
		eq.RightPart.Symbols[0] != symbol.Const("aaaa") || eq.RightPart.Symbols[1] != symbol.Var("bxp") {
		t.Errorf("Test_InitEq_1 failed: wrong eq parsing: ")
		eq.Print()
	}
}

func Test_InitEq_2(t *testing.T) {
	var eq Equation
	err := eq.Init("a a ux = yi b p", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_InitEq_2 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Const("a") ||
		eq.LeftPart.Symbols[2] != symbol.Var("ux") ||
		eq.RightPart.Symbols[0] != symbol.Var("yi") || eq.RightPart.Symbols[1] != symbol.Const("b") ||
		eq.RightPart.Symbols[2] != symbol.Var("p") {
		t.Errorf("Test_InitEq_2 failed: wrong eq parsing: ")
		eq.Print()
	}
}

var constAlphNew = Alphabet{
	words:         []string{"a", "b", "c", "d"},
	size:          4,
	maxWordLength: 1,
}
var varsAlphNew = Alphabet{
	words:         []string{"u", "v", "x", "y"},
	size:          3,
	maxWordLength: 1,
}

func TestEquation_Reduce_1(t *testing.T) {
	var eq Equation
	err := eq.Init("a b x = $ a b v", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Reduce_1 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Const("b") ||
		eq.LeftPart.Symbols[2] != symbol.Var("x") ||
		eq.RightPart.Symbols[0] != symbol.Empty() || eq.RightPart.Symbols[1] != symbol.Const("a") ||
		eq.RightPart.Symbols[2] != symbol.Const("b") || eq.RightPart.Symbols[3] != symbol.Var("v") {
		t.Errorf("TestEquation_Reduce_1 failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	eq.Reduce()
	if eq.LeftPart.Length == 0 || eq.RightPart.Length == 0 {
		t.Errorf("TestEquation_Reduce_1 failed: wrong reduce result: len shouldn\\'t be nil")
	} else {
		if eq.LeftPart.Symbols[0] != symbol.Var("x") || eq.RightPart.Symbols[0] != symbol.Var("v") {
			t.Errorf("TestEquation_Reduce_1 failed: wrong reduce result: ")
			eq.Print()
		}
	}
}

func TestEquation_Reduce_2(t *testing.T) {
	var eq Equation
	err := eq.Init("x $ = x $ a v", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Reduce_2 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Var("x") || eq.LeftPart.Symbols[1] != symbol.Empty() ||
		eq.RightPart.Symbols[0] != symbol.Var("x") || eq.RightPart.Symbols[2] != symbol.Const("a") ||
		eq.RightPart.Symbols[1] != symbol.Empty() || eq.RightPart.Symbols[3] != symbol.Var("v") {
		t.Errorf("TestEquation_Reduce_2 failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	eq.Reduce()
	if eq.LeftPart.Length != 0 {
		t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result: left len should be nil")
		return
	}
	if eq.RightPart.Length != 2 {
		t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result: right len should be: %d", 2)
	} else {
		if eq.RightPart.Symbols[0] != symbol.Const("a") || eq.RightPart.Symbols[1] != symbol.Var("v") {
			t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result: ")
			eq.Print()
		}
	}
}

func TestEquation_Substitute_1(t *testing.T) {
	var eq Equation
	err := eq.Init("a b x = v b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_1 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Const("b") ||
		eq.LeftPart.Symbols[2] != symbol.Var("x") ||
		eq.RightPart.Symbols[0] != symbol.Var("v") || eq.RightPart.Symbols[1] != symbol.Const("b") {
		t.Errorf("TestEquation_Substitute_1 failed: wrong eq parsing: ")
		eq.Print()
		return
	}

	subst := NewSubstitution(eq.RightPart.Symbols[0], []symbol.Symbol{eq.LeftPart.Symbols[0], eq.RightPart.Symbols[0]})
	resultEq := eq.Substitute(subst)

	if resultEq.LeftPart.Length != 2 || resultEq.RightPart.Length != 2 {
		t.Errorf("TestEquation_Substitute_1 failed: wrong reduce result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Const("b") ||
			resultEq.LeftPart.Symbols[1] != symbol.Var("x") ||
			resultEq.RightPart.Symbols[0] != symbol.Var("v") || resultEq.RightPart.Symbols[1] != symbol.Const("b") {
			t.Errorf("TestEquation_Substitute_1 failed: wrong reduce result: ")
			resultEq.Print()
		}
	}
}

func TestEquation_Substitute_2(t *testing.T) {
	var eq Equation
	err := eq.Init("x b x a = v b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_2 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[3] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Const("b") ||
		eq.LeftPart.Symbols[2] != symbol.Var("x") || eq.LeftPart.Symbols[0] != symbol.Var("x") ||
		eq.RightPart.Symbols[0] != symbol.Var("v") || eq.RightPart.Symbols[1] != symbol.Const("b") {
		t.Errorf("TestEquation_Substitute_2 failed: wrong eq parsing: ")
		eq.Print()
		return
	}

	subst := NewSubstitution(eq.LeftPart.Symbols[0], []symbol.Symbol{eq.RightPart.Symbols[0], eq.LeftPart.Symbols[0]})
	resultEq := eq.Substitute(subst)
	if resultEq.LeftPart.Length != 5 || resultEq.RightPart.Length != 1 {
		t.Errorf("TestEquation_Substitute_2 failed: wrong reduce result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Var("x") ||
			resultEq.LeftPart.Symbols[1] != symbol.Const("b") ||
			resultEq.LeftPart.Symbols[2] != symbol.Var("v") || resultEq.LeftPart.Symbols[3] != symbol.Var("x") ||
			resultEq.LeftPart.Symbols[4] != symbol.Const("a") ||
			resultEq.RightPart.Symbols[0] != symbol.Const("b") {
			t.Errorf("TestEquation_Substitute_2 failed: wrong reduce result: ")
			resultEq.Print()
		}
	}
}

func TestEquation_Substitute_3(t *testing.T) {
	var eq Equation
	err := eq.Init("x a v = ", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_3 failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[1] != symbol.Const("a") ||
		eq.LeftPart.Symbols[2] != symbol.Var("v") || eq.LeftPart.Symbols[0] != symbol.Var("x") ||
		eq.RightPart.Symbols[0] != symbol.Empty() {
		t.Errorf("TestEquation_Substitute_3 failed: wrong eq parsing: ")
		eq.Print()
		return
	}

	resultEq, _ := eq.SubstituteVarsWithEmpty()
	if resultEq.LeftPart.Length != 1 || resultEq.RightPart.Length != 1 {
		t.Errorf("TestEquation_Substitute_3 failed: wrong reduce result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Const("a") ||
			resultEq.RightPart.Symbols[0] != symbol.Empty() {
			t.Errorf("TestEquation_Substitute_3 failed: wrong reduce result: ")
			resultEq.Print()
		}
	}
}

func TestEquation_SymbolMap_1(t *testing.T) {
	var eq Equation
	err := eq.Init("a b b a x = u v a u x", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SymbolMap_1 failed: error should be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Const("b") ||
		eq.LeftPart.Symbols[2] != symbol.Const("b") || eq.LeftPart.Symbols[3] != symbol.Const("a") ||
		eq.LeftPart.Symbols[4] != symbol.Var("x") ||
		eq.RightPart.Symbols[0] != symbol.Var("u") || eq.RightPart.Symbols[1] != symbol.Var("v") ||
		eq.RightPart.Symbols[2] != symbol.Const("a") || eq.RightPart.Symbols[3] != symbol.Var("u") ||
		eq.LeftPart.Symbols[4] != symbol.Var("x") {
		t.Errorf("TestEquation_SymbolMap_1 failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	if eq.LeftPart.Structure.Consts()[symbol.Const("a")] != 2 || eq.LeftPart.Structure.Consts()[symbol.Const("b")] != 2 ||
		eq.LeftPart.Structure.Vars()[symbol.Var("x")] != 1 || eq.RightPart.Structure.Vars()[symbol.Var("u")] != 2 ||
		eq.RightPart.Structure.Vars()[symbol.Var("x")] != 1 || eq.RightPart.Structure.Vars()[symbol.Var("v")] != 1 ||
		eq.RightPart.Structure.Consts()[symbol.Const("a")] != 1 {
		t.Errorf("TestEquation_SymbolMap_1 failed: wrong map creation: ")
		eq.LeftPart.Structure.Print()
		fmt.Println()
		eq.RightPart.Structure.Print()
		return
	}
}

func TestEquation_SymbolMap_2(t *testing.T) {
	var eq Equation
	// u a x u b = v b a u => u a x v u b = b a v u
	err := eq.Init("u a x u b = v b a u", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SymbolMap_2 failed: error should be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Var("u") || eq.LeftPart.Symbols[1] != symbol.Const("a") ||
		eq.LeftPart.Symbols[2] != symbol.Var("x") || eq.LeftPart.Symbols[3] != symbol.Var("u") ||
		eq.LeftPart.Symbols[4] != symbol.Const("b") ||
		eq.RightPart.Symbols[0] != symbol.Var("v") || eq.RightPart.Symbols[1] != symbol.Const("b") ||
		eq.RightPart.Symbols[2] != symbol.Const("a") || eq.RightPart.Symbols[3] != symbol.Var("u") {
		t.Errorf("TestEquation_SymbolMap_2 failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	if eq.LeftPart.Structure.Consts()[symbol.Const("a")] != 1 || eq.LeftPart.Structure.Consts()[symbol.Const("b")] != 1 ||
		eq.LeftPart.Structure.Vars()[symbol.Var("x")] != 1 || eq.LeftPart.Structure.Vars()[symbol.Var("u")] != 2 ||
		eq.RightPart.Structure.Vars()[symbol.Var("u")] != 1 || eq.RightPart.Structure.Vars()[symbol.Var("v")] != 1 ||
		eq.RightPart.Structure.Consts()[symbol.Const("a")] != 1 || eq.RightPart.Structure.Consts()[symbol.Const("b")] != 1 {
		t.Errorf("TestEquation_SymbolMap_2 failed: wrong map creation: ")
		eq.LeftPart.Structure.Print()
		fmt.Println()
		eq.RightPart.Structure.Print()
		return
	}
	subst := NewSubstitution(eq.LeftPart.Symbols[0], []symbol.Symbol{eq.RightPart.Symbols[0], eq.LeftPart.Symbols[0]})
	resultEq := eq.Substitute(subst)

	if resultEq.LeftPart.Length != 6 || resultEq.RightPart.Length != 4 {
		t.Errorf("TestEquation_SymbolMap_2 failed: wrong reduce result: ")
		resultEq.Print()
	} else {
		if resultEq.LeftPart.Symbols[0] != symbol.Var("u") || resultEq.LeftPart.Symbols[1] != symbol.Const("a") ||
			resultEq.LeftPart.Symbols[2] != symbol.Var("x") || resultEq.LeftPart.Symbols[4] != symbol.Var("u") ||
			resultEq.LeftPart.Symbols[5] != symbol.Const("b") || resultEq.LeftPart.Symbols[3] != symbol.Var("v") ||
			resultEq.RightPart.Symbols[2] != symbol.Var("v") || resultEq.RightPart.Symbols[0] != symbol.Const("b") ||
			resultEq.RightPart.Symbols[1] != symbol.Const("a") || resultEq.RightPart.Symbols[3] != symbol.Var("u") {
			t.Errorf("TestEquation_SymbolMap_2 failed: wrong substitute result: ")
			resultEq.Print()
			return
		}

		if resultEq.LeftPart.Structure.Consts()[symbol.Const("a")] != 1 || resultEq.LeftPart.Structure.Consts()[symbol.Const("b")] != 1 ||
			resultEq.LeftPart.Structure.Vars()[symbol.Var("x")] != 1 || resultEq.LeftPart.Structure.Vars()[symbol.Var("u")] != 2 ||
			resultEq.LeftPart.Structure.Vars()[symbol.Var("v")] != 1 ||
			resultEq.RightPart.Structure.Vars()[symbol.Var("u")] != 1 || resultEq.RightPart.Structure.Vars()[symbol.Var("v")] != 1 ||
			resultEq.RightPart.Structure.Consts()[symbol.Const("a")] != 1 || resultEq.RightPart.Structure.Consts()[symbol.Const("b")] != 1 {
			t.Errorf("TestEquation_SymbolMap_2 failed: wrong map recreation: ")
			eq.LeftPart.Structure.Print()
			fmt.Println()
			eq.RightPart.Structure.Print()
			return
		}
	}
}

func TestEquation_SplitByEquidecomposability_Simple(t *testing.T) {
	var eq Equation
	// a u v = u b x => v = x ; a u = u b
	err := eq.Init("a u v = u b x", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Simple failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Const("a") || eq.LeftPart.Symbols[1] != symbol.Var("u") ||
		eq.LeftPart.Symbols[2] != symbol.Var("v") ||
		eq.RightPart.Symbols[0] != symbol.Var("u") || eq.RightPart.Symbols[1] != symbol.Const("b") ||
		eq.RightPart.Symbols[2] != symbol.Var("x") {
		t.Errorf("TestEquation_SplitByEquidecomposability_Simple failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	system := eq.SplitByEquidecomposability()
	if system.Size != 2 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Simple failed: wrong split result: ")
		system.PrintInfo()
	} else {
		if system.Equations[1].LeftPart.Symbols[0] != symbol.Const("a") ||
			system.Equations[1].LeftPart.Symbols[1] != symbol.Var("u") ||
			system.Equations[1].RightPart.Symbols[0] != symbol.Var("u") ||
			system.Equations[1].RightPart.Symbols[1] != symbol.Const("b") ||
			system.Equations[0].LeftPart.Symbols[0] != symbol.Var("v") ||
			system.Equations[0].RightPart.Symbols[0] != symbol.Var("x") {
			t.Errorf("TestEquation_SplitByEquidecomposability_Simple failed: wrong split result: ")
			system.PrintInfo()
		}
	}
}

func TestEquation_SplitByEquidecomposability_Backwards(t *testing.T) {
	var eq Equation
	// x c a u = y a u b => x c = y a; a u = u b
	err := eq.Init("x c a u = y a u b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Backwards failed: error shouldn be nil")
		return
	}
	if eq.LeftPart.Symbols[2] != symbol.Const("a") || eq.LeftPart.Symbols[3] != symbol.Var("u") ||
		eq.LeftPart.Symbols[0] != symbol.Var("x") || eq.LeftPart.Symbols[1] != symbol.Const("c") ||
		eq.RightPart.Symbols[2] != symbol.Var("u") || eq.RightPart.Symbols[3] != symbol.Const("b") ||
		eq.RightPart.Symbols[0] != symbol.Var("y") || eq.RightPart.Symbols[1] != symbol.Const("a") {
		t.Errorf("TestEquation_SplitByEquidecomposability_Backwards failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	system := eq.SplitByEquidecomposability()
	system.Print()
	if system.Size != 2 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Backwards failed: wrong split result: ")
		system.PrintInfo()
	} else {
		if system.Equations[1].LeftPart.Symbols[0] != symbol.Const("a") ||
			system.Equations[1].LeftPart.Symbols[1] != symbol.Var("u") ||
			system.Equations[1].RightPart.Symbols[0] != symbol.Var("u") ||
			system.Equations[1].RightPart.Symbols[1] != symbol.Const("b") ||
			system.Equations[0].LeftPart.Symbols[0] != symbol.Var("x") ||
			system.Equations[0].RightPart.Symbols[0] != symbol.Var("y") ||
			system.Equations[0].LeftPart.Symbols[1] != symbol.Const("c") ||
			system.Equations[0].RightPart.Symbols[1] != symbol.Const("a") {
			t.Errorf("TestEquation_SplitByEquidecomposability_Backwards failed: wrong split result: ")
			system.PrintInfo()
		}
	}
}

func TestEquation_SplitByEquidecomposability_Long(t *testing.T) {
	var eq Equation
	// x a y c y b = a x y c b y => y b = b y ; x a = a x
	err := eq.Init("x a y c y b = a x y c b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Long failed: error should be nil")
		return
	}
	if eq.LeftPart.Symbols[0] != symbol.Var("x") || eq.LeftPart.Symbols[1] != symbol.Const("a") ||
		eq.LeftPart.Symbols[2] != symbol.Var("y") || eq.LeftPart.Symbols[3] != symbol.Const("c") ||
		eq.LeftPart.Symbols[4] != symbol.Var("y") || eq.LeftPart.Symbols[5] != symbol.Const("b") ||
		eq.RightPart.Symbols[0] != symbol.Const("a") || eq.RightPart.Symbols[1] != symbol.Var("x") ||
		eq.RightPart.Symbols[2] != symbol.Var("y") || eq.RightPart.Symbols[3] != symbol.Const("c") ||
		eq.RightPart.Symbols[5] != symbol.Var("y") || eq.RightPart.Symbols[4] != symbol.Const("b") {
		t.Errorf("TestEquation_SplitByEquidecomposability_Long failed: wrong eq parsing: ")
		eq.Print()
		return
	}
	system := eq.SplitByEquidecomposability()
	if system.Size != 2 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Long failed: wrong split result: ")
		system.PrintInfo()
	} else {
		if system.Equations[0].LeftPart.Symbols[0] != symbol.Var("y") ||
			system.Equations[0].LeftPart.Symbols[1] != symbol.Const("b") ||
			system.Equations[0].RightPart.Symbols[0] != symbol.Const("b") ||
			system.Equations[0].RightPart.Symbols[1] != symbol.Var("y") ||
			system.Equations[1].LeftPart.Symbols[0] != symbol.Var("x") ||
			system.Equations[1].LeftPart.Symbols[1] != symbol.Const("a") ||
			system.Equations[1].RightPart.Symbols[0] != symbol.Const("a") ||
			system.Equations[1].RightPart.Symbols[1] != symbol.Var("x") {
			t.Errorf("TestEquation_SplitByEquidecomposability_Long failed: wrong split result: ")
			system.PrintInfo()
		}
	}
}

func TestEquation_SplitByEquidecomposability_Consts(t *testing.T) {
	var eq1, eq2, eq3 Equation
	err := eq1.Init("a b = b a", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: error should be nil")
		return
	}
	err = eq2.Init("a = a", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: error should be nil")
		return
	}
	err = eq3.Init("c = a", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: error should be nil")
		return
	}
	system1 := eq1.SplitByEquidecomposability()
	if system1.Size != 1 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: wrong split result: ")
		system1.PrintInfo()
	} else {
		if system1.Equations[0].isEquidecomposable {
			t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: should not be eqidecomposable: ")
			system1.PrintInfo()
		}
	}
	system2 := eq2.SplitByEquidecomposability()
	if system2.Size != 1 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: wrong split result: ")
		system2.PrintInfo()
	} else {
		if !system2.Equations[0].isEquidecomposable {
			t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: should be eqidecomposable: ")
			system2.PrintInfo()
		}
	}
	system3 := eq3.SplitByEquidecomposability()
	if system3.Size != 1 {
		t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: wrong split result: ")
		system3.PrintInfo()
	} else {
		if system3.Equations[0].isEquidecomposable {
			t.Errorf("TestEquation_SplitByEquidecomposability_Consts failed: should not be eqidecomposable: ")
			system3.PrintInfo()
		}
	}
}

func TestEquation_CheckSameness(t *testing.T) {
	var err error
	var eq1, eq2, eq3, eq4 Equation
	err = eq1.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_CheckSameness failed: error should be nil")
		return
	}
	err = eq2.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_CheckSameness failed: error should be nil")
		return
	}

	err = eq3.Init("x $ a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_CheckSameness failed: error should be nil")
		return
	}

	err = eq4.Init("x = b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_CheckSameness failed: error should be nil")
		return
	}

	var same bool
	same = eq1.CheckSameness(&eq2)
	if !same {
		t.Errorf("TestEquation_CheckSameness failed: eq1 and eq2 should be the same")
		return
	}

	same = eq1.CheckSameness(&eq3)
	if !same {
		t.Errorf("TestEquation_CheckSameness failed: eq1 and eq3 should be the same")
		return
	}

	same = eq1.CheckSameness(&eq4)
	if same {
		t.Errorf("TestEquation_CheckSameness failed: eq1 and eq4 should not be the same")
		return
	}
}

func TestEquation_CheckSimpleWord(t *testing.T) {
	var simple bool
	var symbols = []symbol.Symbol{symbol.Const("a")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: a should be simple")
		return
	}
	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("b")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: ab should be simple")
		return
	}
	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("b"), symbol.Const("a"),
		symbol.Const("b"), symbol.Const("c"), symbol.Const("x"), symbol.Const("y")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: ababcxy should be simple")
		return
	}
	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("a"), symbol.Const("a"),
		symbol.Const("a"), symbol.Const("a"), symbol.Const("a")}
	simple = checkSimpleWord(symbols)
	if simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: aaaaaa should not be simple")
		return
	}
	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("a"), symbol.Const("b"),
		symbol.Const("a"), symbol.Const("a"), symbol.Const("a")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: aabaaa should be simple")
		return
	}

	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("a"), symbol.Const("b"),
		symbol.Const("a"), symbol.Const("a"), symbol.Const("a"), symbol.Const("b")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: aabaaab should be simple")
		return
	}

	symbols = []symbol.Symbol{symbol.Const("a"), symbol.Const("a"), symbol.Const("b"),
		symbol.Const("a"), symbol.Const("a"), symbol.Const("b"), symbol.Const("b")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: aabaabb should be simple")
		return
	}

	symbols = []symbol.Symbol{symbol.Const("A"), symbol.Const("C"), symbol.Const("A"),
		symbol.Const("C"), symbol.Const("A"), symbol.Const("A"), symbol.Const("C"),
		symbol.Const("A"), symbol.Const("C"), symbol.Const("A")}
	simple = checkSimpleWord(symbols)
	if simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: ACACAACACA should not be simple")
		return
	}

	symbols = []symbol.Symbol{symbol.Const("A"), symbol.Const("C"), symbol.Const("C"),
		symbol.Const("C"), symbol.Const("A"), symbol.Const("C"), symbol.Const("C"),
		symbol.Const("C")}
	simple = checkSimpleWord(symbols)
	if simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: ACCCACCC should not be simple")
		return
	}
	symbols = []symbol.Symbol{symbol.Const("A"), symbol.Const("C"), symbol.Const("C"),
		symbol.Const("C"), symbol.Const("A"), symbol.Const("C"), symbol.Const("C"),
		symbol.Const("C"), symbol.Const("C")}
	simple = checkSimpleWord(symbols)
	if !simple {
		t.Errorf("TestEquation_CheckSimpleWord failed: ACCCACCCC should be simple")
		return
	}
}

func TestEquation_IsRegularlyOrdered(t *testing.T) {
	// A x T y = x y B A
	eq := NewEquation([]symbol.Symbol{symbol.Const("A"),
		symbol.Var("x"), symbol.LetterVar("T"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("x"),
			symbol.Var("y"), symbol.Const("B"), symbol.Const("A")})

	if !eq.IsRegularlyOrdered() {
		t.Errorf("TestEquation_CheckSimpleWord failed: A x T y = x y B A is regularly ordered")
		return
	}

	// C T y x A y = y D x A T y B
	eq = NewEquation([]symbol.Symbol{symbol.LetterVar("C"),
		symbol.LetterVar("T"), symbol.Var("y"), symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.LetterVar("D"), symbol.Var("x"),
			symbol.Const("A"), symbol.LetterVar("T"), symbol.Var("y"), symbol.Const("B")})

	if !eq.IsRegularlyOrdered() {
		t.Errorf("TestEquation_CheckSimpleWord C T y x A y = y D x A T y B is regularly ordered")
		return
	}
}

var newEq1 = NewEquation([]symbol.Symbol{
	symbol.Const("A"), symbol.Var("y"), symbol.Var("x")},
	[]symbol.Symbol{symbol.Var("x"), symbol.Var("y"), symbol.Const("A")})

func TestEquation_Apply_FirstForward(t *testing.T) {
	var eq, e, newEq Equation
	var applied bool
	var err error
	// y A x x = x x y A
	eq = NewEquation([]symbol.Symbol{symbol.Var("y"),
		symbol.Const("A"), symbol.Var("x"), symbol.Var("x")},
		[]symbol.Symbol{symbol.Var("x"), symbol.Var("x"), symbol.Var("y"), symbol.Const("A")})

	// x A y = y A x
	e = NewEquation([]symbol.Symbol{symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"),
			symbol.Var("x")})
	e.isEquidecomposable = true

	applied, newEq, err = eq.Apply(e)
	if err != nil {
		t.Errorf("TestEquation_Apply_FirstForward err must be nil: %v", err)
		return
	}
	if !applied {
		t.Errorf("TestEquation_Apply_FirstForward must be applied")
		return
	}
	if !newEq.CheckSameness(&newEq1) {
		t.Errorf("TestEquation_Apply_FirstForward new equation must be: %v", newEq1.String())
		return
	}
}

var newEq2 = NewEquation([]symbol.Symbol{
	symbol.Const("A"), symbol.Var("y")},
	[]symbol.Symbol{symbol.Var("y"), symbol.Const("A")})

func TestEquation_Apply_FirstBackwards(t *testing.T) {
	var eq, e, newEq Equation
	var applied bool
	var err error
	// x y A x = x x y A
	eq = NewEquation([]symbol.Symbol{symbol.Var("x"), symbol.Var("y"),
		symbol.Const("A"), symbol.Var("x")},
		[]symbol.Symbol{symbol.Var("x"), symbol.Var("x"), symbol.Var("y"), symbol.Const("A")})

	// x A y = y A x
	e = NewEquation([]symbol.Symbol{symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"),
			symbol.Var("x")})
	e.isEquidecomposable = true

	applied, newEq, err = eq.Apply(e)
	if err != nil {
		t.Errorf("TestEquation_Apply_FirstBackwards err must be nil: %v", err)
		return
	}
	if !applied {
		t.Errorf("TestEquation_Apply_FirstBackwards must be applied")
		return
	}
	if !newEq.CheckSameness(&newEq2) {
		t.Errorf("TestEquation_Apply_FirstBackwards new equation must be: %v", newEq2.String())
		return
	}
}

var newEq3 = NewEquation([]symbol.Symbol{symbol.Const("A"), symbol.Var("z"), symbol.Var("y"),
	symbol.Const("A"), symbol.Var("x")},
	[]symbol.Symbol{symbol.Var("z"), symbol.Var("y"), symbol.Var("x"),
		symbol.Const("A"), symbol.Const("A")})

func TestEquation_Apply_ThirdForward(t *testing.T) {
	var eq, e, newEq Equation
	var applied bool
	var err error
	// A z x A y = z y x A A
	eq = NewEquation([]symbol.Symbol{symbol.Const("A"), symbol.Var("z"), symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("z"), symbol.Var("y"), symbol.Var("x"),
			symbol.Const("A"), symbol.Const("A")})

	// x A y = y A x
	e = NewEquation([]symbol.Symbol{symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"),
			symbol.Var("x")})
	e.isEquidecomposable = true

	applied, newEq, err = eq.Apply(e)
	if err != nil {
		t.Errorf("TestEquation_Apply_ThirdForward err must be nil: %v", err)
		return
	}
	if !applied {
		t.Errorf("TestEquation_Apply_ThirdForward must be applied")
		return
	}
	if !newEq.CheckSameness(&newEq3) {
		t.Errorf("TestEquation_Apply_ThirdForward new equation must be: %v", newEq3.String())
		return
	}
}

var newEq4 = NewEquation([]symbol.Symbol{symbol.Var("y"),
	symbol.Const("A"), symbol.Var("x"), symbol.Var("z"), symbol.Const("A")},
	[]symbol.Symbol{symbol.Var("z"), symbol.Const("A"), symbol.Var("y"), symbol.Var("x"),
		symbol.Const("A"), symbol.Var("z")})

func TestEquation_Apply_ThirdBackwards(t *testing.T) {
	var eq, e, newEq Equation
	var applied bool
	var err error
	// x A y z A =  z A y x A z
	eq = NewEquation([]symbol.Symbol{symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y"), symbol.Var("z"), symbol.Const("A")},
		[]symbol.Symbol{symbol.Var("z"), symbol.Const("A"), symbol.Var("y"), symbol.Var("x"),
			symbol.Const("A"), symbol.Var("z")})

	// x A y = y A x
	e = NewEquation([]symbol.Symbol{symbol.Var("x"),
		symbol.Const("A"), symbol.Var("y")},
		[]symbol.Symbol{symbol.Var("y"), symbol.Const("A"),
			symbol.Var("x")})
	e.isEquidecomposable = true

	applied, newEq, err = eq.Apply(e)
	if err != nil {
		t.Errorf("TestEquation_Apply_ThirdBackwards err must be nil: %v", err)
		return
	}
	if !applied {
		t.Errorf("TestEquation_Apply_ThirdBackwards must be applied")
		return
	}
	if !newEq.CheckSameness(&newEq4) {
		t.Errorf("TestEquation_Apply_ThirdBackwards new equation must be: %v", newEq4.String())
		return
	}
}
