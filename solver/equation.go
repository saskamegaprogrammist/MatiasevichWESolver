package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/symbol"
)

type Equation struct {
	leftPart    []symbol.Symbol
	leftLength  int
	rightPart   []symbol.Symbol
	rightLength int
}

const EQUALS = "="

func (equation *Equation) Init(eq string, constAlphabet *Alphabet, varsAlphabet *Alphabet) error {
	fmt.Println(constAlphabet.words)
	fmt.Print(varsAlphabet.words)
	var err error
	for i := 0; i < len(eq); i++ {
		if string(eq[i]) == EQUALS {
			eqleftPart := eq[0:i]
			var leftSymbols []symbol.Symbol
			if eqleftPart == "" {
				leftSymbols = append(leftSymbols, symbol.Empty())
			} else {
				leftSymbols, err = matchWithAlphabets(eqleftPart, constAlphabet, varsAlphabet)
				if err != nil {
					return fmt.Errorf("error matching alphabet: %v", err)
				}
			}
			eqRightPart := eq[i+1:]
			var rightSymbols []symbol.Symbol
			if eqRightPart == "" {
				rightSymbols = append(rightSymbols, symbol.Empty())
			} else {
				rightSymbols, err = matchWithAlphabets(eqRightPart, constAlphabet, varsAlphabet)
				if err != nil {
					return fmt.Errorf("error matching alphabet: %v", err)
				}
			}
			equation.leftLength = len(leftSymbols)
			equation.leftPart = leftSymbols
			equation.rightLength = len(rightSymbols)
			equation.rightPart = rightSymbols
			return nil
		}
	}
	return fmt.Errorf("invalid equation: %s", eq)
}

func matchWithAlphabets(eqPart string, constAlphabet *Alphabet, varsAlphabet *Alphabet) ([]symbol.Symbol, error) {
	eqLen := len(eqPart)
	var symbols []symbol.Symbol
	var match bool
	var matchType int
	var lastMatchType int
	var lastMatchedWord string
	var currentWord string
	var startIndex int
	var i int
	continueSearch := true
	for {
		currentWord += string(eqPart[i])
		if len(currentWord) <= varsAlphabet.maxWordLength {
			for _, vWord := range varsAlphabet.words {
				if currentWord == vWord {
					matchType = symbol.VARIABLE
					match = true
					break
				}
			}
		}
		if len(currentWord) <= constAlphabet.maxWordLength {
			for _, cWord := range constAlphabet.words {
				if currentWord == cWord {
					matchType = symbol.CONSTANT
					match = true
					break
				}
			}
		}
		if symbol.IsEmptyValue(currentWord) {
			matchType = symbol.EMPTY
			match = true
		}
		nextWordLen := len(currentWord) + 1
		continueSearch = (nextWordLen <= varsAlphabet.maxWordLength || nextWordLen <= constAlphabet.maxWordLength) && eqLen != i+1
		//fmt.Println(nextWordLen, varsAlphabet.maxWordLength, constAlphabet.maxWordLength, continueSearch)
		if match {
			lastMatchedWord = currentWord
			lastMatchType = matchType
			match = false
			if continueSearch {
				i++
			} else {
				symb, err := symbol.NewSymbol(matchType, currentWord)
				if err != nil {
					return symbols, fmt.Errorf("error creating symbol: %v", err)
				}
				symbols = append(symbols, symb)
				currentWord = ""
				lastMatchedWord = ""
				lastMatchType = 0
				i++
				startIndex = i
				if i >= eqLen {
					break
				}
			}
		} else {
			if continueSearch {
				i++
			} else {
				if lastMatchedWord == "" {
					return nil, fmt.Errorf("no match for word: %s", currentWord)
				}
				symb, err := symbol.NewSymbol(lastMatchType, lastMatchedWord)
				if err != nil {
					return symbols, fmt.Errorf("error creating symbol: %v", err)
				}
				symbols = append(symbols, symb)
				currentWord = ""
				lastMatchedWord = ""
				lastMatchType = 0
				i = startIndex + len(lastMatchedWord)
				startIndex = i
				if i >= eqLen {
					break
				}
			}
		}
	}
	fmt.Print(symbols)
	return symbols, nil
}

func (equation *Equation) CheckEquality() bool {
	if equation.rightLength == 0 {
		equation.rightLength++
		equation.rightPart = append(equation.rightPart, symbol.Empty())
	}
	i := 0
	for _, sym := range equation.leftPart {
		if symbol.IsEmpty(sym) {

		} else {
			for i < equation.rightLength && symbol.IsEmpty(equation.rightPart[i]) {
				i++
			}
			if i == equation.rightLength {
				return false
			}
			if sym.Value() != equation.rightPart[i].Value() {
				return false
			} else {
				i++
			}
		}
	}
	for i < equation.rightLength && symbol.IsEmpty(equation.rightPart[i]) {
		i++
	}
	if i != equation.rightLength {
		return false
	}
	return true
}

