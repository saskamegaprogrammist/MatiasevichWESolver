package equation

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

type Substitution struct {
	leftPart  symbol.Symbol
	rightPart []symbol.Symbol
}

func (s *Substitution) Copy() Substitution {
	newSubst := Substitution{}
	newSubst.leftPart = s.leftPart
	newSubst.rightPart = make([]symbol.Symbol, len(s.rightPart))
	copy(newSubst.rightPart, s.rightPart)
	return newSubst
}

func (s *Substitution) SubstitutesToEmpty() bool {
	return len(s.rightPart) == 1 && symbol.IsEmpty(s.rightPart[0])
}

func (s *Substitution) LeftPart() symbol.Symbol {
	return s.leftPart
}

func (s *Substitution) RightPart() []symbol.Symbol {
	return s.rightPart
}

func (s *Substitution) RightPartLength() int {
	return len(s.rightPart)
}

func (s *Substitution) IsToVar() bool {
	return symbol.IsVar(s.leftPart)
}

func (s *Substitution) IsTo() string {
	return s.leftPart.Value()
}

func NewSubstitution(leftPart symbol.Symbol, rightPart []symbol.Symbol) Substitution {
	return Substitution{
		leftPart:  leftPart,
		rightPart: rightPart,
	}
}
