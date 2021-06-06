package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

type VariableValues [][]symbol.Symbol

func NewVariableValues() VariableValues {
	return make([][]symbol.Symbol, 0)
}

func NewVariableValuesArray() []VariableValues {
	return make([]VariableValues, 0)
}

func (vv *VariableValues) IsEmpty() bool {
	return vv == nil || len(*vv) == 0
}

func (vv *VariableValues) AddValue(s []symbol.Symbol) {
	*vv = append(*vv, s)
}

func (vv *VariableValues) AddValueToHead(s []symbol.Symbol) {
	*vv = append([][]symbol.Symbol{s}, *vv...)
}

func (vv *VariableValues) AddToFirstValue(chV VariableValues, s []symbol.Symbol) {
	for i, v := range chV {
		if i == 0 {
			*vv = append(*vv, append(s, v...))
		} else {
			*vv = append(*vv, v)
		}
	}
}

func (vv *VariableValues) AddToEachValue(chV VariableValues, s []symbol.Symbol) {
	for _, v := range chV {
		*vv = append(*vv, append(s, v...))
	}
}

func (vv *VariableValues) AddToSecondValue(chV VariableValues, s []symbol.Symbol) {
	if len(chV) == 1 {
		vv.AddToFirstValue(chV, s)
		return
	}
	for i, v := range chV {
		if i == 1 {
			*vv = append(*vv, append(s, v...))
		} else {
			*vv = append(*vv, v)
		}
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
