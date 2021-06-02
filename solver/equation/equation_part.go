package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

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
