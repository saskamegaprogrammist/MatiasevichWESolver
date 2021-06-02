package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Structure struct {
	lettersLen int
	varsLen    int
	constsLen  int
	letters    map[symbol.Symbol]int
	vars       map[symbol.Symbol]int
	consts     map[symbol.Symbol]int
}

func EmptyStructure() Structure {
	return Structure{
		lettersLen: 0,
		varsLen:    0,
		constsLen:  0,
		letters:    make(map[symbol.Symbol]int),
		vars:       make(map[symbol.Symbol]int),
		consts:     make(map[symbol.Symbol]int),
	}
}

func NewStructure(s []symbol.Symbol) Structure {
	newStructure := EmptyStructure()
	for _, s := range s {
		newStructure.Add(s)
	}
	return newStructure
}

func MergeStructures(left *Structure, right *Structure) Structure {
	merged := EmptyStructure()
	for l, t := range left.letters {
		merged.AddTimes(l, t)
	}
	for v, t := range left.vars {
		merged.AddTimes(v, t)
	}
	for c, t := range left.consts {
		merged.AddTimes(c, t)
	}

	for l, t := range right.letters {
		merged.AddTimes(l, t)
	}
	for v, t := range right.vars {
		merged.AddTimes(v, t)
	}
	for c, t := range right.consts {
		merged.AddTimes(c, t)
	}
	return merged
}

func (str *Structure) Copy() Structure {
	newStructure := Structure{}
	newStructure.letters = make(map[symbol.Symbol]int)
	standart.CopySymbolIntMap(&str.letters, &newStructure.letters)
	newStructure.consts = make(map[symbol.Symbol]int)
	standart.CopySymbolIntMap(&str.consts, &newStructure.consts)
	newStructure.vars = make(map[symbol.Symbol]int)
	standart.CopySymbolIntMap(&str.vars, &newStructure.vars)
	newStructure.constsLen = str.constsLen
	newStructure.lettersLen = str.lettersLen
	newStructure.varsLen = str.varsLen
	return newStructure
}

func (str *Structure) Add(symb symbol.Symbol) {
	str.AddTimes(symb, 1)
}

func (str *Structure) Sub(symb symbol.Symbol) {
	str.AddTimes(symb, -1)
}

func (str *Structure) SubTimes(symb symbol.Symbol, times int) {
	str.AddTimes(symb, -times)
}

func (str *Structure) AddTimes(symb symbol.Symbol, times int) {
	switch {
	case symbol.IsVar(symb):
		str.varsLen += times
		str.vars[symb] += times
		if str.vars[symb] == 0 {
			delete(str.vars, symb)
		}
	case symbol.IsConst(symb):
		str.constsLen += times
		str.consts[symb] += times
		if str.consts[symb] == 0 {
			delete(str.consts, symb)
		}
	case symbol.IsLetter(symb):
		str.lettersLen += times
		str.letters[symb] += times
		if str.letters[symb] == 0 {
			delete(str.letters, symb)
		}
	}
}

func (str *Structure) LettersLen() int {
	return str.lettersLen
}

func (str *Structure) LettersRangeLen() int {
	return len(str.letters)
}

func (str *Structure) VarsLen() int {
	return str.varsLen
}

func (str *Structure) VarsRangeLen() int {
	return len(str.vars)
}

func (str *Structure) VarsAnLettersRangeLen() int {
	return len(str.vars) + len(str.letters)
}

func (str *Structure) Vars() map[symbol.Symbol]int {
	return str.vars
}

func (str *Structure) ConstsLen() int {
	return str.constsLen
}

func (str *Structure) ConstsRangeLen() int {
	return len(str.consts)
}

func (str *Structure) Consts() map[symbol.Symbol]int {
	return str.consts
}

func (str *Structure) Letters() map[symbol.Symbol]int {
	return str.letters
}

func (str *Structure) LettersAndConstsLen() int {
	return str.constsLen + str.lettersLen
}

func (str *Structure) Print() {
	if str.constsLen != 0 {
		fmt.Println("Constants:")
		for s, n := range str.consts {
			fmt.Printf("%s: %d\n", s.Value(), n)
		}
	}
	if str.varsLen != 0 {
		fmt.Println("Vars:")
		for s, n := range str.vars {
			fmt.Printf("%s: %d\n", s.Value(), n)
		}
	}
	if str.lettersLen != 0 {
		fmt.Println("Letters:")
		for s, n := range str.letters {
			fmt.Printf("%s: %d\n", s.Value(), n)
		}
	}
}
