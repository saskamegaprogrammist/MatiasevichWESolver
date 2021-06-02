package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Simplifier struct {
	solver *Solver
}

func (s *Simplifier) Init(constantsAlph string, varsAlph string,
	printOptions PrintOptions, solveOptions SolveOptions) error {
	var err error
	solver := Solver{}
	s.solver = &solver
	err = s.solver.InitWoEquation(constantsAlph, varsAlph, printOptions, solveOptions)
	if err != nil {
		return fmt.Errorf("error initing solver: %v", err)
	}
	return err
}

func (s *Simplifier) Simplify(node *Node) error {
	var err error
	if !node.WasUnfolded() {
		s.solver.clear()
		err = s.solver.solveSystem(node)
		if err != nil {
			return fmt.Errorf("error solving node: %v", err)
		}
	}

	if len(node.value.Equation().Letters()) != 0 {
		var conjunctions = make([]equation.EquationsSystem, 0)
		var resSystem, nodeSystem equation.EquationsSystem
		newNodes, eqs, err := s.checkRulesForLetters(node)
		if err != nil {
			return fmt.Errorf("error checking rules for letters letters: %v", err)
		}
		for i, newNode := range newNodes {
			nodeSystem, err = s.simplifyNode(newNode)
			if err != nil {
				return fmt.Errorf("error simplifying node with letters: %v", err)
			}
			newNode.simplified = nodeSystem
			//nodeSystem.Print()

			conjunctions = append(conjunctions, equation.NewConjunctionSystemFromEquations(append(eqs[i], nodeSystem.GetEquations()...)))
		}
		resSystem = equation.NewDisjunctionSystem(conjunctions)
		node.simplified = resSystem
		resSystem.Print()
	} else {
		var resSystem equation.EquationsSystem
		resSystem, err = s.simplifyNode(node)
		if err != nil {
			return fmt.Errorf("error simplifying node without letters: %v", err)
		}
		node.simplified = resSystem
		resSystem.Print()
	}
	return nil
}

func (s *Simplifier) simplifyNode(node *Node) (equation.EquationsSystem, error) {
	var err error
	var resSystem equation.EquationsSystem
	symbols := standart.SymbolArrayFromIntMap(node.subgraphsSubstituteVars)
	if len(symbols) == 0 {
		return equation.NewSingleEquation(*(node.value.Equation())), nil
	}
	subgraphSymbol := symbols[len(symbols)-1]
	resSystem, err = s.simplify(node, subgraphSymbol)
	if err != nil {
		return resSystem, fmt.Errorf("error during simplification: %v", err)
	}
	return resSystem, nil
}

func (s *Simplifier) checkTrueNodesWoLetterSubstitutions(node *Node) (bool, []LetterSubstitution) {
	var nodesToEmpty = make([]*Node, 0)
	var nonEmpty = make([]LetterSubstitution, 0)
	letterSubstitutions := node.LetterSubstitutions()
	for _, ls := range letterSubstitutions {
		if ls.HasNoSubstitutions() {
			nodesToEmpty = append(nodesToEmpty, ls.nodeToTrue)
		} else {
			nonEmpty = append(nonEmpty, *ls)
		}
	}
	if len(nodesToEmpty) > 0 {
		for _, n := range nonEmpty {
			s.removeTrueNodesWLetters(n.nodeToTrue)
		}
		return true, nonEmpty
	}
	return false, nonEmpty
}

func (s *Simplifier) checkRulesForLetters(node *Node) ([]*Node, [][]equation.Equation, error) {
	var err error
	letters := node.value.Equation().Letters()
	var newLetterSubstitutions []LetterSubstitution
	hasEmpty, nonEmpty := s.checkTrueNodesWoLetterSubstitutions(node)

	if hasEmpty {
		return []*Node{node}, [][]equation.Equation{{}}, nil
	}

	newLetterSubstitutions = s.checkEqualLetterSubstituions(nonEmpty, letters)

	var newNodes []*Node
	var eqs = make([][]equation.Equation, 0)
	for _, nls := range newLetterSubstitutions {
		newEq := nls.NewEquation(node.value.Equation())
		newNode := NewTreeWEquation("0", newEq)
		err = s.solver.solveSystem(&newNode)
		if err != nil {
			return newNodes, eqs, fmt.Errorf("error solving node: %v", err)
		}
		s.checkTrueNodesWoLetterSubstitutions(&newNode)
		newNodes = append(newNodes, &newNode)
		eqs = append(eqs, nls.SubstitutionsAsEquations())
	}
	return newNodes, eqs, nil
}

