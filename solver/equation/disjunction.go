package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

type Disjunction struct {
	variable symbol.Symbol
	values   [][]symbol.Symbol
}

func NewDisjunction() Disjunction {
	return Disjunction{
		variable: nil,
		values:   make([][]symbol.Symbol, 0),
	}
}

func (d *Disjunction) SetVariable(variable symbol.Symbol) {
	d.variable = variable
}

func (d *Disjunction) AddValue(value []symbol.Symbol) {
	d.values = append(d.values, value)
}

func (d *Disjunction) SetValues(values [][]symbol.Symbol) {
	if values != nil {
		d.values = values
	}
}
