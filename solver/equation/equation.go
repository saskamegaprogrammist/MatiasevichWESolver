package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Equation struct {
	LeftPart           EqPart
	RightPart          EqPart
	isEquidecomposable bool
}

const EQUALS = "="

func (equation *Equation) IsEmpty() bool {
	return equation.LeftPart.Length == 0 && equation.RightPart.Length == 0
}

func (equation *Equation) New() {
	equation.LeftPart.New()
	equation.RightPart.New()
}

func (equation *Equation) NewFromParts(leftPart []symbol.Symbol, rightPart []symbol.Symbol) {
	equation.LeftPart.NewFromSymbols(leftPart)
	equation.RightPart.NewFromSymbols(rightPart)
}

func (equation *Equation) Init(eq string, constAlphabet *Alphabet, varsAlphabet *Alphabet) error {
	//fmt.Println(constAlphabet.words)
	//fmt.Print(varsAlphabet.words)
	var err error
	isEq, i := checkEquation(eq)
	if isEq {
		eqleftPart := eq[0 : i-1]
		var leftSymbols []symbol.Symbol
		var leftSymbolsStruct Structure
		if eqleftPart == "" {
			leftSymbols = append(leftSymbols, symbol.Empty())
		} else {
			leftSymbols, leftSymbolsStruct, err = matchWithAlphabetsWithSpace(eqleftPart, constAlphabet, varsAlphabet)
			if err != nil {
				return fmt.Errorf("error matching alphabet: %v", err)
			}
		}
		eqRightPart := eq[i+2:]
		var rightSymbols []symbol.Symbol
		var rightSymbolsStruct Structure
		if eqRightPart == "" {
			rightSymbols = append(rightSymbols, symbol.Empty())
		} else {
			rightSymbols, rightSymbolsStruct, err = matchWithAlphabetsWithSpace(eqRightPart, constAlphabet, varsAlphabet)
			if err != nil {
				return fmt.Errorf("error matching alphabet: %v", err)
			}
		}
		equation.LeftPart = EqPart{
			Length:    len(leftSymbols),
			Symbols:   leftSymbols,
			Structure: leftSymbolsStruct,
		}
		equation.RightPart = EqPart{
			Length:    len(rightSymbols),
			Symbols:   rightSymbols,
			Structure: rightSymbolsStruct,
		}
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

func matchWithAlphabetsWithSpace(eqPart string, constAlphabet *Alphabet, varsAlphabet *Alphabet) ([]symbol.Symbol, Structure, error) {
	var symbols []symbol.Symbol
	var symbolsStructure Structure
	symbolsStructure.New()
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
				return symbols, symbolsStructure, fmt.Errorf("error matching word: %v", err)
			}
			symb, err := symbol.NewSymbol(matchType, word)
			if err != nil {
				return symbols, symbolsStructure, fmt.Errorf("error creating symbol: %v", err)
			}
			symbols = append(symbols, symb)
			symbolsStructure.Add(symb)
			word = ""
		}
	}
	if word != "" {
		matchType, err = matchWord(word, varsAlphabet, constAlphabet)
		if err != nil {
			return symbols, symbolsStructure, fmt.Errorf("error matching word: %v", err)
		}
		symb, err := symbol.NewSymbol(matchType, word)
		if err != nil {
			return symbols, symbolsStructure, fmt.Errorf("error creating symbol: %v", err)
		}
		symbols = append(symbols, symb)
		symbolsStructure.Add(symb)
	}
	return symbols, symbolsStructure, nil
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
		if equation.LeftPart.Length > 0 {
			for _, sym := range equation.LeftPart.Symbols {
				if symbol.IsLetter(sym) || symbol.IsConst(sym) {
					return true
				}
			}
		}
	}
	if equation.IsLeftEmpty() {
		if equation.RightPart.Length > 0 {
			for _, sym := range equation.RightPart.Symbols {
				if symbol.IsLetter(sym) || symbol.IsConst(sym) {
					return true
				}
			}
		}
	}
	return false
}

