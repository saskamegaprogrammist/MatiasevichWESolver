package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/symbol"
	"math"
	"math/rand"
	"time"
)

const cycle_range = 10
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

type Solver struct {
	algorithmType int64
	constantsAlph Alphabet
	varsAlph      Alphabet
	wordsAlph     Alphabet
	equation      Equation
	hasSolution   bool
	cycled        bool
	dotWriter     DotWriter
	fullGraph     bool
}

func (solver *Solver) Init(algorithmType string, constantsAlph string, varsAlph string, equation string, fullGraph bool) error {
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
	err = solver.dotWriter.Init()
	if err != nil {
		return fmt.Errorf("error initing solver: %v", err)
	}
	solver.fullGraph = fullGraph
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
	lenLetters := len(alphLetters)
	if lenLetters == 0 {
		return alphabet, nil
	}
	var currentLetter string
	for i := 0; i < lenLetters; i++ {
		sym := alphLetters[i]
		stringSymbol := string(sym)
		if stringSymbol == COMMA {
			if currentLetter == "" {
				return alphabet, fmt.Errorf("empty constant in alphabet: %s", alphabetStr)
			}
			if i+1 != lenLetters && string(alphLetters[i+1]) != SPACE {
				return alphabet, fmt.Errorf("letters must be separated by space: %s", alphabetStr)
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

func (solver *Solver) getAnswer() string {
	if solver.hasSolution {
		return "TRUE"
	}
	if solver.cycled {
		return "CYCLED"
	}
	return "FALSE"
}

func (solver *Solver) Solve() (string, error) {
	tree := Node{
		Number: "0",
		Value:  solver.equation,
	}
	err := solver.dotWriter.StartDOTDescription()
	if err != nil {
		return "", fmt.Errorf("error writing DOT description: %v", err)
	}
	solver.solve(&tree)
	result := solver.getAnswer()
	err = solver.dotWriter.EndDOTDescription()
	if err != nil {
		return result, fmt.Errorf("error writing DOT description: %v", err)
	}
	return result, nil
}

func (solver *Solver) checkEquality(node *Node) bool {
	return node.Value.CheckEquality()
}

func (solver *Solver) checkHasBeen(node *Node) bool {
	tr := node.Parent
	for tr != nil {
		if node.Value.CheckSameness(&tr.Value) {
			solver.dotWriter.WriteDottedEdge(node, tr)
			return true
		}
		tr = tr.Parent
	}
	return false
}

func randStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (solver *Solver) getWord() symbol.Word {
	i := 1
	for {
		jRange := int(math.Pow(float64(len(letterBytes)), float64(i)))
		for j := 0; j < jRange; j++ {
			str := randStr(i)
			if !solver.wordsAlph.Has(str) {
				solver.wordsAlph.AddWord(str)
				return symbol.WordVar(str)
			}
		}
		i++
	}
}

func (solver *Solver) solve(node *Node) {
	solver.dotWriter.WriteNode(node)
	if !solver.fullGraph && solver.hasSolution {
		return
	}
	if len(node.Number) > cycle_range {
		solver.cycled = true
		return
	}
	//fmt.Println(node.Number)
	if solver.checkEquality(node) {
		trueNode := &TrueNode{
			number: "T_" + node.Number,
		}
		solver.dotWriter.WriteInfoNode(trueNode)
		solver.dotWriter.WriteInfoEdge(node, trueNode)
		solver.hasSolution = true
		//fmt.Println("TRUE")
		//fmt.Println(node.Number)
		return
	}
	if solver.checkHasBeen(node) {
		//fmt.Println("HAS BEEN")
		//fmt.Println(node.Number)
		return
	}
	if solver.algorithmType == FINITE {
		if solver.checkFirstRuleFinite(&node.Value) {
			newVals := []symbol.Symbol{node.Value.rightPart[0]}
			eq := node.Value.Substitute(&node.Value.leftPart[0], newVals)
			child := Node{
				Number: "a" + node.Number,
				Parent: node,
				Value:  eq,
			}
			node.Children = []*Node{&child}
			solver.dotWriter.WriteLabelEdge(node, &child, &node.Value.leftPart[0], newVals)
		}
		if solver.checkSecondRuleLeftFinite(&node.Value) {
			newVals := []symbol.Symbol{node.Value.leftPart[0]}
			eq := node.Value.Substitute(&node.Value.rightPart[0], newVals)
			child := Node{
				Number: "b" + node.Number,
				Parent: node,
				Value:  eq,
			}
			node.Children = []*Node{&child}
			solver.dotWriter.WriteLabelEdge(node, &child, &node.Value.rightPart[0], newVals)
		}
		if solver.checkSecondRuleRightFinite(&node.Value) {
			newVals := []symbol.Symbol{node.Value.rightPart[0]}
			eq := node.Value.Substitute(&node.Value.leftPart[0], newVals)
			child := Node{
				Number: "c" + node.Number,
				Parent: node,
				Value:  eq,
			}
			node.Children = []*Node{&child}
			solver.dotWriter.WriteLabelEdge(node, &child, &node.Value.leftPart[0], newVals)
		}
		if solver.checkFourthRuleLeft(&node.Value) {
			newValsFirst := []symbol.Symbol{symbol.Empty()}
			firstEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsFirst)
			firstChild := Node{
				Number: "d" + node.Number,
				Parent: node,
				Value:  firstEquation,
			}
			newValsSecond := []symbol.Symbol{node.Value.leftPart[0], node.Value.rightPart[0]}
			secondEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsSecond)
			secondChild := Node{
				Number: "e" + node.Number,
				Parent: node,
				Value:  secondEquation,
			}
			node.Children = []*Node{&firstChild, &secondChild}
			solver.dotWriter.WriteLabelEdge(node, &firstChild, &node.Value.rightPart[0], newValsFirst)
			solver.dotWriter.WriteLabelEdge(node, &secondChild, &node.Value.rightPart[0], newValsSecond)
		}
		if solver.checkFourthRuleRight(&node.Value) {
			newValsFirst := []symbol.Symbol{symbol.Empty()}
			firstEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsFirst)
			firstChild := Node{
				Number: "f" + node.Number,
				Parent: node,
				Value:  firstEquation,
			}
			newValsSecond := []symbol.Symbol{node.Value.rightPart[0], node.Value.leftPart[0]}
			secondEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsSecond)
			secondChild := Node{
				Number: "g" + node.Number,
				Parent: node,
				Value:  secondEquation,
			}
			node.Children = []*Node{&firstChild, &secondChild}
			solver.dotWriter.WriteLabelEdge(node, &firstChild, &node.Value.rightPart[0], newValsFirst)
			solver.dotWriter.WriteLabelEdge(node, &secondChild, &node.Value.rightPart[0], newValsSecond)
		}
	}
	if solver.checkFirstRule(&node.Value) {
		var newValsFirst []symbol.Symbol
		if solver.algorithmType == INFINITE {
			newValsFirst = []symbol.Symbol{node.Value.rightPart[0], node.Value.leftPart[0]}
		}
		if solver.algorithmType == FINITE {
			word := solver.getWord()
			newValsFirst = []symbol.Symbol{node.Value.rightPart[0], word, node.Value.leftPart[0]}
		}
		firstEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsFirst)
		firstChild := Node{
			Number: node.Number + "1",
			Parent: node,
			Value:  firstEquation,
		}
		var newValsSecond []symbol.Symbol
		if solver.algorithmType == INFINITE {
			newValsSecond = []symbol.Symbol{node.Value.leftPart[0], node.Value.rightPart[0]}

		}
		if solver.algorithmType == FINITE {
			word := solver.getWord()
			newValsSecond = []symbol.Symbol{node.Value.leftPart[0], word, node.Value.rightPart[0]}
		}
		secondEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsSecond)
		secondChild := Node{
			Number: node.Number + "2",
			Parent: node,
			Value:  secondEquation,
		}
		newValsThird := []symbol.Symbol{node.Value.rightPart[0]}
		thirdEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsThird)
		thirdChild := Node{
			Number: node.Number + "3",
			Parent: node,
			Value:  thirdEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild, &thirdChild}
		solver.dotWriter.WriteLabelEdge(node, &firstChild, &node.Value.leftPart[0], newValsFirst)
		solver.dotWriter.WriteLabelEdge(node, &secondChild, &node.Value.rightPart[0], newValsSecond)
		solver.dotWriter.WriteLabelEdge(node, &thirdChild, &node.Value.leftPart[0], newValsThird)
	}

	if solver.checkSecondRuleLeft(&node.Value) {
		newValsFirst := []symbol.Symbol{symbol.Empty()}
		firstEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsFirst)
		firstChild := Node{
			Number: node.Number + "4",
			Parent: node,
			Value:  firstEquation,
		}
		var newValsSecond []symbol.Symbol
		newValsSecond = []symbol.Symbol{node.Value.leftPart[0], node.Value.rightPart[0]}
		secondEquation := node.Value.Substitute(&node.Value.rightPart[0], newValsSecond)
		secondChild := Node{
			Number: node.Number + "5",
			Parent: node,
			Value:  secondEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild}
		solver.dotWriter.WriteLabelEdge(node, &firstChild, &node.Value.rightPart[0], newValsFirst)
		solver.dotWriter.WriteLabelEdge(node, &secondChild, &node.Value.rightPart[0], newValsSecond)
	}
	if solver.checkSecondRuleRight(&node.Value) {
		newValsFirst := []symbol.Symbol{symbol.Empty()}
		firstEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsFirst)
		firstChild := Node{
			Number: node.Number + "6",
			Parent: node,
			Value:  firstEquation,
		}
		var newValsSecond []symbol.Symbol
		newValsSecond = []symbol.Symbol{node.Value.rightPart[0], node.Value.leftPart[0]}
		secondEquation := node.Value.Substitute(&node.Value.leftPart[0], newValsSecond)
		secondChild := Node{
			Number: node.Number + "7",
			Parent: node,
			Value:  secondEquation,
		}
		node.Children = []*Node{&firstChild, &secondChild}
		solver.dotWriter.WriteLabelEdge(node, &firstChild, &node.Value.leftPart[0], newValsFirst)
		solver.dotWriter.WriteLabelEdge(node, &secondChild, &node.Value.leftPart[0], newValsSecond)

	}
	if solver.checkThirdRuleLeft(&node.Value) || solver.checkThirdRuleRight(&node.Value) {
		eq := node.Value.SubstituteVarsWithEmpty()
		child := Node{
			Number: node.Number + "8",
			Parent: node,
			Value:  eq,
		}
		node.Children = []*Node{&child}
		solver.dotWriter.WriteEdge(node, &child)
	}
	//for _, child := range node.Children {
	//	child.Print()
	//}
	for _, child := range node.Children {
		solver.solve(child)
	}
	if len(node.Children) == 0 {
		falseNode := &FalseNode{number: "F_" + node.Number}
		solver.dotWriter.WriteInfoNode(falseNode)
		solver.dotWriter.WriteInfoEdge(node, falseNode)
	}
}

