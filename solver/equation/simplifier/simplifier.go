package simplifier

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver"
)

type Simplifier struct {
}

func (s *Simplifier) Simplify(node solver.Node) {
	s.walk(&node)
}

func (s *Simplifier) walk(node *solver.Node) {
	//variable := node.Substitution.LeftPart()
	//for _, ch := range node.Children {
	//	ch.HasChildren()
	//	s.walk(ch)
	//}
}
