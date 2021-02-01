package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/symbol"
)

type Solver struct {
	algorithmType int64
	constantsAlph Alphabet
	varsAlph      Alphabet
	equation      Equation
	hasSolution   bool
}

func (solver *Solver) Init(algorithmType string, constantsAlph string, varsAlph string, equation string) error {
	var err error
	intType, err := matchAlgorithmType(algorithmType)
	if err != nil {
		return fmt.Errorf("error matching alphabet type: %v", err)
	}
	solver.algorithmType = intType
	constAlphabet, err := solver.parseAlphabet(constantsAlph)
	if err != nil {
		return fmt.Errorf("error parsing constants: %v", err)
	}
	solver.constantsAlph = constAlphabet
	varsAlphabet, err := solver.parseAlphabet(varsAlph)
	if err != nil {
		return fmt.Errorf("error parsing vars: %v", err)
	}
	solver.varsAlph = varsAlphabet
	err = solver.equation.Init(equation, &constAlphabet, &varsAlphabet)
	if err != nil {
		return fmt.Errorf("error parsing equation: %v", err)
	}
	return nil
}

func (solver *Solver) parseAlphabet(alphabetStr string) (Alphabet, error) {
	var alphabet Alphabet
	var maxWordLength int
	lenAlph := len(alphabetStr)
	if alphabetStr[0:1] != OPENBR || alphabetStr[lenAlph-1:] != CLOSEBR {
		return alphabet, fmt.Errorf("invalid constants alphabet: %s", alphabetStr)
	}
	alphLetters := alphabetStr[1 : lenAlph-1]
	var currentLetter string
	for _, symbol := range alphLetters {
		stringSymbol := string(symbol)
		if stringSymbol == COMMA {
			if currentLetter == "" {
				return alphabet, fmt.Errorf("empty constant in alphabet: %s", alphabetStr)
			}
			alphabet.AddWord(currentLetter)
			if len(currentLetter) > maxWordLength {
				maxWordLength = len(currentLetter)
			}
			currentLetter = ""
		} else {
			currentLetter += stringSymbol
		}
	}
	if currentLetter == "" {
		return alphabet, fmt.Errorf("empty constant in alphabet: %s", alphabetStr)
	}
	alphabet.AddWord(currentLetter)
	if len(currentLetter) > maxWordLength {
		maxWordLength = len(currentLetter)
	}
	alphabet.maxWordLength = maxWordLength
	return alphabet, nil
}

func (solver *Solver) Solve() (bool, error) {
	tree := Node{
		Value: solver.equation,
	}
	if solver.algorithmType == FINITE {
		solver.solveFinite(&tree)
		return solver.hasSolution, nil
	} else if solver.algorithmType == INFINITE {
		return solver.solveInfinite(&tree)
	}
	return false, fmt.Errorf("wrong algorithm type: %d", solver.algorithmType)
}

func (solver *Solver) solveInfinite(node *Node) (bool, error) {

	return false, nil
}

func (solver *Solver) checkEquality(node *Node) bool {
	return node.Value.CheckEquality()
}

func (solver *Solver) checkHasBeen(node *Node) bool {
	tr := node.Parent
	for tr != nil {
		if node.Value.CheckSameness(&tr.Value) {
			return true
		}
		tr = tr.Parent
	}
	return false
}

func (solver *Solver) solveFinite(node *Node) {
	fmt.Println(node.Value)
	if solver.checkEquality(node) {
		solver.hasSolution = true
		fmt.Println("TRUE")
		fmt.Println(node.Value)
		return
	}
	if solver.checkHasBeen(node) {
		return
	}
	if solver.checkFirstRule(&node.Value) {
		newValsFirst := []symbol.Symbol{node.Value.rightPart[0], node.Value.leftPart[0]}
		firstEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsFirst)
		firstChild := Node{
			Parent: node,
			Value:  firstEquation,
		}
		newValsSecond := []symbol.Symbol{node.Value.leftPart[0], node.Value.rightPart[0]}
		secondEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsSecond)
		secondChild := Node{
			Parent: node,
			Value:  secondEquation,
		}
		newValsThird := []symbol.Symbol{node.Value.rightPart[0]}
		thirdEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsThird)
		thirdChild := Node{
			Parent: node,
			Value:  thirdEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild, &thirdChild}
	}
	if solver.checkSecondRuleLeft(&node.Value) {
		newValsFirst := []symbol.Symbol{symbol.Empty()}
		firstEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsFirst)
		firstChild := Node{
			Parent: node,
			Value:  firstEquation,
		}
		newValsSecond := []symbol.Symbol{node.Value.leftPart[0], node.Value.rightPart[0]}
		secondEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsSecond)
		secondChild := Node{
			Parent: node,
			Value:  secondEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild}
	}
	if solver.checkSecondRuleRight(&node.Value) {
		newValsFirst := []symbol.Symbol{symbol.Empty()}
		firstEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsFirst)
		firstChild := Node{
			Parent: node,
			Value:  firstEquation,
		}
		newValsSecond := []symbol.Symbol{node.Value.rightPart[0], node.Value.leftPart[0]}
		secondEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsSecond)
		secondChild := Node{
			Parent: node,
			Value:  secondEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild}
	}
	if solver.checkThirdRuleLeft(&node.Value) || solver.checkThirdRuleRight(&node.Value) {
		eq := node.Value.SubstituteVarsWithEmpty()
		child := Node{
			Parent: node,
			Value:  eq,
		}
		node.Children = []*Node{&child}
	}
	for _, child := range node.Children {
		solver.solveFinite(child)
	}
}

func (solver *Solver) checkFirstRule(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsVar(eq.leftPart[0]) && symbol.IsVar(eq.rightPart[0])
}

func (solver *Solver) checkSecondRuleRight(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsVar(eq.leftPart[0]) && symbol.IsConst(eq.rightPart[0])
}
func (solver *Solver) checkSecondRuleLeft(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsConst(eq.leftPart[0]) && symbol.IsVar(eq.rightPart[0])
}

func (solver *Solver) checkThirdRuleRight(eq *Equation) bool {
	return eq.IsRightEmpty()
}
func (solver *Solver) checkThirdRuleLeft(eq *Equation) bool {
	return eq.IsLeftEmpty()
}
