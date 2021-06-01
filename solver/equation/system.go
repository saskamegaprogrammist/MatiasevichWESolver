package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"strings"
)

type EquationsSystem struct {
	value      Equation
	compounds  []EquationsSystem
	systemType int
}

func (es *EquationsSystem) Size() int {
	if es.IsSingleEquation() {
		return 1
	}
	return len(es.compounds)
}

func (es *EquationsSystem) Compounds() []EquationsSystem {
	return es.compounds
}

func (es *EquationsSystem) GetEquations() []Equation {
	var equations = make([]Equation, 0)
	if es.IsSingleEquation() {
		return []Equation{es.value}
	}
	for _, c := range es.compounds {
		equations = append(equations, c.GetEquations()...)
	}
	return equations
}

func (es *EquationsSystem) Equation() *Equation {
	if es.IsSingleEquation() {
		return &es.value
	}
	return es.compounds[0].Equation()
}

func (es *EquationsSystem) Copy() EquationsSystem {
	newEqSys := EquationsSystem{}
	newEqSys.value = es.value.Copy()
	newEqSys.compounds = make([]EquationsSystem, len(es.compounds))
	copy(newEqSys.compounds, es.compounds)
	newEqSys.systemType = es.systemType
	return newEqSys
}

func (es *EquationsSystem) Print() {
	fmt.Print(es.String())
}

func (es *EquationsSystem) String() string {
	return es.string(0)
}

func (es *EquationsSystem) string(level int) string {
	var result string
	result = strings.Repeat(" ", level)
	if es.value.IsEmpty() {
		result += typesMap[es.systemType]
		result += "\n"
		for _, c := range es.compounds {
			result += c.string(level + 1)
		}
	} else {
		result += es.value.String()
		result += "\n"
	}
	return result
}

func (es *EquationsSystem) IsSingleEquation() bool {
	return es.systemType == SINGLE_EQUATION
}

func (es *EquationsSystem) IsDisjunction() bool {
	return es.systemType == DISJUNCTION
}

func (es *EquationsSystem) IsConjunction() bool {
	return es.systemType == CONJUNCTION
}

func (es *EquationsSystem) CheckInequality() bool {
	if es.IsSingleEquation() {
		return es.value.CheckInequality()
	}
	for _, eq := range es.compounds {
		if eq.CheckInequality() {
			return true
		}
	}
	return false
}

func (es *EquationsSystem) CheckEquality() bool {
	if es.IsSingleEquation() {
		return es.value.CheckEquality()
	}
	for _, eq := range es.compounds {
		if !eq.CheckEquality() {
			return false
		}
	}
	return true
}

func (es *EquationsSystem) HasEqSystem(system EquationsSystem) bool {
	if es.IsSingleEquation() {
		if system.IsSingleEquation() {
			return es.value.CheckSameness(&system.value)
		}
		return false
	}
	for _, eq := range es.compounds {
		if eq.Equals(system) {
			return true
		}
	}
	return false
}

func (es *EquationsSystem) Equals(system EquationsSystem) bool {
	if es.systemType == system.systemType {
		if !es.value.IsEmpty() && !system.value.IsEmpty() {
			return es.value.CheckSameness(&system.value)
		} else if len(es.compounds) == len(system.compounds) {
			compoundsMap := make(map[*EquationsSystem]bool)
			for _, c := range es.compounds {
				for _, sc := range system.compounds {
					if compoundsMap[&sc] == true {
						continue
					}
					if c.Equals(sc) {
						compoundsMap[&sc] = true
						break
					}
				}
			}
			if len(compoundsMap) == len(system.compounds) {
				return true
			}
		}
	}
	return false
}

func NewSingleEquation(eq Equation) EquationsSystem {
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: SINGLE_EQUATION,
	}
}

func CharacteristicEquation(sym symbol.Symbol, values VariableValues, eqType int) (EquationsSystem, error) {
	if values.Size() != 2 {
		return EquationsSystem{}, fmt.Errorf("wrong values len: %d", values.Size())
	}
	var leftPart, rightPart []symbol.Symbol
	if eqType == EQ_TYPE_SIMPLE {
		leftPart = append(leftPart, values[0]...)
		leftPart = append(leftPart, values[1]...)
	} else if eqType == EQ_TYPE_W_EMPTY {
		leftPart = append(leftPart, values[1]...)
		leftPart = append(leftPart, values[0]...)
	}

	leftPart = append(leftPart, sym)
	rightPart = append(rightPart, sym)
	rightPart = append(rightPart, values[1]...)
	rightPart = append(rightPart, values[0]...)

	var eq Equation
	eq.NewFromParts(leftPart, rightPart)
	eq.FullReduceEmpty()
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: SINGLE_EQUATION,
	}, nil
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

func NewConjunctionSystemFromEquations(equations []Equation) EquationsSystem {
	eqSystems := make([]EquationsSystem, 0)
	for _, eq := range equations {
		eqSystems = append(eqSystems, NewSingleEquation(eq))
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  eqSystems,
		systemType: CONJUNCTION,
	}
}

func NewDisjunctionSystemFromEquations(equations []Equation) EquationsSystem {
	eqSystems := make([]EquationsSystem, 0)
	for _, eq := range equations {
		eqSystems = append(eqSystems, NewSingleEquation(eq))
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  eqSystems,
		systemType: DISJUNCTION,
	}
}

func NewConjunctionSystem(equations []EquationsSystem) EquationsSystem {
	if len(equations) == 1 {
		return equations[0]
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  equations,
		systemType: CONJUNCTION,
	}
}

func NewDisjunctionSystem(equations []EquationsSystem) EquationsSystem {
	if len(equations) == 1 {
		return equations[0]
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  equations,
		systemType: DISJUNCTION,
	}
}