func (equation *Equation) CheckSameness(eq *Equation) bool {
	if eq.rightLength == 0 {
		eq.rightLength++
		eq.rightPart = append(eq.rightPart, symbol.Empty())
	}
	if eq.leftLength == 0 {
		eq.leftLength++
		eq.leftPart = append(eq.leftPart, symbol.Empty())
	}
	i := 0
	for _, sym := range equation.leftPart {
		if symbol.IsEmpty(sym) {

		} else {
			for i < eq.leftLength && symbol.IsEmpty(eq.leftPart[i]) {
				i++
			}
			if i == eq.leftLength {
				return false
			}
			if sym.Value() != eq.leftPart[i].Value() {
				return false
			} else {
				i++
			}
		}
	}
	for i < eq.leftLength && symbol.IsEmpty(eq.leftPart[i]) {
		i++
	}
	if i != eq.leftLength {
		return false
	}
	i = 0
	for _, sym := range equation.rightPart {
		if symbol.IsEmpty(sym) {

		} else {
			for i < eq.rightLength && symbol.IsEmpty(eq.rightPart[i]) {
				i++
			}
			if i == eq.rightLength {
				return false
			}
			if sym.Value() != eq.rightPart[i].Value() {
				return false
			} else {
				i++
			}
		}
	}
	for i < eq.rightLength && symbol.IsEmpty(eq.rightPart[i]) {
		i++
	}
	if i != eq.rightLength {
		return false
	}
	return true
}

func (equation *Equation) SubstituteVarsWithEmpty() Equation {
	var resultEquation Equation
	if equation.IsRightEmpty() {
		resultEquation.rightPart = equation.rightPart
		resultEquation.rightLength = equation.rightLength
		for _, sym := range equation.leftPart {
			if symbol.IsConst(sym) {
				resultEquation.leftPart = append(resultEquation.leftPart, sym)
			}
		}
		if len(resultEquation.leftPart) == 0 {
			resultEquation.leftPart = append(resultEquation.leftPart, symbol.Empty())
		}
		resultEquation.leftLength = len(resultEquation.leftPart)
	}
	if equation.IsLeftEmpty() {
		resultEquation.leftPart = equation.leftPart
		resultEquation.leftLength = equation.leftLength
		for _, sym := range equation.rightPart {
			if symbol.IsConst(sym) {
				resultEquation.rightPart = append(resultEquation.rightPart, sym)
			}
		}
		if len(resultEquation.rightPart) == 0 {
			resultEquation.rightPart = append(resultEquation.rightPart, symbol.Empty())
		}
		resultEquation.rightLength = len(resultEquation.rightPart)
	}
	return resultEquation
}

func (equation *Equation) Substitute(symbol *symbol.Symbol, newSymbols []symbol.Symbol) Equation {
	newSymLen := len(newSymbols)
	var resultEquation Equation
	for _, sym := range equation.leftPart {
		if sym.Value() == (*symbol).Value() {
			resultEquation.leftPart = append(resultEquation.leftPart, newSymbols...)
			resultEquation.leftLength += newSymLen
		} else {
			resultEquation.leftPart = append(resultEquation.leftPart, sym)
			resultEquation.leftLength++
		}
	}
	for _, sym := range equation.rightPart {
		if sym.Value() == (*symbol).Value() {
			resultEquation.rightPart = append(resultEquation.rightPart, newSymbols...)
			resultEquation.rightLength += newSymLen
		} else {
			resultEquation.rightPart = append(resultEquation.rightPart, sym)
			resultEquation.rightLength++
		}
	}
	resultEquation.Reduce()
	return resultEquation
}

func (equation *Equation) Reduce() {
	minLen := min(equation.leftLength, equation.rightLength)
	i := 0
	for ; i < minLen; i++ {
		if symbol.IsVar(equation.leftPart[i]) &&
			equation.leftPart[i].Value() == equation.rightPart[i].Value() {

		} else {
			break
		}
	}
	if i > 0 {
		equation.rightPart = equation.rightPart[i:]
		equation.rightLength -= i
		equation.leftPart = equation.leftPart[i:]
		equation.leftLength -= i
	}
	equation.reduceEmpty()
}

func (equation *Equation) reduceEmpty() {
	i := 0
	if equation.leftLength > 1 {
		for ; i < equation.leftLength; i++ {
			if symbol.IsEmpty(equation.leftPart[i]) {

			} else {
				break
			}
		}
		if i > 0 {
			equation.leftPart = equation.leftPart[i:]
			equation.leftLength -= i
		}
	}
	if equation.rightLength > 1 {
		i = 0
		for ; i < equation.rightLength; i++ {
			if symbol.IsEmpty(equation.rightPart[i]) {

			} else {
				break
			}
		}
		if i > 0 {
			equation.leftPart = equation.rightPart[i:]
			equation.rightLength -= i
		}
	}
}

func min(first int, second int) int {
	if first < second {
		return first
	} else {
		return second
	}
}

func (equation *Equation) IsLeftEmpty() bool {
	return equation.leftLength == 1 &&
		symbol.IsEmpty(equation.leftPart[0]) || equation.leftLength == 0
}

func (equation *Equation) IsRightEmpty() bool {
	return equation.rightLength == 1 &&
		symbol.IsEmpty(equation.rightPart[0]) || equation.rightLength == 0
}
