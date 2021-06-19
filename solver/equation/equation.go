package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/standart"
)

type Equation struct {
	LeftPart           EqPart
	RightPart          EqPart
	structure          Structure
	isEquidecomposable bool
	isRegularlyOrdered bool
}

const EQUALS = "="

func (equation *Equation) Structure() *Structure {
	return &equation.structure
}

func (equation *Equation) Letters() []symbol.Symbol {
	var letters = make([]symbol.Symbol, 0)
	for s := range equation.structure.letters {
		letters = append(letters, s)
	}
	return letters
}

func (equation *Equation) Consts() []symbol.Symbol {
	var consts = make([]symbol.Symbol, 0)
	for s := range equation.structure.consts {
		consts = append(consts, s)
	}
	return consts
}

func (equation *Equation) Vars() []symbol.Symbol {
	var vars = make([]symbol.Symbol, 0)
	for s := range equation.structure.vars {
		vars = append(vars, s)
	}
	return vars
}

func (equation *Equation) HasVarOrLetter(s symbol.Symbol) bool {
	return equation.structure.vars[s] != 0 || equation.structure.letters[s] != 0
}

func (equation *Equation) HasVarOrLetterForNormal(s symbol.Symbol) (bool, error) {
	if symbol.IsVar(s) {
		if !(equation.LeftPart.Structure.vars[s] == 1 &&
			equation.RightPart.Structure.vars[s] == 1) {
			return false, nil
		}
	} else if symbol.IsLetter(s) {
		if !(equation.LeftPart.Structure.letters[s] == 1 &&
			equation.RightPart.Structure.letters[s] == 1) {
			return false, nil
		}
	} else {
		return false, fmt.Errorf("symbol is not var or letter: %v", s)
	}

	// checking no more variables in f1 and f2
	if (equation.structure.varsLen + equation.structure.lettersLen) != 2 {
		return false, nil
	}
	var f1, f2 []symbol.Symbol
	if equation.LeftPart.Symbols[0] == s && equation.RightPart.Symbols[equation.RightPart.Length-1] == s {
		f1 = equation.LeftPart.Symbols[1:]
		f2 = equation.RightPart.Symbols[:equation.RightPart.Length-1]
	} else if equation.RightPart.Symbols[0] == s && equation.LeftPart.Symbols[equation.LeftPart.Length-1] == s {
		f2 = equation.RightPart.Symbols[1:]
		f1 = equation.LeftPart.Symbols[:equation.LeftPart.Length-1]
	} else {
		return false, nil
	}

	return checkSimpleWord(f1) && checkSimpleWord(f2), nil
}

func checkSimpleWord(s []symbol.Symbol) bool {
	var sLen = len(s)
	if sLen <= 1 {
		return true
	}
	var currPrefix []symbol.Symbol
	var w = s[:1]
	var wLen = len(w)
	var equal bool
	var newI int

	for k := 1; k < sLen; {
		if k+wLen > sLen {
			return true
		}
		currPrefix = s[k : k+wLen]
		equal = standart.CheckSymbolArraysEquality(w, currPrefix)
		if !equal {
			newI = k + 1
			if newI > sLen {
				return true
			}
			w = s[:newI]
			wLen = len(w)
			k++
		} else {
			k += wLen
		}
	}
	return !equal
}

func (equation *Equation) Check(constsAlphabet *Alphabet, varsAlphabet *Alphabet, lettersAlphabet *Alphabet) error {
	for _, s := range equation.LeftPart.Symbols {
		if symbol.IsVar(s) && !varsAlphabet.Has(s.Value()) {
			return fmt.Errorf("variable doesn't belong to alphabet: %v", s)
		}
		if symbol.IsConst(s) && !constsAlphabet.Has(s.Value()) {
			return fmt.Errorf("const doesn't belong to alphabet: %v", s)
		}
		if symbol.IsLetter(s) && !lettersAlphabet.Has(s.Value()) {
			return fmt.Errorf("letter doesn't belong to alphabet: %v", s)
		}
	}
	return nil
}

