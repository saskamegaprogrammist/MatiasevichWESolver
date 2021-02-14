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
	//fmt.Println(constAlphabet.words)
	//fmt.Print(varsAlphabet.words)
	var err error
	isEq, i := checkEquation(eq)
	if isEq {
		eqleftPart := eq[0 : i-1]
		var leftSymbols []symbol.Symbol
		if eqleftPart == "" {
			leftSymbols = append(leftSymbols, symbol.Empty())
		} else {
			leftSymbols, err = matchWithAlphabetsWithSpace(eqleftPart, constAlphabet, varsAlphabet)
			if err != nil {
				return fmt.Errorf("error matching alphabet: %v", err)
			}
		}
		eqRightPart := eq[i+2:]
		var rightSymbols []symbol.Symbol
		if eqRightPart == "" {
			rightSymbols = append(rightSymbols, symbol.Empty())
		} else {
			rightSymbols, err = matchWithAlphabetsWithSpace(eqRightPart, constAlphabet, varsAlphabet)
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
	return fmt.Errorf("invalid equation: %s", eq)
}

func checkEquation(eq string) (bool, int) {
	eqLen := len(eq) - 1
	for i := 0; i < eqLen; i++ {
		if string(eq[i]) == EQUALS && string(eq[i-1]) == SPACE && string(eq[i+1]) == SPACE {
			return true, i
		}
	}
	return false, 0
}

func matchWithAlphabetsWithSpace(eqPart string, constAlphabet *Alphabet, varsAlphabet *Alphabet) ([]symbol.Symbol, error) {
	var symbols []symbol.Symbol
	var word string
	var err error
	var matchType int
	for _, eqSym := range eqPart {
		eqSymString := string(eqSym)
		if eqSymString != SPACE {
			word += eqSymString
		} else if word != "" {
			matchType, err = matchWord(word, varsAlphabet, constAlphabet)
			if err != nil {
				return symbols, fmt.Errorf("error matching word: %v", err)
			}
			symb, err := symbol.NewSymbol(matchType, word)
			if err != nil {
				return symbols, fmt.Errorf("error creating symbol: %v", err)
			}
			symbols = append(symbols, symb)
			word = ""
		}
	}
	if word != "" {
		matchType, err = matchWord(word, varsAlphabet, constAlphabet)
		if err != nil {
			return symbols, fmt.Errorf("error matching word: %v", err)
		}
		symb, err := symbol.NewSymbol(matchType, word)
		if err != nil {
			return symbols, fmt.Errorf("error creating symbol: %v", err)
		}
		symbols = append(symbols, symb)
	}
	return symbols, nil
}

func matchWord(word string, varsAlphabet *Alphabet, constAlphabet *Alphabet) (int, error) {
	var matchVar bool
	var matchConst bool
	var matchEmpty bool
	var matchType int
	matchVar = findInAlphabet(word, varsAlphabet)
	if matchVar {
		matchType = symbol.VARIABLE
	}
	matchConst = findInAlphabet(word, constAlphabet)
	if matchConst {
		matchType = symbol.CONSTANT
	}

	if matchConst && matchVar {
		return matchType, fmt.Errorf("variable and constant found for word: %s", word)
	}

	if symbol.IsEmptyValue(word) {
		matchEmpty = true
		matchType = symbol.EMPTY
	}
	if !(matchConst || matchEmpty || matchVar) {
		return matchType, fmt.Errorf("no match found with word: %s", word)
	}
	return matchType, nil
}

func findInAlphabet(word string, alphabet *Alphabet) bool {
	if len(word) <= alphabet.maxWordLength {
		for _, vWord := range alphabet.words {
			if word == vWord {
				return true
			}
		}
	}
	return false
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
		match = findInAlphabet(currentWord, varsAlphabet)
		if match {
			matchType = symbol.VARIABLE
		}
		match = findInAlphabet(currentWord, constAlphabet)
		if match {
			matchType = symbol.CONSTANT
		}

		if symbol.IsEmptyValue(currentWord) {
			match = true
			matchType = symbol.EMPTY
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
				i = startIndex + len(lastMatchedWord)
				startIndex = i
				currentWord = ""
				lastMatchedWord = ""
				lastMatchType = 0
				if i >= eqLen {
					break
				}
			}
		}
	}
	//fmt.Print(symbols)
	return symbols, nil
}

func (equation *Equation) CheckInequality() bool {
	//equation.Print()
	if equation.IsRightEmpty() {
		if equation.leftLength > 0 {
			counter := 0
			for _, sym := range equation.leftPart {
				if symbol.IsWord(sym) || symbol.IsConst(sym) {
					counter++
				}
			}
			if counter != 0 {
				return true
			}
		}
	}
	if equation.IsLeftEmpty() {
		if equation.rightLength > 0 {
			counter := 0
			for _, sym := range equation.rightPart {
				if symbol.IsWord(sym) || symbol.IsConst(sym) {
					counter++
				}
			}
			if counter != 0 {
				return true
			}
		}
	}
	return false
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
			if sym != equation.rightPart[i] {
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
	var wordsMap = map[string]string{}
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
			if sym != eq.leftPart[i] {
				if symbol.IsWord(sym) && symbol.IsWord(eq.leftPart[i]) {
					if wordsMap[sym.Value()] == "" {
						wordsMap[sym.Value()] = eq.leftPart[i].Value()
						i++
					} else {
						if wordsMap[sym.Value()] != eq.leftPart[i].Value() {
							return false
						} else {
							i++
						}
					}
				} else {
					return false
				}
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
			if sym != eq.rightPart[i] {
				if symbol.IsWord(sym) && symbol.IsWord(eq.rightPart[i]) {
					if wordsMap[sym.Value()] == "" {
						wordsMap[sym.Value()] = eq.rightPart[i].Value()
						i++
					} else {
						if wordsMap[sym.Value()] != eq.rightPart[i].Value() {
							return false
						} else {
							i++
						}
					}
				} else {
					return false
				}
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
			if symbol.IsConst(sym) || symbol.IsWord(sym) {
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
			if symbol.IsConst(sym) || symbol.IsWord(sym) {
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
		if sym == (*symbol) {
			resultEquation.leftPart = append(resultEquation.leftPart, newSymbols...)
			resultEquation.leftLength += newSymLen
		} else {
			resultEquation.leftPart = append(resultEquation.leftPart, sym)
			resultEquation.leftLength++
		}
	}
	for _, sym := range equation.rightPart {
		if sym == (*symbol) {
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
	equation.reduceEmpty()
	minLen := min(equation.leftLength, equation.rightLength)
	i := 0
	for ; i < minLen; i++ {
		if equation.leftPart[i] == equation.rightPart[i] {

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
			equation.rightPart = equation.rightPart[i:]
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

func (equation *Equation) Print() {
	fmt.Println(equation.String())
}

func (equation *Equation) String() string {
	var result string
	for _, sym := range equation.leftPart {
		result += fmt.Sprintf("%s ", sym.Value())
	}
	result += fmt.Sprintf("%s ", EQUALS)
	for _, sym := range equation.rightPart {
		result += fmt.Sprintf("%s ", sym.Value())
	}
	return result
}
