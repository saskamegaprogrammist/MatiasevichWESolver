package solver

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

// equation rules

func checkFirstRule(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsVar(eq.LeftPart.Symbols[0]) && symbol.IsVar(eq.RightPart.Symbols[0])
}

func checkFirstRuleFinite(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsLetter(eq.LeftPart.Symbols[0]) && symbol.IsLetter(eq.RightPart.Symbols[0])
}

func checkSecondRuleRight(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsVar(eq.LeftPart.Symbols[0]) && symbol.IsConst(eq.RightPart.Symbols[0])
}

func checkSecondRuleLeft(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsConst(eq.LeftPart.Symbols[0]) && symbol.IsVar(eq.RightPart.Symbols[0])
}

func checkSecondRuleRightFinite(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsLetter(eq.LeftPart.Symbols[0]) && symbol.IsConst(eq.RightPart.Symbols[0])
}
func checkSecondRuleLeftFinite(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsConst(eq.LeftPart.Symbols[0]) && symbol.IsLetter(eq.RightPart.Symbols[0])
}

func checkThirdRuleRight(eq *equation.Equation) bool {
	return eq.IsRightEmpty()
}
func checkThirdRuleLeft(eq *equation.Equation) bool {
	return eq.IsLeftEmpty()
}

func checkFourthRuleLeft(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsLetter(eq.LeftPart.Symbols[0]) && symbol.IsVar(eq.RightPart.Symbols[0])
}

func checkFourthRuleRight(eq *equation.Equation) bool {
	return eq.LeftPart.Length > 0 && eq.RightPart.Length > 0 &&
		symbol.IsVar(eq.LeftPart.Symbols[0]) && symbol.IsLetter(eq.RightPart.Symbols[0])
}

// node rules

func checkEquality(node *Node) bool {
	return node.value.Equation().CheckEquality()
}

func checkInequality(node *Node) bool {
	return node.value.Equation().CheckInequality()
}

func checkHasBeen(node *Node) (bool, *Node) {
	tr := node.parent
	for tr != nil {
		if node.value.Equals(tr.value) {
			return true, tr
		}
		tr = tr.parent
	}
	return false, nil
}

// node system rules

func checkSystemEquality(nodeSystem *NodeSystem) bool {
	return nodeSystem.Value.CheckEquality()
}

func checkSystemInequality(nodeSystem *NodeSystem) bool {
	return nodeSystem.Value.CheckInequality()
}

func checkSystemHasBeen(nodeSystem *NodeSystem) (bool, *NodeSystem) {
	tr := nodeSystem.Parent
	for tr != nil {
		if nodeSystem.Value.CheckSameness(&tr.Value) {
			return true, tr
		}
		tr = tr.Parent
	}
	return false, nil
}
