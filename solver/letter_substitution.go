package solver

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"

type LetterSusbstitution struct {
	nodeToTrue    *Node
	substitutions []equation.Substitution
}

func NewLetterSubstitution(node *Node) LetterSusbstitution {
	return LetterSusbstitution{
		nodeToTrue:    node,
		substitutions: make([]equation.Substitution, 0),
	}
}

func (ls *LetterSusbstitution) AddSubstToHead(s equation.Substitution) {
	ls.substitutions = append([]equation.Substitution{s}, ls.substitutions...)
}

func (ls *LetterSusbstitution) Copy(copyLS LetterSusbstitution) {
	ls.nodeToTrue = copyLS.nodeToTrue
	for _, s := range copyLS.substitutions {
		ls.substitutions = append(ls.substitutions, s)
	}
}

func (ls *LetterSusbstitution) HasNoSubstitutions() bool {
	return len(ls.substitutions) == 0
}

func (ls *LetterSusbstitution) NewEquation(oldEquation *equation.Equation) equation.Equation {
	var eq = oldEquation
	var newEq equation.Equation
	for _, s := range ls.substitutions {
		newEq = eq.Substitute(s)
		eq = &newEq
	}
	return newEq
}

func (ls *LetterSusbstitution) SubstitutionsAsEquations() []equation.Equation {
	var eqs = make([]equation.Equation, 0)
	for _, ls := range ls.substitutions {
		eqs = append(eqs, ls.ToEquation())
	}
	return eqs
}
