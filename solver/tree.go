package solver

import "fmt"

type Node struct {
	Number   string
	Parent   *Node
	Children []*Node
	Value    Equation
}

func (node *Node) IsTree() bool {
	return node.Parent == nil
}

func (node *Node) Print() {
	fmt.Printf("%s : ", node.Number)
	node.Value.Print()
	fmt.Println()
}