func (s *Simplifier) checkEqualLetterSubstituions(l []LetterSubstitution, letters []symbol.Symbol) []LetterSubstitution {
	var lsMaps = make([]map[symbol.Symbol]symbol.Symbol, 0)
	var newL = make([]LetterSubstitution, 0)
	for _, l := range l {
		lsMaps = append(lsMaps, s.createMapWithLetters(l, letters))
	}
	s.compareLetterSubstitutionMaps(lsMaps, &l, &newL)
	return newL
}

func (s *Simplifier) compareLetterSubstitutionMaps(maps []map[symbol.Symbol]symbol.Symbol, oldL *[]LetterSubstitution,
	newL *[]LetterSubstitution) {
	var newLMap = make(map[int]bool)
	var mapsLen = len(maps)
	var needsDeletion bool
	var linkToBigger *map[symbol.Symbol]symbol.Symbol
	for i, m := range maps {
		for j := i + 1; j < mapsLen; j++ {
			needsDeletion, linkToBigger = s.compareLetterMaps(m, maps[j])
			if !needsDeletion {
				newLMap[i] = true
				newLMap[j] = true
			} else {
				if linkToBigger == &maps[j] {
					newLMap[i] = false
					newLMap[j] = true
				} else {
					newLMap[i] = true
					newLMap[j] = false
				}
			}
		}
	}
	for i, v := range newLMap {
		if v == true {
			*newL = append(*newL, (*oldL)[i])
		}
	}
}

func (s *Simplifier) compareLetterMaps(fMap map[symbol.Symbol]symbol.Symbol,
	sMap map[symbol.Symbol]symbol.Symbol) (bool, *map[symbol.Symbol]symbol.Symbol) {
	var biggerSubst *map[symbol.Symbol]symbol.Symbol
	for symbolKey, symbolVal := range fMap {
		sVal := sMap[symbolKey]
		if sVal == symbolVal {
			continue
		}
		if symbol.IsConst(symbolVal) {
			if symbol.IsConst(sVal) {
				return false, nil
			}
			if symbol.IsLetter(sVal) {
				if biggerSubst != nil && biggerSubst == &fMap {
					return false, nil
				}
				biggerSubst = &sMap
			}
		} else {
			if symbol.IsConst(sVal) {
				return false, nil
			}
			if symbol.IsLetter(sVal) {
				if biggerSubst != nil && biggerSubst == &sMap {
					return false, nil
				}
				biggerSubst = &fMap
			}
		}
	}
	// two maps are equal, substitutions are equal, choosing any
	if biggerSubst == nil {
		biggerSubst = &fMap
	}
	return true, biggerSubst
}

func (s *Simplifier) unfoldMap(lMap *map[symbol.Symbol]symbol.Symbol, letters []symbol.Symbol) {
	for _, letter := range letters {
		val := (*lMap)[letter]
		if val == nil {
			(*lMap)[letter] = letter
		} else {
			(*lMap)[letter] = unfoldSymbol(val, lMap)
		}
	}
}

func unfoldSymbol(sym symbol.Symbol, lMap *map[symbol.Symbol]symbol.Symbol) symbol.Symbol {
	if symbol.IsConst(sym) {
		return sym
	}
	if symbol.IsLetter(sym) {
		val := (*lMap)[sym]
		if val == nil || val == sym {
			return sym
		}
		(*lMap)[sym] = unfoldSymbol(val, lMap)
	}
	return sym
}

func (s *Simplifier) createMapWithLetters(ls LetterSubstitution, letters []symbol.Symbol) map[symbol.Symbol]symbol.Symbol {
	var lsMap = make(map[symbol.Symbol]symbol.Symbol)
	for _, s := range ls.substitutions {
		lsMap[s.LeftPart()] = s.RightPart()[0]
	}
	s.unfoldMap(&lsMap, letters)
	return lsMap
}

func (s *Simplifier) removeTrueNodesWLetters(node *Node) {
	node.infoChild = nil
	node.SetDoesntHaveTrueChildren()
	tr := node.parent

	for tr != nil {
		tr.FillHelpMapFromChildren()
		tr = tr.parent
	}
}

