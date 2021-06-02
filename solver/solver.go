package solver

import (
	"fmt"
	"github.com/google/logger"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const cycle_range = 100
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

const MAGIC_PREFIX = "MAGIC"
const EMPTY = ""

type Solver struct {
	algorithmType int64
	constantsAlph equation.Alphabet
	varsAlph      equation.Alphabet
	wordsAlph     equation.Alphabet
	equation      equation.Equation
	hasSolution   bool
	cycled        bool
	dotWriter     DotWriter
	printOptions  PrintOptions
	solveOptions  SolveOptions
	simplifier    Simplifier
}

func (solver *Solver) InitWoEquation(constantsAlph string, varsAlph string,
	printOptions PrintOptions,
	solveOptions SolveOptions) error {
	return solver.Init(constantsAlph, varsAlph, "", printOptions, solveOptions)
}

func (solver *Solver) Init(constantsAlph string, varsAlph string, eq string,
	printOptions PrintOptions,
	solveOptions SolveOptions) error {
	var err error
	intType, err := matchAlgorithmType(solveOptions.AlgorithmMode)
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
	if eq != EMPTY {
		err = solver.equation.Init(eq, &constAlphabet, &varsAlphabet)
		if err != nil {
			return fmt.Errorf("error parsing equation: %v", err)
		}
	}

	solver.printOptions = printOptions
	solver.solveOptions = solveOptions

	if solveOptions.NeedsSimplification {
		err = solver.simplifier.Init(constantsAlph, varsAlph, PrintOptions{}, SolveOptions{
			LengthAnalysis:             false,
			SplitByEquidecomposability: false,
			CycleRange:                 10,
			FullGraph:                  true,
			FullSystem:                 false,
			AlgorithmMode:              solver.solveOptions.AlgorithmMode,
			SaveLettersSubstitutions:   true,
			NeedsSimplification:        false,
		})

		if err != nil {
			return fmt.Errorf("error initing simplifier: %v", err)
		}
	}

	// TODO: change const value
	solver.solveOptions.SaveLettersSubstitutions = true

	if solver.solveOptions.CycleRange == 0 {
		solver.solveOptions.CycleRange = cycle_range
	}
	solver.equation.Print()
	fmt.Println(solveOptions.AlgorithmMode)
	return nil
}

func (solver *Solver) SetEquationString(eq string) error {
	var err error
	err = solver.equation.Init(eq, &solver.constantsAlph, &solver.varsAlph)
	if err != nil {
		return fmt.Errorf("error parsing equation: %v", err)
	}
	return nil
}

func (solver *Solver) SetEquation(eq equation.Equation) error {
	if err := eq.Check(&solver.constantsAlph, &solver.varsAlph, &solver.wordsAlph); err != nil {
		return fmt.Errorf("equation doesn't belong: %v", err)
	}
	solver.equation = eq
	return nil
}

func (solver *Solver) setLetterAlphabet(alphabetStr string) error {
	lettersAlphabet, err := solver.parseAlphabet(alphabetStr)
	if err != nil {
		return fmt.Errorf("error parsing constants: %v", err)
	}
	solver.wordsAlph = lettersAlphabet
	return nil
}

func (solver *Solver) parseAlphabet(alphabetStr string) (equation.Alphabet, error) {
	var alphabet equation.Alphabet
	var maxWordLength int
	lenAlph := len(alphabetStr)
	if lenAlph < 2 || alphabetStr[0:1] != equation.OPENBR || alphabetStr[lenAlph-1:] != equation.CLOSEBR {
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
		if stringSymbol == equation.COMMA {
			if currentLetter == "" {
				return alphabet, fmt.Errorf("empty constant in alphabet: %s", alphabetStr)
			}
			if i+1 != lenLetters && string(alphLetters[i+1]) != equation.SPACE {
				return alphabet, fmt.Errorf("letters must be separated by space: %s", alphabetStr)
			} else {
				i++
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
	alphabet.SetMaxWordLength(maxWordLength)
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

// help functions to create letter alphabet

func randStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (solver *Solver) getLetter() symbol.Letter {
	i := 1
	for {
		jRange := int(math.Pow(float64(len(letterBytes)), float64(i)))
		for j := 0; j < jRange; j++ {
			str := randStr(i)
			if !solver.wordsAlph.Has(str) {
				solver.wordsAlph.AddWord(str)
				return symbol.LetterVar(str)
			}
		}
		i++
	}
}

func (solver *Solver) Solve() (string, time.Duration, error) {
	if solver.equation.IsEmpty() {
		return "", 0, fmt.Errorf("no equation was set")
	}
	var duration time.Duration
	var err error
	if solver.printOptions.Dot {
		err := solver.setWriter(solver.equation)
		if err != nil {
			return "", duration, fmt.Errorf("error setting writer: %v", err)
		}
		err = solver.dotWriter.StartDOTDescription()
		if err != nil {
			return "", duration, fmt.Errorf("error writing DOT description: %v", err)
		}
		defer func() {
			err = solver.dotWriter.EndDOTDescription()
			if err != nil {
				logger.Errorf("error writing DOT description: %v", err)
			}
			err = solver.dotWriter.CreateFiles(solver.printOptions.Png)
			if err != nil {
				logger.Errorf("error creating description files: %v", err)
			}
		}()
	}
	if solver.solveOptions.SplitByEquidecomposability {
		if solver.equation.IsQuadratic() {
			var hasSolution = true
			var hasCycled bool
			system := solver.equation.SplitByEquidecomposability()
			var sumTime time.Duration
			for i, eq := range system.Equations {
				duration, err = solver.solveEquationTimes(eq, i)
				if err != nil {
					return "", duration, fmt.Errorf("error solving equation %d: %v", i, err)
				}
				sumTime += duration
				if !solver.hasSolution {
					if !solver.solveOptions.FullSystem {
						return solver.getAnswer(), sumTime, nil
					} else {
						hasSolution = false
					}
				}
				if solver.cycled {
					hasCycled = solver.cycled
				}
			}
			solver.cycled = hasCycled
			solver.hasSolution = hasSolution
			return solver.getAnswer(), sumTime, err
		} else {
			duration, err = solver.solveEquationAsSystem(solver.equation)
			if err != nil {
				return "", duration, fmt.Errorf("error solving regularly ordered equation as system : %v", err)
			}
			result := solver.getAnswer()
			return result, duration, nil
		}
	} else {
		duration, err = solver.solveEquation(solver.equation)
		if err != nil {
			return "", duration, fmt.Errorf("error solving equation: %v", err)
		}
		result := solver.getAnswer()
		return result, duration, nil
	}
}

func (solver *Solver) clear() {
	solver.cycled = false
	solver.hasSolution = false
}

func (solver *Solver) setWriter(equation equation.Equation) error {
	err := solver.dotWriter.Init(solver.solveOptions.AlgorithmMode, equation.String(), solver.printOptions.OutputDir)
	if err != nil {
		return fmt.Errorf("error initing writer: %v", err)
	}
	return nil
}

func (solver *Solver) solveEquationTimes(equation equation.Equation, times int) (time.Duration, error) {
	var err error
	timeStart := time.Now()
	solver.clear() // if we are solving a system, we should clear solving results

	magicPrefix := strings.Repeat(MAGIC_PREFIX, times)
	tree := NewTreeWEquation(magicPrefix+"0", equation)
	err = solver.solveSystem(&tree)
	if err != nil {
		return 0, fmt.Errorf("error solving equation: %v", err)
	}
	tree.SetWasUnfolded()
	measuredTime := time.Since(timeStart)
	//err = solver.simplifier.Simplify(&tree)
	//if err != nil {
	//	return measuredTime, fmt.Errorf("error simplifing eq: %v", err)
	//}

	if solver.printOptions.Dot {
		err = solver.createGraphDescription(&tree)
		if err != nil {
			return measuredTime, fmt.Errorf("error creating graph desc: %v", err)
		}
	}
	return measuredTime, nil
}

func (solver *Solver) solveEquationsSystem(es equation.EquationsSystem) (time.Duration, error) {
	var err error
	timeStart := time.Now()

	tree := NewTreeWEquationsSystem("0", es)
	err = solver.solveSystem(&tree)
	if err != nil {
		return 0, fmt.Errorf("error solving equations system: %v", err)
	}
	tree.SetWasUnfolded()
	measuredTime := time.Since(timeStart)

	if solver.printOptions.Dot {
		err = solver.createGraphDescription(&tree)
		if err != nil {
			return measuredTime, fmt.Errorf("error creating graph desc: %v", err)
		}
	}
	return measuredTime, nil
}

func (solver *Solver) solveEquation(equation equation.Equation) (time.Duration, error) {
	return solver.solveEquationTimes(equation, 0)
}

func (solver *Solver) createGraphDescription(tree *Node) error { // if we are solving a system, we should clear solving results
	var err error
	err = solver.describeGraph(tree)
	if err != nil {
		return fmt.Errorf("error solving equation: %v", err)
	}
	return nil
}

func checkRuleForSystem(nodeSystem *NodeSystem, rule func(eq equation.Equation) bool, each bool) bool {
	for _, eq := range nodeSystem.Value.Equations {
		if rule(eq) {
			if !each {
				return true
			}
		} else if each {
			return false
		}
	}
	return true
}

func (solver *Solver) createFalseNode(node *Node, falseType int) {
	falseNode := &FalseNode{
		number:    "F_" + node.number,
		falseType: falseType,
	}
	node.infoChild = falseNode
	node.SetHasFalseChildren()
	//fmt.Println("___FALSE")
}

func (solver *Solver) createTrueNode(node *Node) {
	trueNode := &TrueNode{
		number: "T_" + node.number,
	}
	node.SetHasTrueChildren()
	node.infoChild = trueNode
	solver.hasSolution = true
	//fmt.Println("TRUE")
	//fmt.Println(node.number)
}

func newEquationSystemWithSubstitution(oldNode *Node, substitution *equation.Substitution, number string) Node {
	newValue := oldNode.value.Substitute(substitution)
	return NewNodeWEquationsSystem(*substitution, number, oldNode, newValue)
}

func (solver *Solver) simplifyNode(node *Node) (bool, error) {
	node.value.Simplify()
	node.value.Reduce()
	regOrdered, simple, err := node.value.SplitIntoRegOrdered()
	if err != nil {
		return false, fmt.Errorf("error splitting into reg ordered and simple: %v", err)
	}
	if regOrdered.IsEmpty() {
		return false, nil
	}
	needsSimplification, err := regOrdered.NeedsSimplification()
	if err != nil {
		return false, fmt.Errorf("error checking if needs simplification: %v", err)
	}
	if !needsSimplification {
		return false, nil
	}
	tree := NewTreeWEquationsSystem("0", regOrdered)
	solver.solveOptions.FullGraph = true
	err = solver.simplifier.Simplify(&tree)
	if err != nil {
		return false, fmt.Errorf("error simplifying: %v", err)
	}
	if !tree.HasTrueChildren() {
		solver.createTrueNode(node)
		return true, nil
	}
	var newEs equation.EquationsSystem
	if tree.simplified.IsConjunction() || tree.simplified.IsSingleEquation() {
		newEs := equation.NewConjunctionSystem([]equation.EquationsSystem{newEs, simple})
		newEs.Simplify()
		child := NewNodeWEquationsSystem(equation.Substitution{},
			"s"+node.number, node, newEs)
		node.SetChildren([]*Node{&child})
	} else if tree.simplified.IsDisjunction() {
		var newChildNodes = make([]*Node, 0)
		for i, c := range tree.simplified.Compounds() {
			newEs := equation.NewConjunctionSystem([]equation.EquationsSystem{c, simple})
			newEs.Simplify()
			child := NewNodeWEquationsSystem(equation.Substitution{},
				"s"+node.number+strconv.Itoa(i), node, newEs)
			newChildNodes = append(newChildNodes, &child)
		}
		node.SetChildren(newChildNodes)
	} else {
		child := NewNodeWEquationsSystem(equation.Substitution{},
			"s"+node.number, node, simple)
		node.SetChildren([]*Node{&child})
	}
	return true, nil
}

func (solver *Solver) solveSystem(node *Node) error {
	var err error
	if !solver.solveOptions.FullGraph && solver.hasSolution {
		return nil
	}
	if len(node.number) > solver.solveOptions.CycleRange {
		solver.cycled = true
		return nil
	}

	//fmt.Println(node.number)
	if checkInequality(node) {
		solver.createFalseNode(node, REGULAR_FALSE)
		return nil
	}
	if checkEquality(node) {
		solver.createTrueNode(node)
		if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
			nls := NewLetterSubstitution(node)
			if node.Substitution().IsToLetter() {
				nls.AddSubstToHead(node.substitution)
			}
			node.letterSubstitutions = append(node.letterSubstitutions, &nls)
		}
		return nil
	}

	hasBeen, tr := checkHasBeen(node)
	if hasBeen {
		node.SetHasBackCycle()
		node.SetChildren([]*Node{tr})
		tr.AddParentFromBackCycle(node)
		return nil
	}

	var simplified bool

	if solver.solveOptions.NeedsSimplification {
		simplified, err = solver.simplifyNode(node)
		if err != nil {
			return fmt.Errorf("error simplifying node: %v", err)
		}
	}

	if !simplified {
		firstEquation := node.value.Equation()

		if solver.algorithmType == FINITE {
			if checkFirstRuleFinite(firstEquation) {
				substitute := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{firstEquation.RightPart.Symbols[0]})
				child := newEquationSystemWithSubstitution(node, &substitute, "a"+node.number)
				node.SetChildren([]*Node{&child})
			}
			if checkSecondRuleLeftFinite(firstEquation) {
				substitute := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], []symbol.Symbol{firstEquation.LeftPart.Symbols[0]})
				child := newEquationSystemWithSubstitution(node, &substitute, "b"+node.number)
				node.SetChildren([]*Node{&child})
			}
			if checkSecondRuleRightFinite(firstEquation) {
				substitute := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{firstEquation.RightPart.Symbols[0]})
				child := newEquationSystemWithSubstitution(node, &substitute, "c"+node.number)
				node.SetChildren([]*Node{&child})
			}
			if checkFourthRuleLeft(firstEquation) {
				substFirst := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
				firstChild := newEquationSystemWithSubstitution(node, &substFirst, "d"+node.number)

				substSecond := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], []symbol.Symbol{firstEquation.LeftPart.Symbols[0], firstEquation.RightPart.Symbols[0]})
				secondChild := newEquationSystemWithSubstitution(node, &substSecond, "e"+node.number)

				node.SetChildren([]*Node{&firstChild, &secondChild})
			}
			if checkFourthRuleRight(firstEquation) {
				substFirst := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
				firstChild := newEquationSystemWithSubstitution(node, &substFirst, "f"+node.number)

				substSecond := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{firstEquation.RightPart.Symbols[0], firstEquation.LeftPart.Symbols[0]})
				secondChild := newEquationSystemWithSubstitution(node, &substSecond, "g"+node.number)

				node.SetChildren([]*Node{&firstChild, &secondChild})
			}
		}
		if checkFirstRule(firstEquation) {
			var newValsFirst []symbol.Symbol
			if solver.algorithmType == STANDARD {
				newValsFirst = []symbol.Symbol{firstEquation.RightPart.Symbols[0], firstEquation.LeftPart.Symbols[0]}
			}
			if solver.algorithmType == FINITE {
				word := solver.getLetter()
				newValsFirst = []symbol.Symbol{firstEquation.RightPart.Symbols[0], word, firstEquation.LeftPart.Symbols[0]}
			}
			substFirst := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], newValsFirst)
			firstChild := newEquationSystemWithSubstitution(node, &substFirst, node.number+"1")

			var newValsSecond []symbol.Symbol
			if solver.algorithmType == STANDARD {
				newValsSecond = []symbol.Symbol{firstEquation.LeftPart.Symbols[0], firstEquation.RightPart.Symbols[0]}
			}
			if solver.algorithmType == FINITE {
				word := solver.getLetter()
				newValsSecond = []symbol.Symbol{firstEquation.LeftPart.Symbols[0], word, firstEquation.RightPart.Symbols[0]}
			}

			substSecond := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], newValsSecond)
			secondChild := newEquationSystemWithSubstitution(node, &substSecond, node.number+"2")

			substThird := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{firstEquation.RightPart.Symbols[0]})
			thirdChild := newEquationSystemWithSubstitution(node, &substThird, node.number+"3")

			node.SetChildren([]*Node{&thirdChild, &firstChild, &secondChild})
		}

		if checkSecondRuleLeft(firstEquation) {
			substFirst := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
			firstChild := newEquationSystemWithSubstitution(node, &substFirst, node.number+"4")

			substSecond := equation.NewSubstitution(firstEquation.RightPart.Symbols[0], []symbol.Symbol{firstEquation.LeftPart.Symbols[0], firstEquation.RightPart.Symbols[0]})
			secondChild := newEquationSystemWithSubstitution(node, &substSecond, node.number+"5")

			node.SetChildren([]*Node{&firstChild, &secondChild})
		}
		if checkSecondRuleRight(firstEquation) {
			substFirst := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
			firstChild := newEquationSystemWithSubstitution(node, &substFirst, node.number+"6")

			substSecond := equation.NewSubstitution(firstEquation.LeftPart.Symbols[0], []symbol.Symbol{firstEquation.RightPart.Symbols[0], firstEquation.LeftPart.Symbols[0]})
			secondChild := newEquationSystemWithSubstitution(node, &substSecond, node.number+"7")

			node.SetChildren([]*Node{&firstChild, &secondChild})
		}
		if checkThirdRuleLeft(firstEquation) || checkThirdRuleRight(firstEquation) {
			newES, subsVars := node.value.SubstituteVarsWithEmpty()
			var childL *Node
			for v, _ := range subsVars {
				if childL != nil {
					node = childL
				}
				subst := equation.NewSubstitution(v, []symbol.Symbol{symbol.Empty()})
				// Writing original equation for every node
				child := NewNodeWEquationsSystem(subst, node.number+"8", node, node.value)
				node.SetChildren([]*Node{&child})
				childL = &child
			}
			// Writing equation with all substituted vars
			node.children[0].value = newES

		}
	}

	//node.Print()
	//for i, child := range node.children {
	//	fmt.Printf(" %d  :", i)
	//	child.Print()
	//}
	if len(node.children) == 0 {
		solver.createFalseNode(node, REGULAR_FALSE)
		return nil
	}
	for _, child := range node.children {
		err = solver.solveSystem(child)
		if err != nil {
			return fmt.Errorf("error solving for child: %v", err)
		}
	}
	node.FillHelpMapFromChildren()
	node.FillSubstituteMapsFromChildren()
	// needs to go afterwards because of the map filling
	for _, child := range node.children {
		if child.HasTrueChildren() {
			node.AddSubstituteVar(child.substitution.LeftPart())
		}
		if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
			if child.HasTrueChildren() {
				node.letterSubstitutions = append(node.letterSubstitutions, child.letterSubstitutions...)
			}
		}
	}

	if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
		if node.Substitution().IsToLetter() {
			for _, ls := range node.letterSubstitutions {
				ls.AddSubstToHead(node.substitution)
			}
		}
	}

	return nil
}

