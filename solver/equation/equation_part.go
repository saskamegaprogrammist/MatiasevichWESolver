package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type EqPart struct {
	Symbols   []symbol.Symbol
	Length    int
	Structure Structure
}

func EmptyEqPart() EqPart {
	return EqPart{
		Symbols:   make([]symbol.Symbol, 0),
		Structure: EmptyStructure(),
	}
}

func (eqPart *EqPart) GetSymbolMode(i int, mode int) (symbol.Symbol, error) {
	switch mode {
	case FORWARD:
		return eqPart.GetSymbol(i)
	case BACKWARDS:
		return eqPart.GetSymbolFromEnd(i)
	default:
		return nil, fmt.Errorf("wrong mode: %v", mode)
	}
}

func (eqPart *EqPart) GetSymbol(i int) (symbol.Symbol, error) {
	if i >= eqPart.Length {
		return nil, fmt.Errorf("index is out of range")
	}
	return eqPart.Symbols[i], nil
}

func (eqPart *EqPart) GetSymbolFromEnd(i int) (symbol.Symbol, error) {
	if i >= eqPart.Length {
		return nil, fmt.Errorf("index is out of range")
	}
	return eqPart.Symbols[eqPart.Length-1-i], nil
}

func (eqPart *EqPart) Copy() EqPart {
	newEqPart := EqPart{}
	newEqPart.Symbols = make([]symbol.Symbol, eqPart.Length)
	copy(newEqPart.Symbols, eqPart.Symbols)
	newEqPart.Length = eqPart.Length
	newEqPart.Structure = eqPart.Structure.Copy()
	return newEqPart
}

func (eqPart *EqPart) IsEmpty() bool {
	return eqPart.Length == 0 || (eqPart.Length == 1 && symbol.IsEmpty(eqPart.Symbols[0]))
}

func NewEqPartFromSymbols(symbols []symbol.Symbol) EqPart {
	eqPart := EmptyEqPart()
	eqPart.Length = len(symbols)
	eqPart.Symbols = symbols
	for _, s := range symbols {
		eqPart.Structure.Add(s)
	}
	return eqPart
}