func (s *Simplifier) simplify(node *Node, symbolVar symbol.Symbol) (equation.EquationsSystem, error) {
	var resultEquationSystem equation.EquationsSystem
	var err error
	var eqSystems []equation.EquationsSystem
	err = s.walk(node, &eqSystems, symbolVar)
	if err != nil {
		return equation.EquationsSystem{}, fmt.Errorf("error walking node: %v", err)
	}
	if len(eqSystems) == 0 {
		return resultEquationSystem, nil
	}
	disjunctions := getAllDisjunctions(eqSystems)
	newGraphs := make([]Node, 0)
	var varMap = make(map[symbol.Symbol]bool)
	for _, disj := range disjunctions {
		disj.Print()
		newGraph := Node{}
		err = copyGraph(node, &newGraph, disj, symbolVar)
		if err != nil {
			return equation.EquationsSystem{}, fmt.Errorf("error copying graph: %v", err)
		}
		newGraphs = append(newGraphs, newGraph)
		standart.MergeMapsInt(&varMap, newGraph.subgraphsSubstituteVars)
	}
	symbols := standart.SymbolArrayFromBoolMap(varMap)
	if len(symbols) == 0 {
		return equation.NewDisjunctionSystem(disjunctions), nil
	}
	subgraphSymbol := symbols[len(symbols)-1]
	var newDisjunctions []equation.EquationsSystem
	for i, graph := range newGraphs {
		var es, newEs equation.EquationsSystem
		es, err = s.simplify(&graph, subgraphSymbol)
		if err != nil {
			return equation.EquationsSystem{}, fmt.Errorf("error simplifying children graph: %v", err)
		}
		newEs = equation.NewConjunctionSystem([]equation.EquationsSystem{es, disjunctions[i]})
		newDisjunctions = append(newDisjunctions, newEs)
	}
	return equation.NewDisjunctionSystem(newDisjunctions), nil
}

func copyGraph(node *Node, copyNode *Node, disjunctionComponent equation.EquationsSystem, currSymbol symbol.Symbol) error {
	var err error
	copyNode.Copy(node)
	if node.simplified.HasEqSystem(disjunctionComponent) {
		copyNode.SetHasTrueChildren()
		trueNode := &TrueNode{
			number: "T_" + copyNode.number,
		}
		copyNode.infoChild = trueNode
		//var sym symbol.Symbol
		//sym, err = node.SubstituteVar()
		//if err != nil {
		//	return fmt.Errorf("error getting substitute var: %v", err)
		//}
		//copyNode.RemoveSubstituteVar(sym, len(node.children))
		return nil
	}
	if len(node.children) == 0 {
		return nil
	}
	if node.LeadsToBackCycle() {
		tr := copyNode.parent
		for tr != nil {
			if node.value.Equation().CheckSameness(tr.value.Equation()) {
				break
			}
			tr = tr.parent
		}
		copyNode.children = append(copyNode.children, tr)
		copyNode.SetHasBackCycle()
		return nil
	}
	for _, ch := range node.children {
		if ch.HasSingleSubstituteVar() {
			var sym symbol.Symbol
			sym, err = ch.SubstituteVar()
			if err != nil {
				return fmt.Errorf("error getting substitute var: %v", err)
			}
			if sym.Value() == currSymbol.Value() && !ch.simplified.HasEqSystem(disjunctionComponent) {
				continue
			}
		}
		newChildNode := Node{}
		newChildNode.parent = copyNode
		err = copyGraph(ch, &newChildNode, disjunctionComponent, currSymbol)
		if err != nil {
			return fmt.Errorf("error copying child graph: %v", err)
		}
		if newChildNode.HasTrueChildrenAndBackCycles() {
			copyNode.AddSubstituteVar(newChildNode.substitution.LeftPart())
			copyNode.children = append(copyNode.children, &newChildNode)
		}
	}
	copyNode.FillSubstituteMapsFromChildren()
	copyNode.FillHelpMapFromChildren()
	return nil
}

func getAllDisjunctions(eqSystems []equation.EquationsSystem) []equation.EquationsSystem {
	allDisjunctions := make([]equation.EquationsSystem, 0)
	for _, eqSys := range eqSystems {
		if eqSys.IsSingleEquation() {
			allDisjunctions = append(allDisjunctions, eqSys)
		} else {
			allDisjunctions = append(allDisjunctions, getAllDisjunctions(eqSys.Compounds())...)
		}
	}
	return allDisjunctions
}