func (equation *Equation) IsEmpty() bool {
	return equation.LeftPart.IsEmpty() && equation.RightPart.IsEmpty()
}

func (equation *Equation) Copy() Equation {
	newEq := Equation{}
	newEq.LeftPart = equation.LeftPart.Copy()
	newEq.RightPart = equation.RightPart.Copy()
	newEq.isEquidecomposable = equation.isEquidecomposable
	newEq.structure = equation.structure.Copy()
	return newEq
}

func NewEquation(leftPart []symbol.Symbol, rightPart []symbol.Symbol) Equation {
	var eq Equation
	eq.NewFromParts(leftPart, rightPart)
	return eq
}

func (equation *Equation) New() {
	equation.LeftPart = EmptyEqPart()
	equation.RightPart = EmptyEqPart()
	equation.structure = EmptyStructure()
}

func (equation *Equation) NewFromParts(leftPart []symbol.Symbol, rightPart []symbol.Symbol) {
	equation.LeftPart = NewEqPartFromSymbols(leftPart)
	equation.RightPart = NewEqPartFromSymbols(rightPart)
	equation.structure = MergeStructures(&equation.LeftPart.Structure, &equation.RightPart.Structure)
}

func (equation *Equation) Init(eq string, constAlphabet *Alphabet, varsAlphabet *Alphabet) error {
	//fmt.Println(constAlphabet.words)
	//fmt.Print(varsAlphabet.words)
	var err error
	isEq, i := checkEquation(eq)
	if !isEq {
		return fmt.Errorf("invalid equation: %s", eq)

	}
	equation.structure = EmptyStructure()
	eqleftPart := eq[0 : i-1]
	var leftSymbols []symbol.Symbol
	var leftSymbolsStruct = EmptyStructure()
	if eqleftPart == "" {
		leftSymbols = append(leftSymbols, symbol.Empty())
	} else {
		leftSymbols, err = matchWithAlphabetsWithSpace(eqleftPart, constAlphabet, varsAlphabet,
			&leftSymbolsStruct, &equation.structure)
		if err != nil {
			return fmt.Errorf("error matching alphabet: %v", err)
		}
	}
	eqRightPart := eq[i+2:]
	var rightSymbols []symbol.Symbol
	var rightSymbolsStruct = EmptyStructure()
	if eqRightPart == "" {
		rightSymbols = append(rightSymbols, symbol.Empty())
	} else {
		rightSymbols, err = matchWithAlphabetsWithSpace(eqRightPart, constAlphabet, varsAlphabet,
			&rightSymbolsStruct, &equation.structure)
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

func checkEquation(eq string) (bool, int) {
	eqLen := len(eq) - 1
	for i := 0; i < eqLen; i++ {
		if string(eq[i]) == EQUALS && string(eq[i-1]) == SPACE && string(eq[i+1]) == SPACE {
			return true, i
		}
	}
	return false, 0
}

func matchWithAlphabetsWithSpace(eqPart string, constAlphabet *Alphabet,
	varsAlphabet *Alphabet, symbolsStructure *Structure, equationStructure *Structure) ([]symbol.Symbol, error) {
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
			symbolsStructure.Add(symb)
			equationStructure.Add(symb)
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
		symbolsStructure.Add(symb)
		equationStructure.Add(symb)
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
				resultEquation.structure.Add(sym)
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
				resultEquation.structure.Add(sym)
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
			resultEquation.structure.Add(sym)
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
			resultEquation.structure.Add(sym)
		}
	}
	for _, sym := range substitution.RightPart() {
		resultEquation.LeftPart.Structure.AddTimes(sym, lCounter)
		resultEquation.RightPart.Structure.AddTimes(sym, rCounter)
		resultEquation.structure.AddTimes(sym, lCounter+rCounter)
	}
	resultEquation.Reduce()
	return resultEquation
}