func (solver *Solver) checkFirstRule(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsVar(eq.leftPart[0]) && symbol.IsVar(eq.rightPart[0])
}

func (solver *Solver) checkFirstRuleFinite(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsWord(eq.leftPart[0]) && symbol.IsWord(eq.rightPart[0])
}

func (solver *Solver) checkSecondRuleRight(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsVar(eq.leftPart[0]) && symbol.IsConst(eq.rightPart[0])
}

func (solver *Solver) checkSecondRuleLeft(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsConst(eq.leftPart[0]) && symbol.IsVar(eq.rightPart[0])
}

func (solver *Solver) checkSecondRuleRightFinite(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsWord(eq.leftPart[0]) && symbol.IsConst(eq.rightPart[0])
}
func (solver *Solver) checkSecondRuleLeftFinite(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsConst(eq.leftPart[0]) && symbol.IsWord(eq.rightPart[0])
}

func (solver *Solver) checkThirdRuleRight(eq *Equation) bool {
	return eq.IsRightEmpty()
}
func (solver *Solver) checkThirdRuleLeft(eq *Equation) bool {
	return eq.IsLeftEmpty()
}

func (solver *Solver) checkFourthRuleLeft(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsWord(eq.leftPart[0]) && symbol.IsVar(eq.rightPart[0])
}

func (solver *Solver) checkFourthRuleRight(eq *Equation) bool {
	return eq.leftLength > 0 && eq.rightLength > 0 &&
		symbol.IsVar(eq.leftPart[0]) && symbol.IsWord(eq.rightPart[0])
}
