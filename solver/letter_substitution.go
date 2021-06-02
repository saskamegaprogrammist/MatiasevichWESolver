package solver

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"

type LetterSubstitution struct {
	nodeToTrue    *Node
	substitutions []equation.Substitution
}

func NewLetterSubstitution(node *Node) LetterSubstitution {
	return LetterSubstitution{
		nodeToTrue:    node,
		substitutions: make([]equation.Substitution, 0),
	}
}

func (ls *LetterSubstitution) AddSubstToHead(s equation.Substitution) {
	ls.substitutions = append([]equation.Substitution{s}, ls.substitutions...)
}

func (ls *LetterSubstitution) Copy(copyLS LetterSubstitution) {
	ls.nodeToTrue = copyLS.nodeToTrue
	for _, s := range copyLS.substitutions {
		ls.substitutions = append(ls.substitutions, s)
	}
}

func (ls *LetterSubstitution) HasNoSubstitutions() bool {
	return len(ls.substitutions) == 0
}

func (ls *LetterSubstitution) NewEquation(oldEquation *equation.Equation) equation.Equation {
	var eq = oldEquation
	var newEq equation.Equation
	for _, s := range ls.substitutions {
		newEq = eq.Substitute(s)
		eq = &newEq
	}
	return newEq
}

func (ls *LetterSubstitution) SubstitutionsAsEquations() []equation.Equation {
	var eqs = make([]equation.Equation, 0)
	for _, ls := range ls.substitutions {
		eqs = append(eqs, ls.ToEquation())
	}
	return eqs
}
