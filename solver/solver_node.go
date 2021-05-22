package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type Node struct {
	number          string
	substitution    equation.Substitution
	childrenSubVars map[symbol.Symbol]bool
	parent          *Node
	children        []*Node
	value           equation.Equation
}

func (node *Node) HasChildren() bool {
	return node.children != nil && len(node.children) > 0
}

func (node *Node) SetNumber(number string) {
	node.number = number
}

func (node *Node) IsTree() bool {
	return node.parent == nil
}

func (node *Node) Print() {
	fmt.Printf("%s : ", node.number)
	node.value.Print()
	fmt.Println()
}

func (node *Node) SetChildren(children []*Node) {
	if children != nil {
		node.children = children
	}
}

func (node *Node) AddSubVar(subVar symbol.Symbol) {
	node.childrenSubVars[subVar] = true
}

func NewNode(sub equation.Substitution, number string, parent *Node, val equation.Equation) Node {
	return Node{
		number:          number,
		substitution:    sub,
		parent:          parent,
		children:        make([]*Node, 0),
		value:           val,
		childrenSubVars: make(map[symbol.Symbol]bool),
	}
}

func NewTree(number string, val equation.Equation) Node {
	return Node{
		number:          number,
		children:        make([]*Node, 0),
		value:           val,
		childrenSubVars: make(map[symbol.Symbol]bool),
	}
}

func EmptyNode() Node {
	return Node{
		childrenSubVars: make(map[symbol.Symbol]bool),
		children:        make([]*Node, 0),
	}
}

type NodeSystem struct {
	number   string
	Parent   *NodeSystem
	Children []*NodeSystem
	Value    equation.EqSystem
}

func (nodeSystem *NodeSystem) IsTree() bool {
	return nodeSystem.Parent == nil
}

func (nodeSystem *NodeSystem) Print() {
	fmt.Printf("%s : ", nodeSystem.number)
	nodeSystem.Value.Print()
	fmt.Println()
}
