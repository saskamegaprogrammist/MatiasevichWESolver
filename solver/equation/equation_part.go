package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

type EqPart struct {
	Symbols   []symbol.Symbol
	Length    int
	Structure Structure
}

func (eqPart *EqPart) New() {
	eqPart.Structure.New()
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

func (eqPart *EqPart) NewFromSymbols(symbols []symbol.Symbol) {
	eqPart.Length = len(symbols)
	eqPart.Symbols = symbols
	eqPart.Structure.New()
	for _, s := range symbols {
		eqPart.Structure.Add(s)
	}
}
