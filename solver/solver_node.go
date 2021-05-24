package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

type Node struct {
	helpMap         map[int]bool
	number          string
	substitution    equation.Substitution
	childrenSubVars map[symbol.Symbol]bool
	parent          *Node
	children        []*Node
	value           equation.Equation
	simplified      equation.EquationsSystem
}

func (node *Node) Children() []*Node {
	return node.children
}

func (node *Node) Substitution() *equation.Substitution {
	return &node.substitution
}

func (node *Node) HasChildren() bool {
	return node.children != nil && len(node.children) > 0
}

func (node *Node) SetNumber(number string) {
	node.number = number
}

func (node *Node) LeadsToBackCycle() bool {
	return len(node.children) == 1 &&
		len(node.children[0].number) < len(node.number)
}

func (node *Node) IsTree() bool {
	return node.parent == nil
}

func (node *Node) Print() {
	fmt.Printf("%s : ", node.number)
	node.value.Print()
	fmt.Println()
}

func (node *Node) FillHelpMapFromChildren() {
	for _, ch := range node.children {
		if len(node.helpMap) == RANGE {
			break
		}
		if !node.helpMap[HAS_TRUE] && ch.helpMap[HAS_TRUE] {
			node.helpMap[HAS_TRUE] = true
			continue
		}
		if !node.helpMap[HAS_FALSE] && ch.helpMap[HAS_FALSE] {
			node.helpMap[HAS_FALSE] = true
			continue
		}
		if !node.helpMap[HAS_BACK_CYCLE] && ch.helpMap[HAS_BACK_CYCLE] {
			node.helpMap[HAS_BACK_CYCLE] = true
			continue
		}
	}
}

func (node *Node) SetHasTrueChildren() {
	node.helpMap[HAS_TRUE] = true
}

func (node *Node) SetHasFalseChildren() {
	node.helpMap[HAS_FALSE] = true
}

func (node *Node) SetHasBackCycle() {
	node.helpMap[HAS_BACK_CYCLE] = true
}

func (node *Node) HasTrueChildren() bool {
	return node.helpMap[HAS_TRUE]
}

func (node *Node) HasFalseChildren() bool {
	return node.helpMap[HAS_FALSE]
}

func (node *Node) HasBackCycle() bool {
	return node.helpMap[HAS_BACK_CYCLE]
}

func (node *Node) SetChildren(children []*Node) {
	if children != nil {
		node.children = children
	}
}

func (node *Node) AddSubstituteVar(subVar symbol.Symbol) {
	node.childrenSubVars[subVar] = true
}

func (node *Node) SetSimplifiedRepresentation(es equation.EquationsSystem) {
	node.simplified = es
}

func (node *Node) RemoveSubstituteVar(subVar symbol.Symbol) {
	tr := node
	for tr != nil {
		delete(node.childrenSubVars, subVar)
		tr = tr.parent
	}
}

func (node *Node) HasSingleSubstituteVar() bool {
	return len(node.childrenSubVars) == 1
}

func (node *Node) SubstituteVar() (symbol.Symbol, error) {
	if !node.HasSingleSubstituteVar() {
		return nil, fmt.Errorf("node has several substitute vars")
	}
	return node.children[0].substitution.LeftPart(), nil
}

func NewNode(sub equation.Substitution, number string, parent *Node, val equation.Equation) Node {
	return Node{
		number:          number,
		substitution:    sub,
		parent:          parent,
		children:        make([]*Node, 0),
		value:           val,
		childrenSubVars: make(map[symbol.Symbol]bool),
		helpMap:         make(map[int]bool),
	}
}

func NewTree(number string, val equation.Equation) Node {
	return Node{
		number:          number,
		children:        make([]*Node, 0),
		value:           val,
		childrenSubVars: make(map[symbol.Symbol]bool),
		helpMap:         make(map[int]bool),
	}
}

func EmptyNode() Node {
	return Node{
		childrenSubVars: make(map[symbol.Symbol]bool),
		children:        make([]*Node, 0),
		helpMap:         make(map[int]bool),
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
