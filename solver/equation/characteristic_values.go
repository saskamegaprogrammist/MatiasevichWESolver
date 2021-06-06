package equation

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type CharacteristicValues struct {
	F1 []symbol.Symbol
	F2 []symbol.Symbol
}

func NewCharacteristicValues() CharacteristicValues {
	return CharacteristicValues{
		F1: make([]symbol.Symbol, 0),
		F2: make([]symbol.Symbol, 0),
	}
}

func NewCharacteristicValuesFromArrays(F1 []symbol.Symbol, F2 []symbol.Symbol) CharacteristicValues {
	return CharacteristicValues{
		F1: F1,
		F2: F2,
	}
}

func NewCharacteristicValuesArray() []CharacteristicValues {
	return make([]CharacteristicValues, 0)
}

func (ch *CharacteristicValues) AddToF1(s symbol.Symbol) {
	ch.F1 = append(ch.F1, s)
}

func (ch *CharacteristicValues) AddToF1Head(s symbol.Symbol) {
	ch.F1 = append([]symbol.Symbol{s}, ch.F1...)
}

func (ch *CharacteristicValues) AddToF2(s symbol.Symbol) {
	ch.F2 = append(ch.F2, s)
}

func (ch *CharacteristicValues) AddToF2Head(s symbol.Symbol) {
	ch.F2 = append([]symbol.Symbol{s}, ch.F2...)
}

func (ch *CharacteristicValues) IsEmpty() bool {
	return len(ch.F1) == 0 && len(ch.F2) == 0
}

func (ch *CharacteristicValues) Append(chv CharacteristicValues) {
	ch.F1 = append(ch.F1, chv.F1...)
	ch.F2 = append(ch.F2, chv.F2...)
}

func (ch *CharacteristicValues) Compare(chv CharacteristicValues) (bool, bool) {
	ch.ReduceEmptyVars()
	chv.ReduceEmptyVars()
	return standart.CheckSymbolArraysEquality(ch.F1, chv.F1),
		standart.CheckSymbolArraysEquality(ch.F2, chv.F2)
}

func (ch *CharacteristicValues) ReduceEmptyVars() CharacteristicValues {
	newChv := NewCharacteristicValues()
	for _, v := range ch.F1 {
		if !symbol.IsEmpty(v) {
			newChv.F1 = append(newChv.F1, v)
		}
	}

	for _, v := range ch.F2 {
		if !symbol.IsEmpty(v) {
			newChv.F2 = append(newChv.F2, v)
		}
	}
	return newChv
}
