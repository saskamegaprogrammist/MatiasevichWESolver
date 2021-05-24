package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type EquationsSystem struct {
	value      Equation
	compounds  []EquationsSystem
	systemType int
}

func (es *EquationsSystem) Print() {
	if es.value.IsEmpty() {
		for _, c := range es.compounds {
			c.Print()
		}
	} else {
		es.value.Print()
		fmt.Println()
	}
}

func NewSingleEquation(eq Equation) EquationsSystem {
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: SINGLE_EQUATION,
	}
}

func SystemFromValues(leftSym symbol.Symbol, rightPart VariableValues) EquationsSystem {
	if rightPart.Size() == 1 {
		var eq Equation
		eq.NewFromParts([]symbol.Symbol{leftSym}, rightPart[0])
		return EquationsSystem{
			value:      eq,
			compounds:  nil,
			systemType: SINGLE_EQUATION,
		}
	} else {
		eqSystems := make([]EquationsSystem, 0)
		for _, vv := range rightPart {
			var eq Equation
			eq.NewFromParts([]symbol.Symbol{leftSym}, vv)
			eqSystems = append(eqSystems, EquationsSystem{
				value:      eq,
				compounds:  nil,
				systemType: SINGLE_EQUATION,
			})
		}
		return EquationsSystem{
			value:      Equation{},
			compounds:  eqSystems,
			systemType: DISJUNCTION,
		}

	}
}

func NewConjunctionSystem(eq Equation) EquationsSystem {
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: CONJUNCTION,
	}
}

func NewDisjunctionSystem(eq Equation) EquationsSystem {
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: DISJUNCTION,
	}
}
