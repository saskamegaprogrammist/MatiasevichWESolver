package symbol

import (
	"fmt"
	"reflect"
)

const (
	CONSTANT    = 1
	VARIABLE    = 2
	EMPTY       = 3
	WORD        = 4
	emptySymbol = "$"
	len         = 1
)

type Symbol interface {
	Value() string
}

type Constant struct {
	value string
}

func (constant Constant) Value() string {
	return constant.value
}

type Variable struct {
	value string
}

func (variable Variable) Value() string {
	return variable.value
}

type EmptySymbol struct {
}

func (empty EmptySymbol) Value() string {
	return emptySymbol
}

type Word struct {
	value string
}

func (word Word) Value() string {
	return word.value
}

func (word Word) Len() int {
	return len
}

func Empty() EmptySymbol {
	return EmptySymbol{}
}

func Const(value string) Constant {
	return Constant{value: value}
}

func Var(value string) Variable {
	return Variable{value: value}
}

func WordVar(value string) Word {
	return Word{value: value}
}

func IsEmptyValue(value string) bool {
	return value == emptySymbol
}

func IsEmpty(sym Symbol) bool {
	return reflect.TypeOf(sym) == reflect.TypeOf(EmptySymbol{})
}

func IsConst(sym Symbol) bool {
	return reflect.TypeOf(sym) == reflect.TypeOf(Constant{})
}

func IsVar(sym Symbol) bool {
	return reflect.TypeOf(sym) == reflect.TypeOf(Variable{})
}

func IsWord(sym Symbol) bool {
	return reflect.TypeOf(sym) == reflect.TypeOf(Word{})
}

func IsVarOrWord(sym Symbol) bool {
	return IsWord(sym) || IsVar(sym)
}

func NewSymbol(symbolType int, value string) (Symbol, error) {
	switch symbolType {
	case CONSTANT:
		return Const(value), nil
	case VARIABLE:
		return Var(value), nil
	case WORD:
		return WordVar(value), nil
	case EMPTY:
		return Empty(), nil
	default:
		return nil, fmt.Errorf("invalid symbol type: %d", symbolType)
	}
}