func (s *Simplifier) walk(node *Node, eqSystems *[]equation.EquationsSystem, subgraphSymbol symbol.Symbol) error {
	var err error
	if node.HasFalseChildrenAndBackCycles() || node.LeadsToBackCycle() {
		return nil
	}
	if node.HasSingleSubstituteVar() {
		var sVar symbol.Symbol
		sVar, err = node.SubstituteVar()
		if err != nil {
			return fmt.Errorf("error getting substitution var: %v", err)
		}
		if sVar != subgraphSymbol {
			return nil
		}
		if !node.HasBackCycle() {
			values := s.walkWithSymbol(node)
			values.ReduceEmptyVars()

			es := equation.SystemFromValues(sVar, values)
			node.SetSimplifiedRepresentation(es)
			node.simplified.Print()

			*eqSystems = append(*eqSystems, es)
			return nil
		} else {
			node.SetIsSubgraphRoot()

			values, _, _, needsReduce := s.walkWithSymbolBackCycled(node)
			var valuesLen = len(values)
			var es equation.EquationsSystem
			var eqType int = equation.EQ_TYPE_SIMPLE
			var valuesToEq = values[valuesLen-1]
			for _, v := range values {
				v.ReduceEmptyVars()
			}
			if needsReduce && valuesLen > 1 {
				if valuesLen != 2 {
					return fmt.Errorf("values length must be 2")
				}
				if !standart.CheckSymbolArraysEquality(values[0][1], values[1][1]) {
					return fmt.Errorf("can't create characteristic equation")
				}
				if len(values[0][0]) == 1 && symbol.IsEmpty(values[0][0][0]) {
					eqType = equation.EQ_TYPE_W_EMPTY
				} else if !standart.CheckSymbolArraysEquality(values[0][0], values[1][0]) {
					return fmt.Errorf("can't create characteristic equation")
				}
			}
			es, err = equation.CharacteristicEquation(sVar, valuesToEq, eqType)
			if err != nil {
				return fmt.Errorf("can't create characteristic equation: %v", err)

			}
			node.SetSimplifiedRepresentation(es)
			node.simplified.Print()

			node.UnsetIsSubgraphRoot()
			*eqSystems = append(*eqSystems, es)
			return nil
		}

	}
	for _, ch := range node.Children() {
		err = s.walk(ch, eqSystems, subgraphSymbol)
		if err != nil {
			return fmt.Errorf("error walking child: %v", err)
		}
	}
	return nil
}

func (s *Simplifier) walkWithSymbol(node *Node) equation.VariableValues {
	var values = equation.NewVariableValues()
	if node.LeadsToBackCycle() {
		return values
	}
	var chValues equation.VariableValues
	for _, ch := range node.Children() {
		if ch.HasOnlyFalseChildren() {
			continue
		}
		chValues = s.walkWithSymbol(ch)
		if chValues.IsEmpty() {
			values.AddValue([]symbol.Symbol{ch.Substitution().RightPart()[0]})
		} else {
			values.AddToEachValue(chValues, []symbol.Symbol{ch.Substitution().RightPart()[0]})
		}
	}
	return values
}

func (s *Simplifier) walkWithSymbolBackCycled(node *Node) ([]equation.VariableValues, *Node, bool, bool) {
	var values = equation.NewVariableValuesArray()
	values = append(values, equation.NewVariableValues())
	if node.LeadsToBackCycle() {
		// returning node to which it has back cycle to
		return values, node.children[0], false, false
	}
	//var chValues []equation.VariableValues
	var parentNode *Node
	var metEmptySubstNode bool
	var index int
	var filteredChildren = make([][]equation.VariableValues, 0)
	var currParentNodes = make([]*Node, 0)
	var metEmptyBefore = make([]bool, 0)
	var hasOneCharValues = make([]bool, 0)
	var newLetters = make([]symbol.Symbol, 0)
	var size = 0

	for _, ch := range node.Children() {
		if ch.HasOnlyFalseChildren() {
			continue
		}
		size++
		chValues, currParentNode, metEmptySubstNodebefore, hasOneChar := s.walkWithSymbolBackCycled(ch)
		filteredChildren = append(filteredChildren, chValues)
		currParentNodes = append(currParentNodes, currParentNode)
		metEmptyBefore = append(metEmptyBefore, metEmptySubstNodebefore)
		hasOneCharValues = append(hasOneCharValues, hasOneChar)
		newLetters = append(newLetters, ch.NewLetter())

		if currParentNode != nil {
			parentNode = currParentNode
		}
	}
	for _, v := range hasOneCharValues {
		if v {
			index = 1
			values = append(values, equation.NewVariableValues())
			break
		}
	}
	for i := 0; i < size; i++ {

		if filteredChildren[i][index].IsEmpty() {
			// means it was child who leads to true node with empty substitution
			if currParentNodes[i] == nil {
				// adding empty symbol
				metEmptySubstNode = true
				values[index].AddValueToHead([]symbol.Symbol{newLetters[i]})
			} else {
				// means it was child who leads to cycle head
				values[index].AddValue([]symbol.Symbol{newLetters[i]})
			}

		} else {
			if metEmptyBefore[i] {
				metEmptySubstNode = true
				values[index].AddToFirstValue(filteredChildren[i][index], []symbol.Symbol{newLetters[i]})
			} else {
				values[index].AddToSecondValue(filteredChildren[i][index], []symbol.Symbol{newLetters[i]})
			}
		}
	}
	if index == 1 || (parentNode != nil && parentNode.number == node.number) {
		return values, nil, metEmptySubstNode, true
	}
	return values, nil, metEmptySubstNode, false
}