func (solver *Solver) solve(node *Node) error {
	var err error
	if !solver.solveOptions.FullGraph && solver.hasSolution {
		return nil
	}
	if len(node.number) > solver.solveOptions.CycleRange {
		solver.cycled = true
		return nil
	}

	// checking length

	if solver.solveOptions.LengthAnalysis {
		checkedLength, replaceSymbol, replaceLen := checkLengthRules(node.value.Equation())
		if checkedLength {
			if replaceSymbol != nil {
				var newLetters []symbol.Symbol
				for i := 0; i < replaceLen; i++ {
					newLetters = append(newLetters, solver.getLetter())
				}
				substitute := equation.NewSubstitution(replaceSymbol, newLetters)
				eq := node.value.Equation().Substitute(substitute)
				if solver.algorithmType == STANDARD && !(eq.RightPart.IsEmpty() || eq.LeftPart.IsEmpty()) {

				} else {
					child := NewNodeWEquation(substitute, "r"+node.number, node, eq)
					node.SetChildren([]*Node{&child})
					node = &child
				}
			}
		} else {
			solver.createFalseNode(node, FAILED_LENGTH_ANALISYS)
			return nil
		}
	}

	//fmt.Println(node.number)
	if checkInequality(node) {
		solver.createFalseNode(node, REGULAR_FALSE)
		return nil
	}
	if checkEquality(node) {
		solver.createTrueNode(node)
		if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
			nls := NewLetterSubstitution(node)
			if node.Substitution().IsToLetter() {
				nls.AddSubstToHead(node.substitution)
			}
			node.letterSubstitutions = append(node.letterSubstitutions, &nls)
		}
		return nil
	}

	hasBeen, tr := checkHasBeen(node)
	if hasBeen {
		node.SetHasBackCycle()
		node.SetChildren([]*Node{tr})
		tr.AddParentFromBackCycle(node)
		return nil
	}

	if solver.algorithmType == FINITE {
		if checkFirstRuleFinite(node.value.Equation()) {
			substitute := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{node.value.Equation().RightPart.Symbols[0]})
			eq := node.value.Equation().Substitute(substitute)
			child := NewNodeWEquation(substitute, "a"+node.number, node, eq)
			node.SetChildren([]*Node{&child})
		}
		if checkSecondRuleLeftFinite(node.value.Equation()) {
			substitute := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], []symbol.Symbol{node.value.Equation().LeftPart.Symbols[0]})
			eq := node.value.Equation().Substitute(substitute)
			child := NewNodeWEquation(substitute, "b"+node.number, node, eq)
			node.SetChildren([]*Node{&child})
		}
		if checkSecondRuleRightFinite(node.value.Equation()) {
			substitute := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{node.value.Equation().RightPart.Symbols[0]})
			eq := node.value.Equation().Substitute(substitute)
			child := NewNodeWEquation(substitute, "c"+node.number, node, eq)

			node.SetChildren([]*Node{&child})
		}
		if checkFourthRuleLeft(node.value.Equation()) {
			substFirst := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
			firstEquation := node.value.Equation().Substitute(substFirst)
			firstChild := NewNodeWEquation(substFirst, "d"+node.number, node, firstEquation)

			substSecond := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], []symbol.Symbol{node.value.Equation().LeftPart.Symbols[0], node.value.Equation().RightPart.Symbols[0]})

			secondEquation := node.value.Equation().Substitute(substSecond)
			secondChild := NewNodeWEquation(substSecond, "e"+node.number, node, secondEquation)

			node.SetChildren([]*Node{&firstChild, &secondChild})
		}
		if checkFourthRuleRight(node.value.Equation()) {
			substFirst := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{symbol.Empty()})

			firstEquation := node.value.Equation().Substitute(substFirst)
			firstChild := NewNodeWEquation(substFirst, "f"+node.number, node, firstEquation)

			substSecond := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{node.value.Equation().RightPart.Symbols[0], node.value.Equation().LeftPart.Symbols[0]})

			secondEquation := node.value.Equation().Substitute(substSecond)
			secondChild := NewNodeWEquation(substSecond, "g"+node.number, node, secondEquation)

			node.SetChildren([]*Node{&firstChild, &secondChild})
		}
	}
	if checkFirstRule(node.value.Equation()) {
		var newValsFirst []symbol.Symbol
		if solver.algorithmType == STANDARD {
			newValsFirst = []symbol.Symbol{node.value.Equation().RightPart.Symbols[0], node.value.Equation().LeftPart.Symbols[0]}
		}
		if solver.algorithmType == FINITE {
			word := solver.getLetter()
			newValsFirst = []symbol.Symbol{node.value.Equation().RightPart.Symbols[0], word, node.value.Equation().LeftPart.Symbols[0]}
		}
		substFirst := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], newValsFirst)

		firstEquation := node.value.Equation().Substitute(substFirst)
		firstChild := NewNodeWEquation(substFirst, node.number+"1", node, firstEquation)

		var newValsSecond []symbol.Symbol
		if solver.algorithmType == STANDARD {
			newValsSecond = []symbol.Symbol{node.value.Equation().LeftPart.Symbols[0], node.value.Equation().RightPart.Symbols[0]}
		}
		if solver.algorithmType == FINITE {
			word := solver.getLetter()
			newValsSecond = []symbol.Symbol{node.value.Equation().LeftPart.Symbols[0], word, node.value.Equation().RightPart.Symbols[0]}
		}

		substSecond := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], newValsSecond)
		secondEquation := node.value.Equation().Substitute(substSecond)
		secondChild := NewNodeWEquation(substSecond, node.number+"2", node, secondEquation)

		substThird := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{node.value.Equation().RightPart.Symbols[0]})
		thirdEquation := node.value.Equation().Substitute(substThird)
		thirdChild := NewNodeWEquation(substThird, node.number+"3", node, thirdEquation)
		node.SetChildren([]*Node{&thirdChild, &firstChild, &secondChild})
	}

	if checkSecondRuleLeft(node.value.Equation()) {
		substFirst := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
		firstEquation := node.value.Equation().Substitute(substFirst)
		firstChild := NewNodeWEquation(substFirst, node.number+"4", node, firstEquation)

		substSecond := equation.NewSubstitution(node.value.Equation().RightPart.Symbols[0], []symbol.Symbol{node.value.Equation().LeftPart.Symbols[0], node.value.Equation().RightPart.Symbols[0]})

		secondEquation := node.value.Equation().Substitute(substSecond)
		secondChild := NewNodeWEquation(substSecond, node.number+"5", node, secondEquation)
		node.SetChildren([]*Node{&firstChild, &secondChild})
	}
	if checkSecondRuleRight(node.value.Equation()) {
		substFirst := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{symbol.Empty()})
		firstEquation := node.value.Equation().Substitute(substFirst)
		firstChild := NewNodeWEquation(substFirst, node.number+"6", node, firstEquation)

		substSecond := equation.NewSubstitution(node.value.Equation().LeftPart.Symbols[0], []symbol.Symbol{node.value.Equation().RightPart.Symbols[0], node.value.Equation().LeftPart.Symbols[0]})

		secondEquation := node.value.Equation().Substitute(substSecond)
		secondChild := NewNodeWEquation(substSecond, node.number+"7", node, secondEquation)
		node.SetChildren([]*Node{&firstChild, &secondChild})
	}
	if checkThirdRuleLeft(node.value.Equation()) || checkThirdRuleRight(node.value.Equation()) {
		eq, subsVars := node.value.Equation().SubstituteVarsWithEmpty()
		var child Node
		for v, _ := range subsVars {
			if child.number != "" {
				*node = child
			}
			subst := equation.NewSubstitution(v, []symbol.Symbol{symbol.Empty()})
			// Writing original equation for every node
			child = NewNodeWEquationsSystem(subst, node.number+"8", node, node.value)
			node.SetChildren([]*Node{&child})
		}
		// Writing equation with all substituted vars
		node.children[0].value = equation.NewSingleEquation(eq)

	}
	//node.Print()
	//for i, child := range node.children {
	//	fmt.Printf(" %d  :", i)
	//	child.Print()
	//}
	if len(node.children) == 0 {
		solver.createFalseNode(node, REGULAR_FALSE)
		return nil
	}
	for _, child := range node.children {
		err = solver.solve(child)
		if err != nil {
			return fmt.Errorf("error solving for child: %v", err)
		}
	}
	node.FillHelpMapFromChildren()
	node.FillSubstituteMapsFromChildren()
	// needs to go afterwards because of the map filling
	for _, child := range node.children {
		if child.HasTrueChildren() {
			node.AddSubstituteVar(child.substitution.LeftPart())
		}
		if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
			if child.HasTrueChildren() {
				node.letterSubstitutions = append(node.letterSubstitutions, child.letterSubstitutions...)
			}
		}
	}

	if solver.algorithmType == FINITE && solver.solveOptions.SaveLettersSubstitutions {
		if node.Substitution().IsToLetter() {
			for _, ls := range node.letterSubstitutions {
				ls.AddSubstToHead(node.substitution)
			}
		}
	}

	return nil
}

