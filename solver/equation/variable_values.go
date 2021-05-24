package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

type VariableValues [][]symbol.Symbol

func NewVariableValues() VariableValues {
	return make([][]symbol.Symbol, 0)
}

func (vv *VariableValues) SetVariableValues(values [][]symbol.Symbol) {
	if values != nil {
		*vv = values
	}
}

func (vv *VariableValues) IsEmpty() bool {
	return vv == nil || len(*vv) == 0
}

func (vv *VariableValues) AddValue(s []symbol.Symbol) {
	*vv = append(*vv, s)
}

func (vv *VariableValues) AddToEachValue(chV VariableValues, s []symbol.Symbol) {
	for _, v := range chV {
		*vv = append(*vv, append(s, v...))
	}
}

func (vv *VariableValues) Size() int {
	return len(*vv)
}

func (vv *VariableValues) ReduceEmptyVars() {
	var newSymbols [][]symbol.Symbol
	for _, val := range *vv {
		var newS []symbol.Symbol
		for _, v := range val {
			if !symbol.IsEmpty(v) {
				newS = append(newS, v)
			}
		}

		if len(newS) == 0 {
			newS = append(newS, symbol.EmptySymbol{})
		}
		newSymbols = append(newSymbols, newS)
	}
	*vv = newSymbols
}
