package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
	"sort"
	"strings"
)

type EquationsSystem struct {
	value             Equation
	compounds         []EquationsSystem
	systemType        int
	hasOnlyRegOrdered bool
}

func (es *EquationsSystem) Size() int {
	if es.IsSingleEquation() {
		return 1
	}
	return len(es.compounds)
}

func (es *EquationsSystem) SplitByEquideComposability() (bool, EquationsSystem) {
	var splitted bool
	if es.IsEmpty() {
		return splitted, EquationsSystem{}
	}
	var newEs EquationsSystem
	if es.IsSingleEquation() {
		var newEqs = make([]Equation, 0)
		system := es.value.SplitByEquidecomposability()
		if system.Size > 1 {
			splitted = true
		}
		for _, s := range system.Equations {
			if !s.CheckEquality() {
				newEqs = append(newEqs, s)
			}
		}
		if len(newEqs) == 0 {
			return splitted, EquationsSystem{}
		}
		if len(newEqs) == 1 {
			return splitted, NewSingleEquation(newEqs[0])
		}
		newEs = NewConjunctionSystemFromEquations(newEqs)

	}
	if es.IsConjunction() {
		var newCompounds = make([]EquationsSystem, 0)
		for _, c := range es.compounds {
			wasSplitted, newCompound := c.SplitByEquideComposability()
			splitted = splitted || wasSplitted
			if !newCompound.IsEmpty() {
				newCompounds = append(newCompounds, newCompound)
			}
		}
		newEs = NewConjunctionSystem(newCompounds)
	}
	if es.IsDisjunction() {
		var newCompounds = make([]EquationsSystem, 0)
		for _, c := range es.compounds {
			wasSplitted, newCompound := c.SplitByEquideComposability()
			splitted = splitted || wasSplitted
			if !newCompound.IsEmpty() {
				newCompounds = append(newCompounds, newCompound)
			}
		}
		newEs = NewDisjunctionSystem(newCompounds)
	}
	newEs.RemoveEqual()
	return splitted, newEs
}

func (es *EquationsSystem) Compounds() []EquationsSystem {
	if es.IsSingleEquation() {
		return []EquationsSystem{*es}
	}
	return es.compounds
}

func (es *EquationsSystem) Reorder() {
	var newCompounds = make([]EquationsSystem, len(es.compounds))
	var indexes = make([]int, 0)
	var eqs = make([]Equation, 0)
	if es.IsConjunction() {
		for i, c := range es.compounds {
			if !c.IsSingleEquation() || c.value.isEquidecomposable || c.value.isRegularlyOrdered {
				newCompounds[i] = c
			} else {
				indexes = append(indexes, i)
				eqs = append(eqs, c.value)
			}
		}
		sort.Sort(EquationsByLength(eqs))
		for i, ind := range indexes {
			newCompounds[ind] = NewSingleEquation(eqs[i])
		}
		es.compounds = newCompounds
	}
}