func (solver *Solver) describeGraph(node *Node) error {
	var err error
	err = solver.dotWriter.WriteNode(node)
	if err != nil {
		return fmt.Errorf("error writing node: %v", err)
	}

	if len(node.children) == 0 && node.HasInfoChild() {
		infoChild := node.InfoChild()
		err = solver.dotWriter.WriteInfoNode(infoChild)
		if err != nil {
			return fmt.Errorf("error writing info node: %v", err)
		}
		if IsFalseNode(infoChild) {
			err = solver.dotWriter.WriteInfoEdgeWithLabel(node, (infoChild).(*FalseNode))
		} else {
			err = solver.dotWriter.WriteInfoEdge(node, infoChild)
		}
		if err != nil {
			return fmt.Errorf("error writing info edge: %v", err)
		}
	}

	if node.LeadsToBackCycle() {
		err = solver.dotWriter.WriteDottedEdge(node, node.children[0])
		if err != nil {
			return fmt.Errorf("error writing dotted edge: %v", err)
		}
		return nil
	}

	for _, child := range node.children {
		subst := child.Substitution()
		if subst.IsEmpty() {
			err = solver.dotWriter.WriteLabelEdgeBold(node, child)
			if err != nil {
				return fmt.Errorf("error writing label edge bold: %v", err)
			}
		} else {
			leftSym := subst.LeftPart()
			rightPart := subst.RightPart()
			err = solver.dotWriter.WriteLabelEdge(node, child, &leftSym, rightPart)
			if err != nil {
				return fmt.Errorf("error writing label edge: %v", err)
			}
		}

		err = solver.describeGraph(child)
		if err != nil {
			return fmt.Errorf("error solving for child: %v", err)
		}
	}
	return nil
}
