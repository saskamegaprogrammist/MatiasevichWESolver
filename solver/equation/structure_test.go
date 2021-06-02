package equation

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"testing"
)

func Test_MergeStructures_1(t *testing.T) {
	var s1 = NewStructure([]symbol.Symbol{symbol.Const("a"), symbol.Var("x"), symbol.Const("a"),
		symbol.Const("b"), symbol.LetterVar("t"), symbol.LetterVar("m")})
	var s2 = NewStructure([]symbol.Symbol{symbol.Const("p"), symbol.Var("x"), symbol.Const("a"),
		symbol.Const("b"), symbol.LetterVar("t"), symbol.Var("y")})
	var m = MergeStructures(&s1, &s2)
	if m.letters[symbol.LetterVar("t")] != 2 || m.letters[symbol.LetterVar("m")] != 1 ||
		m.vars[symbol.Var("x")] != 2 || m.vars[symbol.Var("y")] != 1 ||
		m.consts[symbol.Const("a")] != 3 || m.consts[symbol.Const("b")] != 2 ||
		m.consts[symbol.Const("p")] != 1 {
		t.Errorf("Test_MergeStructures_1 failed: wrong merged struct")
		return
	}
	if m.varsLen != 3 || m.lettersLen != 3 || m.constsLen != 6 {
		t.Errorf("Test_MergeStructures_1 failed: wrong merged struct")
		return
	}
	if m.VarsRangeLen() != 2 || m.LettersRangeLen() != 2 || m.ConstsRangeLen() != 3 {
		t.Errorf("Test_MergeStructures_1 failed: wrong merged struct")
		return
	}
}
