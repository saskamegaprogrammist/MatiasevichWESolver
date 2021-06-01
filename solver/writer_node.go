package solver

import (
	"reflect"
)

const (
	TRUE  = "TRUE"
	FALSE = "FALSE"

	REGULAR_FALSE          = 0
	FAILED_LENGTH_ANALISYS = 1
)

var falseTypeMap = map[int]string{REGULAR_FALSE: "", FAILED_LENGTH_ANALISYS: "failed length analysis"}

type InfoNode interface {
	GetValue() string
	GetNumber() string
}

type TrueNode struct {
	value  string
	number string
}

type FalseNode struct {
	value     string
	number    string
	falseType int
}

func (trueNode TrueNode) GetNumber() string {
	return trueNode.number
}

func (trueNode TrueNode) GetValue() string {
	return TRUE
}

func (falseNode FalseNode) GetValue() string {
	return FALSE
}

func (falseNode FalseNode) GetNumber() string {
	return falseNode.number
}

func (falseNode FalseNode) GetInfoLabel() string {
	return falseTypeMap[falseNode.falseType]
}

func IsTrueNode(in InfoNode) bool {
	return reflect.TypeOf(in) == reflect.TypeOf(&TrueNode{})
}

func IsFalseNode(in InfoNode) bool {
	return reflect.TypeOf(in) == reflect.TypeOf(&FalseNode{})
}
