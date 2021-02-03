package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/symbol"
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

var test2InitEqErrorMessage = "error matching alphabet: no match for word: o"

func Test_InitEq_Error_2(t *testing.T) {
	var eq Equation
	err := eq.Init("a=o", &constAlph, &varsAlph)
	if err == nil {
		t.Errorf("Test_InitEq_Error_2 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test2InitEqErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_InitEq_Error_2 failed: wrong error message")
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
	err := eq.Init("aaaaaaabx=aaaabxp", &constAlphMisleading, &varsAlphMisleading)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_InitEq_1 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[0] != symbol.Var("aaaaa") || eq.leftPart[1] != symbol.Const("aabx") ||
		eq.rightPart[0] != symbol.Const("aaaa") || eq.rightPart[1] != symbol.Var("bxp") {
		t.Errorf("Test_InitEq_1 failed: wrong eq parsing : ")
		eq.Print()
	}
}

func Test_InitEq_2(t *testing.T) {
	var eq Equation
	err := eq.Init("aaux=yibp", &constAlph, &varsAlph)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Test_InitEq_2 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[0] != symbol.Const("a") || eq.leftPart[1] != symbol.Const("a") ||
		eq.leftPart[2] != symbol.Var("ux") ||
		eq.rightPart[0] != symbol.Var("yi") || eq.rightPart[1] != symbol.Const("b") ||
		eq.rightPart[2] != symbol.Var("p") {
		t.Errorf("Test_InitEq_2 failed: wrong eq parsing : ")
		eq.Print()
	}
}

var constAlphNew = Alphabet{
	words:         []string{"a", "b", "c", "d"},
	size:          4,
	maxWordLength: 1,
}
var varsAlphNew = Alphabet{
	words:         []string{"u", "v", "x"},
	size:          3,
	maxWordLength: 1,
}

func TestEquation_Reduce_1(t *testing.T) {
	var eq Equation
	err := eq.Init("abx=$abv", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Reduce_1 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[0] != symbol.Const("a") || eq.leftPart[1] != symbol.Const("b") ||
		eq.leftPart[2] != symbol.Var("x") ||
		eq.rightPart[0] != symbol.Empty() || eq.rightPart[1] != symbol.Const("a") ||
		eq.rightPart[2] != symbol.Const("b") || eq.rightPart[3] != symbol.Var("v") {
		t.Errorf("TestEquation_Reduce_1 failed: wrong eq parsing : ")
		eq.Print()
		return
	}
	eq.Reduce()
	if eq.leftLength == 0 || eq.rightLength == 0 {
		t.Errorf("TestEquation_Reduce_1 failed: wrong reduce result : len shouldn\\'t be nil")
	} else {
		if eq.leftPart[0] != symbol.Var("x") || eq.rightPart[0] != symbol.Var("v") {
			t.Errorf("TestEquation_Reduce_1 failed: wrong reduce result : ")
			eq.Print()
		}
	}
}

func TestEquation_Reduce_2(t *testing.T) {
	var eq Equation
	err := eq.Init("x$=x$av", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Reduce_2 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[0] != symbol.Var("x") || eq.leftPart[1] != symbol.Empty() ||
		eq.rightPart[0] != symbol.Var("x") || eq.rightPart[2] != symbol.Const("a") ||
		eq.rightPart[1] != symbol.Empty() || eq.rightPart[3] != symbol.Var("v") {
		t.Errorf("TestEquation_Reduce_2 failed: wrong eq parsing : ")
		eq.Print()
		return
	}
	eq.Reduce()
	if eq.leftLength != 0 {
		t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result : left len should be nil")
		return
	}
	if eq.rightLength != 2 {
		t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result : right len should be: %d", 2)
	} else {
		if eq.rightPart[0] != symbol.Const("a") || eq.rightPart[1] != symbol.Var("v") {
			t.Errorf("TestEquation_Reduce_2 failed: wrong reduce result : ")
			eq.Print()
		}
	}
}

func TestEquation_Substitute_1(t *testing.T) {
	var eq Equation
	err := eq.Init("abx=vb", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_1 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[0] != symbol.Const("a") || eq.leftPart[1] != symbol.Const("b") ||
		eq.leftPart[2] != symbol.Var("x") ||
		eq.rightPart[0] != symbol.Var("v") || eq.rightPart[1] != symbol.Const("b") {
		t.Errorf("TestEquation_Substitute_1 failed: wrong eq parsing : ")
		eq.Print()
		return
	}
	resultEq := eq.Substitute(&eq.rightPart[0], []symbol.Symbol{eq.leftPart[0], eq.rightPart[0]})
	if resultEq.leftLength != 2 || resultEq.rightLength != 2 {
		t.Errorf("TestEquation_Substitute_1 failed: wrong reduce result : ")
		resultEq.Print()
	} else {
		if resultEq.leftPart[0] != symbol.Const("b") ||
			resultEq.leftPart[1] != symbol.Var("x") ||
			resultEq.rightPart[0] != symbol.Var("v") || resultEq.rightPart[1] != symbol.Const("b") {
			t.Errorf("TestEquation_Substitute_1 failed: wrong reduce result : ")
			resultEq.Print()
		}
	}
}

func TestEquation_Substitute_2(t *testing.T) {
	var eq Equation
	err := eq.Init("xbxa=vb", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_2 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[3] != symbol.Const("a") || eq.leftPart[1] != symbol.Const("b") ||
		eq.leftPart[2] != symbol.Var("x") || eq.leftPart[0] != symbol.Var("x") ||
		eq.rightPart[0] != symbol.Var("v") || eq.rightPart[1] != symbol.Const("b") {
		t.Errorf("TestEquation_Substitute_2 failed: wrong eq parsing : ")
		eq.Print()
		return
	}
	resultEq := eq.Substitute(&eq.leftPart[0], []symbol.Symbol{eq.rightPart[0], eq.leftPart[0]})
	if resultEq.leftLength != 5 || resultEq.rightLength != 1 {
		t.Errorf("TestEquation_Substitute_2 failed: wrong reduce result : ")
		resultEq.Print()
	} else {
		if resultEq.leftPart[0] != symbol.Var("x") ||
			resultEq.leftPart[1] != symbol.Const("b") ||
			resultEq.leftPart[2] != symbol.Var("v") || resultEq.leftPart[3] != symbol.Var("x") ||
			resultEq.leftPart[4] != symbol.Const("a") ||
			resultEq.rightPart[0] != symbol.Const("b") {
			t.Errorf("TestEquation_Substitute_2 failed: wrong reduce result : ")
			resultEq.Print()
		}
	}
}

func TestEquation_Substitute_3(t *testing.T) {
	var eq Equation
	err := eq.Init("xav=", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquation_Substitute_3 failed: error shouldn be nil")
		return
	}
	if eq.leftPart[1] != symbol.Const("a") ||
		eq.leftPart[2] != symbol.Var("v") || eq.leftPart[0] != symbol.Var("x") ||
		eq.rightPart[0] != symbol.Empty() {
		t.Errorf("TestEquation_Substitute_3 failed: wrong eq parsing : ")
		eq.Print()
		return
	}
	resultEq := eq.SubstituteVarsWithEmpty()
	if resultEq.leftLength != 1 || resultEq.rightLength != 1 {
		t.Errorf("TestEquation_Substitute_3 failed: wrong reduce result : ")
		resultEq.Print()
	} else {
		if resultEq.leftPart[0] != symbol.Const("a") ||
			resultEq.rightPart[0] != symbol.Empty() {
			t.Errorf("TestEquation_Substitute_3 failed: wrong reduce result : ")
			resultEq.Print()
		}
	}
}
