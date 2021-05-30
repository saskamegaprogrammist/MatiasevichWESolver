package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Node struct {
	helpMap                 map[int]bool
	number                  string
	substitution            equation.Substitution
	childrenSubstituteVars  map[symbol.Symbol]int
	subgraphsSubstituteVars map[symbol.Symbol]int
	parent                  *Node
	parentsFromBackCycles   []*Node
	children                []*Node
	value                   equation.Equation
	simplified              equation.EquationsSystem
	subgraphRoot            bool
}

func (node *Node) Copy(original *Node) {
	// not copying links to parent and children !!!
	node.number = original.number
	node.substitution = original.substitution.Copy()
	node.value = original.value.Copy()
	node.simplified = original.simplified.Copy()
	node.helpMap = make(map[int]bool)
	standart.CopyIntBoolMap(&original.helpMap, &node.helpMap)
	node.childrenSubstituteVars = make(map[symbol.Symbol]int)
	//standart.CopySymbolIntMap(&node.childrenSubstituteVars, &newNode.childrenSubstituteVars)
	node.subgraphsSubstituteVars = make(map[symbol.Symbol]int)
	//standart.CopySymbolIntMap(&node.subgraphsSubstituteVars, &newNode.subgraphsSubstituteVars)
}

func (node *Node) Children() []*Node {
	return node.children
}

func (node *Node) NewLetter() symbol.Symbol {
	return node.substitution.RightPart()[0]
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

func (node *Node) HasCycleParents() bool {
	return len(node.parentsFromBackCycles) > 0
}

func (node *Node) HasCycleParent(parent *Node) bool {
	for _, p := range node.parentsFromBackCycles {
		if p.number == parent.number {
			return true
		}
	}
	return false
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
		}
		if !node.helpMap[HAS_FALSE] && ch.helpMap[HAS_FALSE] {
			node.helpMap[HAS_FALSE] = true
		}
		if !node.helpMap[HAS_BACK_CYCLE] && ch.helpMap[HAS_BACK_CYCLE] {
			node.helpMap[HAS_BACK_CYCLE] = true
		}
	}
}

func (node *Node) FillSubstituteMapsFromChildren() {
	for _, ch := range node.children {
		for k, v := range ch.childrenSubstituteVars {
			node.childrenSubstituteVars[k] += v
		}
		for k, v := range ch.subgraphsSubstituteVars {
			node.subgraphsSubstituteVars[k] += v
		}
	}
}

func (node *Node) ClearSubstituteMapsFromChildren() {
	node.childrenSubstituteVars = make(map[symbol.Symbol]int)
	node.subgraphsSubstituteVars = make(map[symbol.Symbol]int)
}

func (node *Node) SetIsSubgraphRoot() {
	node.subgraphRoot = true
}

func (node *Node) UnsetIsSubgraphRoot() {
	node.subgraphRoot = false
}

func (node *Node) IsSubgraphRoot() bool {
	return node.subgraphRoot
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

func (node *Node) HasOnlyFalseChildren() bool {
	return node.helpMap[HAS_FALSE] && !node.helpMap[HAS_TRUE] && !node.helpMap[HAS_BACK_CYCLE]
}

func (node *Node) HasFalseChildrenAndBackCycles() bool {
	return node.helpMap[HAS_FALSE] && node.helpMap[HAS_BACK_CYCLE] && !node.helpMap[HAS_TRUE]
}

func (node *Node) HasOnlyTrueChildren() bool {
	return node.helpMap[HAS_TRUE] && !node.helpMap[HAS_FALSE]
}

func (node *Node) HasBackCycle() bool {
	return node.helpMap[HAS_BACK_CYCLE]
}

func (node *Node) SetChildren(children []*Node) {
	if children != nil {
		node.children = children
	}
}

func (node *Node) AddParentFromBackCycle(child *Node) {
	node.parentsFromBackCycles = append(node.parentsFromBackCycles, child)
}

func (node *Node) AddSubstituteVar(subVar symbol.Symbol) {
	node.childrenSubstituteVars[subVar]++
	if node.HasSingleSubstituteVar() {
		node.subgraphsSubstituteVars[subVar]++
	}
}

func (node *Node) SetSimplifiedRepresentation(es equation.EquationsSystem) {
	node.simplified = es
}

func (node *Node) RemoveSubstituteVar(subVar symbol.Symbol, len int) {
	tr := node
	for tr != nil {
		node.childrenSubstituteVars[subVar] -= len
		if node.childrenSubstituteVars[subVar] == 0 {
			delete(node.childrenSubstituteVars, subVar)
		}
		node.removeSubgraphSubstituteVar(subVar, len)
		tr = tr.parent
	}
}

func (node *Node) removeSubgraphSubstituteVar(subVar symbol.Symbol, len int) {
	if node.subgraphsSubstituteVars[subVar] == 0 {
		if node.HasSingleSubstituteVar() {
			node.subgraphsSubstituteVars[subVar]++
		}
	} else {
		node.subgraphsSubstituteVars[subVar] -= len
		if node.subgraphsSubstituteVars[subVar] == 0 {
			delete(node.subgraphsSubstituteVars, subVar)
		}
	}
}

func (node *Node) HasSingleSubstituteVar() bool {
	return len(node.childrenSubstituteVars) == 1
}

func (node *Node) SubstituteVar() (symbol.Symbol, error) {
	if !node.HasSingleSubstituteVar() {
		return nil, fmt.Errorf("node has several substitute vars")
	}
	return node.children[0].substitution.LeftPart(), nil
}

func NewNode(sub equation.Substitution, number string, parent *Node, val equation.Equation) Node {
	return Node{
		number:                  number,
		substitution:            sub,
		parent:                  parent,
		children:                make([]*Node, 0),
		parentsFromBackCycles:   make([]*Node, 0),
		value:                   val,
		childrenSubstituteVars:  make(map[symbol.Symbol]int),
		subgraphsSubstituteVars: make(map[symbol.Symbol]int),
		helpMap:                 make(map[int]bool),
	}
}

func NewTree(number string, val equation.Equation) Node {
	return Node{
		number:                  number,
		children:                make([]*Node, 0),
		parentsFromBackCycles:   make([]*Node, 0),
		value:                   val,
		childrenSubstituteVars:  make(map[symbol.Symbol]int),
		subgraphsSubstituteVars: make(map[symbol.Symbol]int),
		helpMap:                 make(map[int]bool),
	}
}

func EmptyNode() Node {
	return Node{
		childrenSubstituteVars:  make(map[symbol.Symbol]int),
		subgraphsSubstituteVars: make(map[symbol.Symbol]int),
		children:                make([]*Node, 0),
		parentsFromBackCycles:   make([]*Node, 0),
		helpMap:                 make(map[int]bool),
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