func (equation *Equation) Reduce() bool {
	var reduced bool
	i := 0
	j := 0
	for i < equation.LeftPart.Length && j < equation.RightPart.Length {
		if symbol.IsEmpty(equation.LeftPart.Symbols[i]) {
			i++
			continue
		}
		if symbol.IsEmpty(equation.RightPart.Symbols[j]) {
			j++
			continue
		}
		if equation.LeftPart.Symbols[i] != equation.RightPart.Symbols[j] {
			break
		}
		s := equation.LeftPart.Symbols[i]
		equation.LeftPart.Structure.Sub(s)
		equation.RightPart.Structure.Sub(s)
		equation.structure.SubTimes(s, 2)
		i++
		j++
	}
	if i > 0 && j > 0 {
		equation.RightPart.Symbols = equation.RightPart.Symbols[j:]
		equation.RightPart.Length -= j
		equation.LeftPart.Symbols = equation.LeftPart.Symbols[i:]
		equation.LeftPart.Length -= i
		reduced = true
	}
	i = equation.LeftPart.Length - 1
	j = equation.RightPart.Length - 1
	for i >= 0 && j >= 0 {
		if symbol.IsEmpty(equation.LeftPart.Symbols[i]) {
			i--
			continue
		}
		if symbol.IsEmpty(equation.RightPart.Symbols[j]) {
			j--
			continue
		}
		if equation.LeftPart.Symbols[i] != equation.RightPart.Symbols[j] {
			break
		}
		s := equation.LeftPart.Symbols[i]
		equation.LeftPart.Structure.Sub(s)
		equation.RightPart.Structure.Sub(s)
		equation.structure.SubTimes(s, 2)
		i--
		j--
	}
	if i != equation.LeftPart.Length-1 && j != equation.RightPart.Length-1 {
		equation.RightPart.Symbols = equation.RightPart.Symbols[:j+1]
		equation.RightPart.Length -= j
		equation.LeftPart.Symbols = equation.LeftPart.Symbols[:i+1]
		equation.LeftPart.Length -= i
		reduced = true
	}
	equation.FullReduceEmpty()
	return reduced
}

func (equation *Equation) FullReduceEmpty() {
	var equationLeftPart = make([]symbol.Symbol, 0)
	var equationRightPart = make([]symbol.Symbol, 0)
	if equation.LeftPart.Length > 1 {
		for _, sym := range equation.LeftPart.Symbols {
			if !symbol.IsEmpty(sym) {
				equationLeftPart = append(equationLeftPart, sym)
			}
		}
		equation.LeftPart.Symbols = equationLeftPart
		equation.LeftPart.Length = len(equationLeftPart)
	}
	if equation.RightPart.Length > 1 {
		for _, sym := range equation.RightPart.Symbols {
			if !symbol.IsEmpty(sym) {
				equationRightPart = append(equationRightPart, sym)
			}
		}
		equation.RightPart.Symbols = equationRightPart
		equation.RightPart.Length = len(equationRightPart)
	}
}

func (equation *Equation) SplitByEquidecomposability() EqSystem {
	equation.isEquidecomposable, _ = checkEquidecomposability(equation.LeftPart.Symbols, equation.RightPart.Symbols)
	defaultSystem := SystemFromEquation(*equation)
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
		equideComposable, equal := checkEquidecomposability(firstPart, secondPart)
		if equideComposable {
			if equal {
				return createSystemEqual(firstPartEnd, secondPartEnd)
			} else {
				return createSystem(firstPart, secondPart, firstPartEnd, secondPartEnd)
			}
		}
	}
	// backwards order
	for i := 1; i < minLen; i++ {
		firstPart = equation.LeftPart.Symbols[minLen-i:]
		secondPart = equation.RightPart.Symbols[minLen-i:]
		firstPartEnd = equation.LeftPart.Symbols[:minLen-i]
		secondPartEnd = equation.RightPart.Symbols[:minLen-i]
		equideComposable, equal := checkEquidecomposability(firstPart, secondPart)
		if equideComposable {
			if equal {
				return createSystemEqual(firstPartEnd, secondPartEnd)
			} else {
				return createSystem(firstPart, secondPart, firstPartEnd, secondPartEnd)
			}
		}
	}

	return defaultSystem
}

