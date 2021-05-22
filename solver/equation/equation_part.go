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

func (eqPart *EqPart) NewFromSymbols(symbols []symbol.Symbol) {
	eqPart.Length = len(symbols)
	eqPart.Symbols = symbols
	eqPart.Structure.New()
	for _, s := range symbols {
		eqPart.Structure.Add(s)
	}
}