func (equation *Equation) CheckEquality() bool {
	if equation.RightPart.Length == 0 {
		equation.RightPart.Length++
		equation.RightPart.Symbols = append(equation.RightPart.Symbols, symbol.Empty())
	}
	i := 0
	for _, sym := range equation.LeftPart.Symbols {
		if symbol.IsEmpty(sym) {

		} else {
			for i < equation.RightPart.Length && symbol.IsEmpty(equation.RightPart.Symbols[i]) {
				i++
			}
			if i == equation.RightPart.Length {
				return false
			}
			if sym != equation.RightPart.Symbols[i] {
				return false
			} else {
				i++
			}
		}
	}
	for i < equation.RightPart.Length && symbol.IsEmpty(equation.RightPart.Symbols[i]) {
		i++
	}
	if i != equation.RightPart.Length {
		return false
	}
	return true
}

func (equation *Equation) CheckSameness(eq *Equation) bool {
	var wordsMap = map[string]string{}
	if eq.RightPart.Length == 0 {
		eq.RightPart.Length++
		eq.RightPart.Symbols = append(eq.RightPart.Symbols, symbol.Empty())
	}
	if eq.LeftPart.Length == 0 {
		eq.LeftPart.Length++
		eq.LeftPart.Symbols = append(eq.LeftPart.Symbols, symbol.Empty())
	}
	i := 0
	for _, sym := range equation.LeftPart.Symbols {
		if symbol.IsEmpty(sym) {

		} else {
			for i < eq.LeftPart.Length && symbol.IsEmpty(eq.LeftPart.Symbols[i]) {
				i++
			}
			if i == eq.LeftPart.Length {
				return false
			}
			if sym != eq.LeftPart.Symbols[i] {
				if symbol.IsLetter(sym) && symbol.IsLetter(eq.LeftPart.Symbols[i]) {
					if wordsMap[sym.Value()] == "" {
						wordsMap[sym.Value()] = eq.LeftPart.Symbols[i].Value()
						i++
					} else {
						if wordsMap[sym.Value()] != eq.LeftPart.Symbols[i].Value() {
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
	for i < eq.LeftPart.Length && symbol.IsEmpty(eq.LeftPart.Symbols[i]) {
		i++
	}
	if i != eq.LeftPart.Length {
		return false
	}
	i = 0
	for _, sym := range equation.RightPart.Symbols {
		if symbol.IsEmpty(sym) {

		} else {
			for i < eq.RightPart.Length && symbol.IsEmpty(eq.RightPart.Symbols[i]) {
				i++
			}
			if i == eq.RightPart.Length {
				return false
			}
			if sym != eq.RightPart.Symbols[i] {
				if symbol.IsLetter(sym) && symbol.IsLetter(eq.RightPart.Symbols[i]) {
					if wordsMap[sym.Value()] == "" {
						wordsMap[sym.Value()] = eq.RightPart.Symbols[i].Value()
						i++
					} else {
						if wordsMap[sym.Value()] != eq.RightPart.Symbols[i].Value() {
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
	for i < eq.RightPart.Length && symbol.IsEmpty(eq.RightPart.Symbols[i]) {
		i++
	}
	if i != eq.RightPart.Length {
		return false
	}
	return true
}

func (equation *Equation) SubstituteVarsWithEmpty() (Equation, map[symbol.Symbol]bool) {
	var varsMap = make(map[symbol.Symbol]bool)
	var resultEquation Equation
	resultEquation.New()
	if equation.IsRightEmpty() {
		resultEquation.RightPart = equation.RightPart
		resultEquation.RightPart.Length = equation.RightPart.Length
		for _, sym := range equation.LeftPart.Symbols {
			if symbol.IsConst(sym) || symbol.IsLetter(sym) {
				resultEquation.LeftPart.Symbols = append(resultEquation.LeftPart.Symbols, sym)
				resultEquation.LeftPart.Structure.Add(sym)
			} else {
				varsMap[sym] = true
			}
		}
		if len(resultEquation.LeftPart.Symbols) == 0 {
			resultEquation.LeftPart.Symbols = append(resultEquation.LeftPart.Symbols, symbol.Empty())
		}
		resultEquation.LeftPart.Length = len(resultEquation.LeftPart.Symbols)
	}
	if equation.IsLeftEmpty() {
		resultEquation.LeftPart = equation.LeftPart
		resultEquation.LeftPart.Length = equation.LeftPart.Length
		for _, sym := range equation.RightPart.Symbols {
			if symbol.IsConst(sym) || symbol.IsLetter(sym) {
				resultEquation.RightPart.Symbols = append(resultEquation.RightPart.Symbols, sym)
				resultEquation.RightPart.Structure.Add(sym)
			} else {
				varsMap[sym] = true
			}
		}
		if len(resultEquation.RightPart.Symbols) == 0 {
			resultEquation.RightPart.Symbols = append(resultEquation.RightPart.Symbols, symbol.Empty())
		}
		resultEquation.RightPart.Length = len(resultEquation.RightPart.Symbols)
	}
	return resultEquation, varsMap
}

func (equation *Equation) Substitute(substitution Substitution) Equation {
	newSymLen := substitution.RightPartLength()
	var resultEquation Equation
	resultEquation.New()
	var lCounter, rCounter int
	for _, sym := range equation.LeftPart.Symbols {
		if sym == substitution.LeftPart() {
			lCounter++
			resultEquation.LeftPart.Symbols = append(resultEquation.LeftPart.Symbols, substitution.RightPart()...)
			resultEquation.LeftPart.Length += newSymLen
		} else {
			resultEquation.LeftPart.Symbols = append(resultEquation.LeftPart.Symbols, sym)
			resultEquation.LeftPart.Length++
			resultEquation.LeftPart.Structure.Add(sym)
		}
	}
	for _, sym := range equation.RightPart.Symbols {
		if sym == substitution.LeftPart() {
			rCounter++
			resultEquation.RightPart.Symbols = append(resultEquation.RightPart.Symbols, substitution.RightPart()...)
			resultEquation.RightPart.Length += newSymLen
		} else {
			resultEquation.RightPart.Symbols = append(resultEquation.RightPart.Symbols, sym)
			resultEquation.RightPart.Length++
			resultEquation.RightPart.Structure.Add(sym)
		}
	}
	for _, sym := range substitution.RightPart() {
		resultEquation.LeftPart.Structure.AddTimes(sym, lCounter)
		resultEquation.RightPart.Structure.AddTimes(sym, rCounter)
	}
	resultEquation.Reduce()
	return resultEquation
}

func (equation *Equation) Reduce() {
	equation.reduceEmpty()
	minLen := standart.Min(equation.LeftPart.Length, equation.RightPart.Length)
	i := 0
	for ; i < minLen; i++ {
		if equation.LeftPart.Symbols[i] != equation.RightPart.Symbols[i] {
			break
		}
		equation.LeftPart.Structure.Sub(equation.LeftPart.Symbols[i])
		equation.RightPart.Structure.Sub(equation.RightPart.Symbols[i])
	}
	if i > 0 {
		equation.RightPart.Symbols = equation.RightPart.Symbols[i:]
		equation.RightPart.Length -= i
		equation.LeftPart.Symbols = equation.LeftPart.Symbols[i:]
		equation.LeftPart.Length -= i
	}
	equation.reduceEmpty()
}

func (equation *Equation) reduceEmpty() {
	i := 0
	if equation.LeftPart.Length > 1 {
		for ; i < equation.LeftPart.Length; i++ {
			if !symbol.IsEmpty(equation.LeftPart.Symbols[i]) {
				break
			}
		}
		if i > 0 {
			equation.LeftPart.Symbols = equation.LeftPart.Symbols[i:]
			equation.LeftPart.Length -= i
		}
	}
	if equation.RightPart.Length > 1 {
		i = 0
		for ; i < equation.RightPart.Length; i++ {
			if !symbol.IsEmpty(equation.RightPart.Symbols[i]) {
				break
			}
		}
		if i > 0 {
			equation.RightPart.Symbols = equation.RightPart.Symbols[i:]
			equation.RightPart.Length -= i
		}
	}
}

func (equation *Equation) SplitByEquidecomposability() EqSystem {
	equation.isEquidecomposable = checkEquidecomposability(equation.LeftPart.Symbols, equation.RightPart.Symbols)
	defaultSystem := EqSystem{
		Equations: []Equation{*equation}, Size: 1,
	}
	if equation.LeftPart.Length <= 1 || equation.RightPart.Length <= 1 {
		return defaultSystem
	}
	var firstPart, secondPart []symbol.Symbol
	var firstPartEnd, secondPartEnd []symbol.Symbol
	minLen := standart.Min(equation.LeftPart.Length, equation.RightPart.Length)
	// forward order
	for i := 1; i < minLen; i++ {
		firstPart = equation.LeftPart.Symbols[:i]
		secondPart = equation.RightPart.Symbols[:i]
		firstPartEnd = equation.LeftPart.Symbols[i:]
		secondPartEnd = equation.RightPart.Symbols[i:]
		if checkEquidecomposability(firstPart, secondPart) {
			return createSystem(firstPart, secondPart, firstPartEnd, secondPartEnd)
		}
	}
	// backwards order
	for i := 1; i < minLen; i++ {
		firstPart = equation.LeftPart.Symbols[minLen-i:]
		secondPart = equation.RightPart.Symbols[minLen-i:]
		firstPartEnd = equation.LeftPart.Symbols[:minLen-i]
		secondPartEnd = equation.RightPart.Symbols[:minLen-i]
		if checkEquidecomposability(firstPart, secondPart) {
			return createSystem(firstPart, secondPart, firstPartEnd, secondPartEnd)
		}
	}

	return defaultSystem
}

func createSystem(firstPart []symbol.Symbol, secondPart []symbol.Symbol, firstPartEnd []symbol.Symbol, secondPartEnd []symbol.Symbol) EqSystem {
	var fLeftPart, fRightPart EqPart
	var sLeftPart, sRightPart EqPart
	fLeftPart.NewFromSymbols(firstPart)
	fRightPart.NewFromSymbols(secondPart)
	sLeftPart.NewFromSymbols(firstPartEnd)
	sRightPart.NewFromSymbols(secondPartEnd)
	resultSystem := EqSystem{
		Equations: []Equation{
			{
				LeftPart:           fLeftPart,
				RightPart:          fRightPart,
				isEquidecomposable: true,
			},
		},
		Size: 1,
	}
	endEq := Equation{
		LeftPart:  sLeftPart,
		RightPart: sRightPart,
	}
	endSystem := endEq.SplitByEquidecomposability()
	resultSystem.Equations = append(endSystem.Equations, resultSystem.Equations...)
	resultSystem.Size += endSystem.Size
	return resultSystem
}

func checkEquidecomposability(firstPart []symbol.Symbol, secondPart []symbol.Symbol) bool {
	if len(firstPart) != len(secondPart) {
		return false
	}
	var length = len(firstPart)
	var firstVars = make(map[symbol.Symbol]int)
	var secondVars = make(map[symbol.Symbol]int)
	var firstConsts = make([]symbol.Symbol, 0)
	var secondConsts = make([]symbol.Symbol, 0)
	var isEquidecomposable = true
	var firstSym, secondSym symbol.Symbol
	for i := 0; i < length; i++ {
		firstSym = firstPart[i]
		if symbol.IsVar(firstSym) || symbol.IsLetter(firstSym) {
			firstVars[firstSym] += 1
		} else if symbol.IsConst(firstSym) {
			firstConsts = append(firstConsts, firstSym)
		}
		secondSym = secondPart[i]
		if symbol.IsVar(secondSym) || symbol.IsLetter(secondSym) {
			secondVars[secondSym] += 1
		} else if symbol.IsConst(secondSym) {
			secondConsts = append(secondConsts, secondSym)
		}
	}
	// if equation doesn't have vars or letters, checking parts are equal
	if len(firstVars) == 0 && len(secondVars) == 0 {
		if len(firstConsts) != len(secondConsts) {
			return false
		}
		for lS, lN := range firstConsts {
			if lN != secondConsts[lS] {
				isEquidecomposable = false
				break
			}
		}
	}
	if len(firstVars) != len(secondVars) {
		return false
	}
	for lS, lN := range firstVars {
		if lN != secondVars[lS] {
			isEquidecomposable = false
			break
		}
	}
	if isEquidecomposable {
		return true
	}
	return false
}

func (equation *Equation) IsQuadratic() bool {
	var helpMap = make(map[symbol.Symbol]int)
	for lS, lN := range equation.LeftPart.Structure.Vars() {
		if lN > 2 {
			return false
		}
		helpMap[lS] += lN
	}
	for lS, lN := range equation.RightPart.Structure.Vars() {
		helpMap[lS] += lN
		if helpMap[lS] > 2 {
			return false
		}
	}
	return true
}

func (equation *Equation) IsRegularlyOrdered() bool {
	var i, j int
	var symL, symR symbol.Symbol
	for {
		for ; i < equation.LeftPart.Length; i++ {
			symL = equation.LeftPart.Symbols[i]
			if !symbol.IsConst(symL) && !symbol.IsEmpty(symL) {
				break
			}
		}
		for ; j < equation.RightPart.Length; j++ {
			symR = equation.RightPart.Symbols[j]
			if !symbol.IsConst(symR) && !symbol.IsEmpty(symR) {
				break
			}
		}
		if i == equation.LeftPart.Length {
			break
		}
		if j == equation.RightPart.Length {
			break
		}
		symL = equation.LeftPart.Symbols[i]
		symR = equation.RightPart.Symbols[j]
		if symL.Value() != symR.Value() {
			return false
		}
		i++
		j++
	}
	for ; i < equation.LeftPart.Length; i++ {
		symL = equation.LeftPart.Symbols[i]
		if symbol.IsLetterOrVar(symL) {
			return false
		}
	}
	for ; j < equation.RightPart.Length; j++ {
		symR = equation.RightPart.Symbols[j]
		if symbol.IsLetterOrVar(symR) {
			return false
		}
	}
	if i == equation.LeftPart.Length && j == equation.RightPart.Length {
		return true
	}
	return false
}

func (equation *Equation) IsEquidecomposable() bool {
	return equation.isEquidecomposable
}

func (equation *Equation) HasEqVarsLen() bool {
	if equation.LeftPart.Structure.VarsRangeLen() != equation.RightPart.Structure.VarsRangeLen() {
		return false
	}
	return true
}

func (equation *Equation) HasEqConstsLen() bool {
	if equation.LeftPart.Structure.ConstsRangeLen() != equation.RightPart.Structure.ConstsRangeLen() {
		return false
	}
	return true
}

func (equation *Equation) IsLeftEmpty() bool {
	return equation.LeftPart.Length == 1 &&
		symbol.IsEmpty(equation.LeftPart.Symbols[0]) || equation.LeftPart.Length == 0
}

func (equation *Equation) IsRightEmpty() bool {
	return equation.RightPart.Length == 1 &&
		symbol.IsEmpty(equation.RightPart.Symbols[0]) || equation.RightPart.Length == 0
}

func (equation *Equation) Print() {
	fmt.Println(equation.String())
}

func (equation *Equation) String() string {
	var result string
	for _, sym := range equation.LeftPart.Symbols {
		result += fmt.Sprintf("%s ", sym.Value())
	}
	result += fmt.Sprintf("%s ", EQUALS)
	for _, sym := range equation.RightPart.Symbols {
		result += fmt.Sprintf("%s ", sym.Value())
	}
	return result
}
