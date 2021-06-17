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
		if !node.HasTrueChildren() {
			node.simplified = node.value
			return nil
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
		//resSystem.Print()
	} else {
		var resSystem equation.EquationsSystem
		resSystem, err = s.simplifyNode(node)
		if err != nil {
			return fmt.Errorf("error simplifying node without letters: %v", err)
		}
		node.simplified = resSystem
		//resSystem.Print()
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
	subgraphSymbol := symbols[0]
	resSystem, err = s.simplify(node, subgraphSymbol, []symbol.Symbol{subgraphSymbol})
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

func (s *Simplifier) simplify(node *Node, symbolVar symbol.Symbol, hasAlreadyBeen []symbol.Symbol) (equation.EquationsSystem, error) {
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
	ds := equation.NewDisjunctionSystem(disjunctions)
	ds.Simplify()
	ds.Reduce()
	newGraphs := make([]Node, 0)
	var varMap = make(map[symbol.Symbol]bool)
	for _, disj := range ds.Compounds() {
		disj.Print()
		fmt.Println()
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
	var subgraphSymbol symbol.Symbol
	for _, s := range symbols {
		for _, h := range hasAlreadyBeen {
			if s == h {
				break
			}
		}
		subgraphSymbol = s
		hasAlreadyBeen = append(hasAlreadyBeen, s)
		break
	}
	var newDisjunctions []equation.EquationsSystem
	if !symbol.IsLetterOrVar(subgraphSymbol) {
		return equation.NewDisjunctionSystem(newDisjunctions), nil
	}
	for i, graph := range newGraphs {
		var es, newEs equation.EquationsSystem
		es, err = s.simplify(&graph, subgraphSymbol, hasAlreadyBeen)
		if err != nil {
			return equation.EquationsSystem{}, fmt.Errorf("error simplifying children graph: %v", err)
		}
		newEs = equation.NewConjunctionSystem([]equation.EquationsSystem{disjunctions[i], es})
		newDisjunctions = append(newDisjunctions, newEs)
	}
	return equation.NewDisjunctionSystem(newDisjunctions), nil
}

func copyGraph(node *Node, copyNode *Node, disjunctionComponent equation.EquationsSystem, currSymbol symbol.Symbol) error {
	var err error
	copyNode.Copy(node)
	// TODO: возможно надо создавать узел False явно
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
			if sym != nil && sym.Value() == currSymbol.Value() && !ch.simplified.HasEqSystem(disjunctionComponent) {
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
			if !newChildNode.Substitution().IsEmpty() &&
				// new rule
				newChildNode.substitution.LeftPart() != currSymbol {
				copyNode.AddSubstituteVar(newChildNode.substitution.LeftPart())
			}
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
		//node.Print()

		nodeToWalk := node
		for len(nodeToWalk.children) == 1 && nodeToWalk.children[0].substitution.IsEmpty() {
			nodeToWalk = node.children[0]
		}

		if !nodeToWalk.HasBackCycle() {
			values := s.walkWithSymbol(nodeToWalk)
			values.ReduceEmptyVars()

			es := equation.SystemFromValues(sVar, values)
			node.SetSimplifiedRepresentation(es)
			//node.simplified.Print()

			*eqSystems = append(*eqSystems, es)
			return nil
		} else {

			nodeToWalk.SetIsSubgraphRoot()

			f1, f2, _, _ := s.walkWithSymbolBackCycledRefactored(nodeToWalk, sVar)
			fmt.Println(f1)
			fmt.Println(f2)

			var es equation.EquationsSystem
			var eqType int = equation.EQ_TYPE_SIMPLE

			if f2.IsEmpty() {
				eqType = equation.EQ_TYPE_SIMPLE
			} else {
				f1E, f2E := f1.Compare(f2)
				if f2E {
					if f1E {
						eqType = equation.EQ_TYPE_SIMPLE
					} else {
						eqType = equation.EQ_TYPE_W_EMPTY
					}
				} else {
					return fmt.Errorf("can't create characteristic equation")
				}
			}

			es, err = equation.CharacteristicEquationRefactored(sVar, f1, eqType)

			if err != nil {
				return fmt.Errorf("can't create characteristic equation: %v", err)

			}
			//es.Print()
			node.SetSimplifiedRepresentation(es)
			//node.simplified.Print()

			nodeToWalk.UnsetIsSubgraphRoot()
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
		if !ch.Substitution().IsEmpty() {
			if chValues.IsEmpty() {
				values.AddValue([]symbol.Symbol{ch.Substitution().RightPart()[0]})
			} else {
				values.AddToEachValue(chValues, []symbol.Symbol{ch.Substitution().RightPart()[0]})
			}
		}
	}
	return values
}

func (s *Simplifier) walkWithSymbolBackCycledRefactored(node *Node, sVar symbol.Symbol) (equation.CharacteristicValues,
	equation.CharacteristicValues, *Node, bool) {
	if node.LeadsToBackCycle() {
		// returning node to which it has back cycle to
		if !node.Substitution().IsEmpty() {
			return equation.NewCharacteristicValuesFromArrays([]symbol.Symbol{}, []symbol.Symbol{node.NewLetter()}),
				equation.NewCharacteristicValues(), node.children[0], false
		} else {
			return equation.NewCharacteristicValues(), equation.NewCharacteristicValues(), node.children[0], false
		}
	}
	if len(node.children) == 0 && node.infoChild != nil {
		if !node.Substitution().IsEmpty() {
			return equation.NewCharacteristicValuesFromArrays([]symbol.Symbol{node.NewLetter()}, []symbol.Symbol{}),
				equation.NewCharacteristicValues(), nil, true
		} else {
			return equation.NewCharacteristicValues(), equation.NewCharacteristicValues(), nil, true
		}
	}
	var cycleToNode *Node
	var chValuesMainArray = equation.NewCharacteristicValuesArray()
	var chValuesHelpArray = equation.NewCharacteristicValuesArray()
	var metTrueNodeAll bool
	for _, ch := range node.Children() {
		if !ch.substitution.IsEmpty() && ch.substitution.IsTo() != sVar.Value() {
			continue
		}
		if ch.HasOnlyFalseChildren() {
			continue
		}
		chValuesMain, chValuesHelp, currParentNode, metTrueNode := s.walkWithSymbolBackCycledRefactored(ch, sVar)
		//ch.Print()
		//fmt.Println(chValuesMain)
		//fmt.Println(chValuesHelp)
		chValuesMainArray = append(chValuesMainArray, chValuesMain)
		chValuesHelpArray = append(chValuesHelpArray, chValuesHelp)
		if currParentNode != nil {
			cycleToNode = currParentNode
		}
		if metTrueNode {
			metTrueNodeAll = true
		}
	}
	var newChvalues = equation.NewCharacteristicValues()
	var newChvaluesHelp = equation.NewCharacteristicValues()
	for _, v := range chValuesMainArray {
		newChvalues.Append(v)
	}
	for _, v := range chValuesHelpArray {
		newChvaluesHelp.Append(v)
	}
	var newLetter symbol.Symbol
	if !node.substitution.IsEmpty() {
		newLetter = node.NewLetter()
	}
	if cycleToNode == node {
		if node.IsSubgraphRoot() {
			return newChvalues, equation.NewCharacteristicValues(), nil, false
		}
		if !node.substitution.IsEmpty() {
			return equation.NewCharacteristicValuesFromArrays([]symbol.Symbol{newLetter}, []symbol.Symbol{}),
				newChvalues, nil, false
		}
		return equation.NewCharacteristicValuesFromArrays([]symbol.Symbol{}, []symbol.Symbol{}),
			newChvalues, nil, false
	} else {
		if !node.HasTrueChildren() {
			if !node.substitution.IsEmpty() {
				newChvalues.AddToF2Head(newLetter)
			}
			return newChvalues, equation.NewCharacteristicValues(), cycleToNode, false
		} else {
			if newChvaluesHelp.IsEmpty() {
				if !node.substitution.IsEmpty() {
					if metTrueNodeAll {
						newChvalues.AddToF1Head(newLetter)
					} else {
						newChvalues.AddToF2Head(newLetter)
					}
				}
				return newChvalues, equation.NewCharacteristicValues(), cycleToNode, metTrueNodeAll
			} else {
				f1E, f2E := newChvaluesHelp.Compare(newChvalues)
				if f1E && f2E {
					if !node.substitution.IsEmpty() {
						if metTrueNodeAll {
							newChvalues.AddToF1Head(newLetter)
						} else {
							newChvalues.AddToF2Head(newLetter)
						}
					}
					if node.IsSubgraphRoot() {
						return newChvalues, equation.NewCharacteristicValues(), cycleToNode, false
					} else {
						return equation.NewCharacteristicValues(), newChvalues, cycleToNode, false
					}
				} else {
					if !node.substitution.IsEmpty() {
						if metTrueNodeAll {
							newChvalues.AddToF1Head(newLetter)
						} else {
							newChvalues.AddToF2Head(newLetter)
						}
					}
					return newChvalues, newChvaluesHelp, cycleToNode, metTrueNodeAll
				}
			}
		}
	}
}
