package solver

import (
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"math"
)

func checkLengthRules(eq *equation.Equation) (bool, symbol.Symbol, int) {
	if eq.IsEquidecomposable() {
		return analiseMultiplicity(eq), nil, 0
	}
	var diffL, diffR int
	var diffVL, diffVR float64
	var diffSymL, diffSymR symbol.Symbol
	s1 := eq.LeftPart.Structure.LettersAndConstsLen()
	s2 := eq.RightPart.Structure.LettersAndConstsLen()
	var hasBeenMap = make(map[symbol.Symbol]bool)
	var leftMap, rightMap map[symbol.Symbol]int
	leftMap = eq.LeftPart.Structure.Vars()
	rightMap = eq.RightPart.Structure.Vars()
	var cont, res bool
	for sym, numL := range leftMap {
		numR := rightMap[sym]
		if numR != 0 {
			hasBeenMap[sym] = true
		}
		cont, res = checkRules(sym, s1, s2, numL, numR, &diffL, &diffR, &diffSymL, &diffSymR, &diffVL, &diffVR)
		if !cont {
			return res, nil, 0
		}
	}
	for sym, numR := range rightMap {
		if hasBeenMap[sym] {
			continue
		}
		numL := leftMap[sym]
		cont, res = checkRules(sym, s1, s2, numL, numR, &diffL, &diffR, &diffSymL, &diffSymR, &diffVL, &diffVR)
		if !cont {
			return res, nil, 0
		}
	}
	if diffL == 1 && s2 >= s1 {
		var newLetters = float64(s2-s1) / diffVL
		if newLetters-math.Trunc(newLetters) == 0 {
			return true, diffSymL, int(newLetters)
		}
		return false, nil, 0
	}
	if diffR == 1 && s2 <= s1 {
		var newLetters = float64(s1-s2) / diffVR
		if newLetters-math.Trunc(newLetters) == 0 {
			return true, diffSymR, int(newLetters)
		}
		return false, nil, 0
	}
	if diffR == 0 && diffL == 0 && s1 == s2 {
		return analiseMultiplicity(eq), nil, 0
	}
	return true, nil, 0
}

func checkRules(sym symbol.Symbol, s1, s2, numL, numR int, diffL, diffR *int, diffSymL, diffSymR *symbol.Symbol, diffVL, diffVR *float64) (bool, bool) {
	// TODO: check this case
	if s1 == s2 && (*diffL > 1 || *diffR > 1) {
		return false, true
	}
	if numL > numR {
		*diffVL = float64(numL - numR)
		*diffL++
		*diffSymL = sym
		if s1 > s2 {
			return false, false
		}
	} else if numL < numR {
		*diffVR = float64(numR - numL)
		*diffR++
		*diffSymR = sym
		if s2 > s1 {
			return false, false
		}
	}
	return true, true
}

func analiseMultiplicity(eq *equation.Equation) bool {
	r1 := eq.LeftPart.Structure.LettersLen()
	r2 := eq.RightPart.Structure.LettersLen()
	var leftMap, rightMap map[symbol.Symbol]int
	rightMap = eq.RightPart.Structure.Consts()
	leftMap = eq.LeftPart.Structure.Consts()
	var hasBeenMap = make(map[symbol.Symbol]bool)
	var leftLen, rightLen int
	for sym, numL := range leftMap {
		numR := rightMap[sym]
		if numR != 0 {
			hasBeenMap[sym] = true
		}
		if numL > numR {
			leftLen += numL - numR
		} else if numR > numL {
			rightLen += numR - numL
		}
	}
	for sym, numR := range rightMap {
		if hasBeenMap[sym] {
			continue
		}
		numL := leftMap[sym]
		if numL > numR {
			leftLen += numL - numR
		} else if numR > numL {
			rightLen += numR - numL
		}
	}
	if leftLen > r2 {
		return false
	}
	if rightLen > r1 {
		return false
	}
	return true
}
