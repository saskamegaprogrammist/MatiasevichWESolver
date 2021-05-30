package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Simplifier struct {
}

func (s *Simplifier) Simplify(node *Node) error {
	var err error
	symbols := standart.SymbolArrayFromIntMap(node.subgraphsSubstituteVars)
	if len(symbols) == 0 {
		return nil
	}
	subgraphSymbol := symbols[len(symbols)-1]
	var resSystem equation.EquationsSystem
	resSystem, err = s.simplify(node, subgraphSymbol)
	if err != nil {
		return fmt.Errorf("error during simplification: %v", err)
	}
	node.simplified = resSystem
	resSystem.Print()

	return nil
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
		newGraph := Node{}
		err = copyGraph(node, &newGraph, disj, symbolVar)
		if err != nil {
			return equation.EquationsSystem{}, fmt.Errorf("error copying graph: %v", err)
		}
		newGraphs = append(newGraphs, newGraph)
		standart.MergeMaps(&varMap, newGraph.subgraphsSubstituteVars)
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
			if node.value.CheckSameness(&tr.value) {
				break
			}
			tr = tr.parent
		}
		copyNode.children = append(copyNode.children, tr)
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
		copyNode.AddSubstituteVar(newChildNode.substitution.LeftPart())
		copyNode.children = append(copyNode.children, &newChildNode)

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
				if !compareSymbols(values[0][1], values[1][1]) {
					return fmt.Errorf("can't create characteristic equation")
				}
				if len(values[0][0]) == 1 && symbol.IsEmpty(values[0][0][0]) {
					eqType = equation.EQ_TYPE_W_EMPTY
				} else if !compareSymbols(values[0][0], values[1][0]) {
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

func compareSymbols(original []symbol.Symbol, new []symbol.Symbol) bool {
	if len(original) != len(new) {
		return false
	}
	for i, s := range original {
		if s != new[i] {
			return false
		}
	}
	return true
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
