package equation

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type Substitution struct {
	leftPart  symbol.Symbol
	rightPart []symbol.Symbol
	sType     int
}

func (s *Substitution) Copy() Substitution {
	newSubst := Substitution{}
	newSubst.leftPart = s.leftPart
	newSubst.rightPart = make([]symbol.Symbol, len(s.rightPart))
	copy(newSubst.rightPart, s.rightPart)
	return newSubst
}

func (s *Substitution) IsEmpty() bool {
	return len(s.rightPart) == 0 && s.leftPart == nil
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
func (s *Substitution) IsToLetter() bool {
	return symbol.IsLetter(s.leftPart)
}
func (s *Substitution) IsTo() string {
	return s.leftPart.Value()
}

func (s *Substitution) ToEquation() Equation {
	if len(s.rightPart) == 0 {
		return NewEquation([]symbol.Symbol{s.leftPart}, []symbol.Symbol{symbol.EmptySymbol{}})
	} else {
		return NewEquation([]symbol.Symbol{s.leftPart}, s.rightPart[:1])
	}
}

func (s *Substitution) String() string {
	if s.sType == STANDARD {
		value := s.leftPart.Value() + "->"
		for _, sym := range s.rightPart {
			value += sym.Value()
		}
		return value
	} else {
		return sTypesMap[s.sType]
	}
}

func NewSubstitution(leftPart symbol.Symbol, rightPart []symbol.Symbol) Substitution {
	return Substitution{
		leftPart:  leftPart,
		rightPart: rightPart,
		sType:     STANDARD,
	}
}

func NewSubstitutionSplit() Substitution {
	return Substitution{
		sType: SPLITTING,
	}
}

func NewSubstitutionReduce() Substitution {
	return Substitution{
		sType: REDUCING,
	}
}

func NewSubstituionApply() Substitution {
	return Substitution{
		sType: APPLYING,
	}
}
