package solver

type Node struct {
	Parent   *Node
	Children []*Node
	Value    Equation
}

func (node *Node) IsTree() bool {
	return node.Parent == nil
}
