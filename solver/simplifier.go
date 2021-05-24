package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type Simplifier struct {
}

func (s *Simplifier) Simplify(node Node) error {
	err := s.walk(&node)
	if err != nil {
		return fmt.Errorf("error walking node: %v", err)
	}
	return nil
}

func (s *Simplifier) walk(node *Node) error {
	var err error
	if node.HasFalseChildren() {
		return nil
	}
	if node.HasSingleSubstituteVar() {
		values := s.walkWithSymbol(node)
		values.ReduceEmptyVars()
		var sVar symbol.Symbol
		sVar, err = node.SubstituteVar()
		if err != nil {
			return fmt.Errorf("error getting substitution var: %v", err)
		}
		es := equation.SystemFromValues(sVar, values)
		node.SetSimplifiedRepresentation(es)
		node.RemoveSubstituteVar(sVar)
		node.simplified.Print()
		return nil
	}
	for _, ch := range node.Children() {
		err = s.walk(ch)
		if err != nil {
			return fmt.Errorf("error walking child: %v", err)
		}
		return nil
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
		chValues = s.walkWithSymbol(ch)
		if chValues.IsEmpty() {
			values.AddValue([]symbol.Symbol{ch.Substitution().RightPart()[0]})
		} else {
			values.AddToEachValue(chValues, []symbol.Symbol{ch.Substitution().RightPart()[0]})
		}
	}
	return values
}