func createSystemEqual(firstPartEnd []symbol.Symbol, secondPartEnd []symbol.Symbol) EqSystem {
	endEq := NewEquation(firstPartEnd, secondPartEnd)
	endSystem := endEq.SplitByEquidecomposability()
	return endSystem
}

func createSystem(firstPart []symbol.Symbol, secondPart []symbol.Symbol, firstPartEnd []symbol.Symbol, secondPartEnd []symbol.Symbol) EqSystem {
	resultSystem := SystemFromEquation(NewEquation(firstPart, secondPart))
	endEq := NewEquation(firstPartEnd, secondPartEnd)
	endSystem := endEq.SplitByEquidecomposability()
	resultSystem.Equations = append(endSystem.Equations, resultSystem.Equations...)
	resultSystem.Size += endSystem.Size
	return resultSystem
}

func (equation *Equation) CheckEquidecomposability() bool {
	if equation.isEquidecomposable {
		return true
	} else {
		equation.isEquidecomposable, _ = checkEquidecomposability(equation.LeftPart.Symbols, equation.RightPart.Symbols)
	}
	return equation.isEquidecomposable
}

func checkEquidecomposability(firstPart []symbol.Symbol, secondPart []symbol.Symbol) (bool, bool) {
	if len(firstPart) != len(secondPart) {
		return false, false
	}
	var length = len(firstPart)
	var firstVars = make(map[symbol.Symbol]int)
	var secondVars = make(map[symbol.Symbol]int)
	var firstConsts = make([]symbol.Symbol, 0)
	var secondConsts = make([]symbol.Symbol, 0)
	var isEquidecomposable = true
	var isEqual = true
	var firstSym, secondSym symbol.Symbol
	for i := 0; i < length; i++ {
		firstSym = firstPart[i]
		secondSym = secondPart[i]
		if isEqual && firstSym != secondSym {
			isEqual = false
		}
		if symbol.IsVar(firstSym) || symbol.IsLetter(firstSym) {
			firstVars[firstSym] += 1
		} else if symbol.IsConst(firstSym) {
			firstConsts = append(firstConsts, firstSym)
		}

		if symbol.IsVar(secondSym) || symbol.IsLetter(secondSym) {
			secondVars[secondSym] += 1
		} else if symbol.IsConst(secondSym) {
			secondConsts = append(secondConsts, secondSym)
		}
	}
	if isEqual {
		return true, true
	}
	// if equation doesn't have vars or letters, checking parts are equal
	if len(firstVars) == 0 && len(secondVars) == 0 {
		if len(firstConsts) != len(secondConsts) {
			return false, false
		}
		for lS, lN := range firstConsts {
			if lN != secondConsts[lS] {
				return false, false
			}
		}
		return true, true
	}
	if len(firstVars) != len(secondVars) {
		return false, false
	}
	for lS, lN := range firstVars {
		if lN != secondVars[lS] {
			isEquidecomposable = false
			break
		}
	}
	if isEquidecomposable {
		return true, false
	}
	return false, false
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
			if symbol.IsVar(symL) {
				break
			}
		}
		for ; j < equation.RightPart.Length; j++ {
			symR = equation.RightPart.Symbols[j]
			if symbol.IsVar(symR) {
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
		if symbol.IsVar(symL) {
			return false
		}
	}
	for ; j < equation.RightPart.Length; j++ {
		symR = equation.RightPart.Symbols[j]
		if symbol.IsVar(symR) {
			return false
		}
	}
	if i == equation.LeftPart.Length && j == equation.RightPart.Length {
		equation.isRegularlyOrdered = true
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