func (es *EquationsSystem) Simplify() {
	if es.IsEmpty() || es.IsSingleEquation() {
		return
	}
	var newCompounds = make([]EquationsSystem, 0)
	if es.IsConjunction() {
		for _, c := range es.compounds {
			if c.IsEmpty() {
				continue
			}
			c.Simplify()
			if c.IsConjunction() {
				newCompounds = append(newCompounds, c.compounds...)
			} else {
				newCompounds = append(newCompounds, c)
			}
		}
	}
	if es.IsDisjunction() {
		for _, c := range es.compounds {
			if c.IsEmpty() {
				continue
			}
			c.Simplify()
			if c.IsDisjunction() {
				newCompounds = append(newCompounds, c.compounds...)
			} else {
				newCompounds = append(newCompounds, c)
			}
		}
	}
	es.compounds = newCompounds
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

// gets the first equation in system
func (es *EquationsSystem) Equation() *Equation {
	if es.IsEmpty() {
		return nil
	}
	if es.IsSingleEquation() {
		return &es.value
	}
	return es.compounds[0].Equation()
}

func (es *EquationsSystem) Substitute(substitute *Substitution) EquationsSystem {
	if es.IsSingleEquation() {
		newValue := es.value.Substitute(*substitute)
		if newValue.IsEmpty() {
			return EquationsSystem{}
		}
		return NewSingleEquation(newValue)
	}
	var newCompounds = make([]EquationsSystem, 0)
	for _, c := range es.compounds {
		newCompound := c.Substitute(substitute)
		newCompound.RemoveEqual()
		if newCompound.IsEmpty() {
			continue
		}
		newCompounds = append(newCompounds, newCompound)
	}
	if len(newCompounds) == 0 {
		return EquationsSystem{}
	}
	newEs := EquationsSystem{
		value:             Equation{},
		compounds:         newCompounds,
		systemType:        es.systemType,
		hasOnlyRegOrdered: false,
	}
	newEs.RemoveEqual()
	return newEs
}

func (es *EquationsSystem) SubstituteVarsWithEmpty() (EquationsSystem, map[symbol.Symbol]bool) {
	if es.IsSingleEquation() {
		newValue, vars := es.value.SubstituteVarsWithEmpty()
		if newValue.IsEmpty() {
			return EquationsSystem{}, vars
		}
		return NewSingleEquation(newValue), vars
	}
	var newCompounds = make([]EquationsSystem, 0)
	var vars = make(map[symbol.Symbol]bool)
	for _, c := range es.compounds {
		newCompound, cvars := c.SubstituteVarsWithEmpty()
		standart.MergeMapsBool(&vars, cvars)
		newCompound.RemoveEqual()
		if newCompound.IsEmpty() {
			continue
		}
		newCompounds = append(newCompounds, newCompound)
	}
	newEs := EquationsSystem{
		value:             Equation{},
		compounds:         newCompounds,
		systemType:        es.systemType,
		hasOnlyRegOrdered: false,
	}
	newEs.RemoveEqual()
	return newEs, vars
}

func (es *EquationsSystem) Reduce() (bool, EquationsSystem) {
	var reduced, reducedCurr bool
	if es.IsEmpty() {
		return false, *es
	}
	if es.IsSingleEquation() {
		newEq := es.value.Copy()
		reduced = newEq.Reduce()
		return reduced, NewSingleEquation(newEq)
	}
	var newEqSystems []EquationsSystem
	var n EquationsSystem
	for _, c := range es.compounds {
		reducedCurr, n = c.Reduce()
		reduced = reduced || reducedCurr

		newEqSystems = append(newEqSystems, n)
	}
	if es.IsConjunction() {
		return reduced, NewConjunctionSystem(newEqSystems)
	}
	return reduced, NewDisjunctionSystem(newEqSystems)
}

func (es *EquationsSystem) Apply() (bool, error) {
	if es.IsEmpty() {
		return false, nil
	}
	if es.IsSingleEquation() {
		return false, nil
	}
	var applied bool
	var err error
	var app bool
	if es.IsConjunction() {
		var notEqs []EquationsSystem
		var equidecomposable []Equation
		var notEquidecomposable []Equation
		for _, c := range es.compounds {
			if !c.IsSingleEquation() {
				notEqs = append(notEqs, c)
			} else {
				eq := c.Equation()
				if eq.CheckEquidecomposability() {
					equidecomposable = append(equidecomposable, *eq)
				} else {
					notEquidecomposable = append(notEquidecomposable, *eq)
				}
			}
		}

		var newEq Equation

		var alreadyApplied = make(map[int]bool, 0)

		for i, c := range equidecomposable {
			for j := 0; j < len(notEquidecomposable); j++ {
				app, newEq, err = notEquidecomposable[j].Apply(c)
				if err != nil {
					return applied, fmt.Errorf("error applying: %v", err)
				}
				applied = applied || app
				if app {
					notEquidecomposable[j] = newEq
				}
			}
			for j := 0; j < len(equidecomposable); j++ {
				if i == j || alreadyApplied[i] || alreadyApplied[j] {
					continue
				}
				app, newEq, err = equidecomposable[j].Apply(c)
				if err != nil {
					return applied, fmt.Errorf("error applying: %v", err)
				}
				applied = applied || app
				if app {
					equidecomposable[j] = newEq
					alreadyApplied[j] = true
				}
			}
		}

		for _, c := range notEqs {
			app, err = c.Apply()
			if err != nil {
				return applied, fmt.Errorf("error applying: %v", err)
			}
			applied = applied || app
		}
		if len(notEqs) != 0 {
			newConj := NewConjunctionSystemFromEquations(append(notEquidecomposable, equidecomposable...))
			*es = NewConjunctionSystem(append(notEqs, newConj))

		} else {
			*es = NewConjunctionSystemFromEquations(append(notEquidecomposable, equidecomposable...))
		}

	}
	if es.IsDisjunction() {
		for _, c := range es.compounds {
			app, err = c.Apply()
			if err != nil {
				return applied, fmt.Errorf("error applying: %v", err)
			}
			applied = applied || app
		}
	}
	es.RemoveEqual()
	es.Simplify()
	return applied, nil
}

func (es *EquationsSystem) RemoveEqual() {
	if es.IsEmpty() || es.IsSingleEquation() {
		return
	}
	var cmap = make(map[int]bool)
	clen := len(es.compounds)
OUTER:
	for i, c := range es.compounds {
		for j := i + 1; j < clen; j++ {
			if c.Equals(es.compounds[j]) {
				continue OUTER
			}
		}
		if !(c.IsSingleEquation() && c.value.IsEmpty()) {
			cmap[i] = true
		}
	}
	var newCompounds = make([]EquationsSystem, 0)
	for i, elem := range es.compounds {
		if cmap[i] {
			newCompounds = append(newCompounds, elem)
		}
	}
	if len(newCompounds) == 1 && newCompounds[0].IsSingleEquation() {
		es.systemType = SINGLE_EQUATION
		es.value = newCompounds[0].value
		es.hasOnlyRegOrdered = newCompounds[0].hasOnlyRegOrdered
		es.compounds = nil
	} else {
		es.compounds = newCompounds
	}
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

func (es *EquationsSystem) IsEmpty() bool {
	return es.systemType == EMPTY
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
		if es.IsSingleEquation() && es.IsSingleEquation() {
			if es.value.IsEmpty() && !system.value.IsEmpty() ||
				!es.value.IsEmpty() && system.value.IsEmpty() {
				return false
			} else {
				return es.value.CheckSameness(&system.value)
			}
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

func (es *EquationsSystem) SplitIntoRegOrdered() (EquationsSystem, EquationsSystem, error) {
	if es.IsSingleEquation() {
		if es.value.IsRegularlyOrdered() {
			es.hasOnlyRegOrdered = true
			return *es, NewEmptySystem(), nil
		} else {
			return NewEmptySystem(), *es, nil
		}
	}
	if !es.IsConjunction() {
		return EquationsSystem{}, EquationsSystem{}, fmt.Errorf("can only work with conjunctions")
	}
	if es.hasOnlyRegOrdered {
		return *es, NewEmptySystem(), nil
	}
	eqs := es.GetEquations()
	var regOrderedEqs = make([]Equation, 0)
	var simpleEqs = make([]Equation, 0)
	for _, e := range eqs {
		if e.IsRegularlyOrdered() {
			regOrderedEqs = append(regOrderedEqs, e)
		} else {
			simpleEqs = append(simpleEqs, e)
		}
	}
	return NewRegOrderedConjunctionSystemFromEquations(regOrderedEqs),
		NewConjunctionSystemFromEquations(simpleEqs), nil
}

func (es *EquationsSystem) NeedsSimplification() (bool, error) {
	if !es.hasOnlyRegOrdered {
		return false, fmt.Errorf("equation system doesn't consist of regulary ordered equations")
	}
	if es.IsDisjunction() {
		return false, fmt.Errorf("equation system doesn't must be conjunction or single equation")
	}
	equations := es.GetEquations()
	checked, varsAndLetters, err := checkSingleVarForEquation(equations)
	if err != nil {
		return true, fmt.Errorf("error checking normal form: %v", err)
	}
	if !checked {
		return true, nil
	}
	if len(varsAndLetters) != len(equations) {
		return true, fmt.Errorf("vars and eqs has different length: %d and %d",
			len(varsAndLetters), len(equations))
	}
	var normal bool
	for i, e := range equations {
		normal, err = e.HasVarOrLetterForNormal(varsAndLetters[i])
		if err != nil {
			return true, fmt.Errorf("error checking normal form for equation: %v", err)
		}
		if !normal {
			return true, nil
		}
	}
	return false, nil
}

func checkSingleVarForEquation(eqs []Equation) (bool, []symbol.Symbol, error) {
	var varsAndLettersMap = make(map[symbol.Symbol]bool, 0)
	var varsAndLetters = make([]symbol.Symbol, 0)

	for _, e := range eqs {
		if e.structure.VarsAnLettersRangeLen() != 1 {
			return false, varsAndLetters, nil
		}
		for l := range e.structure.letters {
			// means this letter is already in another equation
			if varsAndLettersMap[l] {
				return false, varsAndLetters, nil
			}
			varsAndLettersMap[l] = true
			varsAndLetters = append(varsAndLetters, l)
		}
		for v := range e.structure.vars {
			// means this var is already in another equation
			if varsAndLettersMap[v] {
				return false, varsAndLetters, nil
			}
			varsAndLettersMap[v] = true
			varsAndLetters = append(varsAndLetters, v)
		}
	}
	if len(varsAndLetters) != len(varsAndLettersMap) {
		return false, varsAndLetters, fmt.Errorf("array and map has different length: %d and %d",
			len(varsAndLetters), len(varsAndLettersMap))
	}
	return true, varsAndLetters, nil
}

func NewSingleEquation(eq Equation) EquationsSystem {
	return EquationsSystem{
		value:      eq,
		compounds:  nil,
		systemType: SINGLE_EQUATION,
	}
}

func NewRegOrderedSingleEquation(eq Equation) EquationsSystem {
	req := NewSingleEquation(eq)
	req.hasOnlyRegOrdered = true
	return req
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

func CharacteristicEquationRefactored(sym symbol.Symbol, values CharacteristicValues, eqType int) (EquationsSystem, error) {
	var leftPart, rightPart []symbol.Symbol
	if eqType == EQ_TYPE_SIMPLE {
		leftPart = append(leftPart, values.F1...)
		leftPart = append(leftPart, values.F2...)
	} else if eqType == EQ_TYPE_W_EMPTY {
		leftPart = append(leftPart, values.F2...)
		leftPart = append(leftPart, values.F1...)
	}

	leftPart = append(leftPart, sym)
	rightPart = append(rightPart, sym)
	rightPart = append(rightPart, values.F2...)
	rightPart = append(rightPart, values.F1...)

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
	if len(equations) == 0 {
		return NewEmptySystem()
	}
	eqSystems := make([]EquationsSystem, 0)
	for _, eq := range equations {
		eqSystems = append(eqSystems, NewSingleEquation(eq))
	}
	if len(eqSystems) == 1 {
		return eqSystems[0]
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  eqSystems,
		systemType: CONJUNCTION,
	}
}

func NewRegOrderedConjunctionSystemFromEquations(equations []Equation) EquationsSystem {
	if len(equations) == 0 {
		return NewEmptySystem()
	}
	eqSystems := make([]EquationsSystem, 0)
	for _, eq := range equations {
		eqSystems = append(eqSystems, NewRegOrderedSingleEquation(eq))
	}
	if len(eqSystems) == 1 {
		return eqSystems[0]
	}
	return EquationsSystem{
		value:             Equation{},
		compounds:         eqSystems,
		systemType:        CONJUNCTION,
		hasOnlyRegOrdered: true,
	}
}

func NewDisjunctionSystemFromEquations(equations []Equation) EquationsSystem {
	if len(equations) == 0 {
		return NewEmptySystem()
	}
	eqSystems := make([]EquationsSystem, 0)
	for _, eq := range equations {
		eqSystems = append(eqSystems, NewSingleEquation(eq))
	}
	if len(eqSystems) == 1 {
		return eqSystems[0]
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  eqSystems,
		systemType: DISJUNCTION,
	}
}

func NewConjunctionSystem(equations []EquationsSystem) EquationsSystem {
	if len(equations) == 0 {
		return NewEmptySystem()
	}
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
	if len(equations) == 0 {
		return NewEmptySystem()
	}
	if len(equations) == 1 {
		return equations[0]
	}
	return EquationsSystem{
		value:      Equation{},
		compounds:  equations,
		systemType: DISJUNCTION,
	}
}

func NewEmptySystem() EquationsSystem {
	return EquationsSystem{
		value:      Equation{},
		systemType: EMPTY,
	}
}
