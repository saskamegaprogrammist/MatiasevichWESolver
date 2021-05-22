package solver

const (
	TRUE  = "TRUE"
	FALSE = "FALSE"
)

type InfoNode interface {
	GetValue() string
	GetNumber() string
}

type TrueNode struct {
	value  string
	number string
}

type FalseNode struct {
	value  string
	number string
}

func (trueNode *TrueNode) GetNumber() string {
	return trueNode.number
}

func (trueNode *TrueNode) GetValue() string {
	return TRUE
}

func (falseNode *FalseNode) GetValue() string {
	return FALSE
}

func (falseNode *FalseNode) GetNumber() string {
	return falseNode.number
}
