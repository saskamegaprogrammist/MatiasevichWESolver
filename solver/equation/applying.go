package equation

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
)

func (equation *Equation) Apply(e Equation) (bool, Equation, error) {
	var err error
	if !e.isEquidecomposable {
		return false, Equation{}, fmt.Errorf("applied equation %s is not equidecomposable", e.String())
	}
	if equation.Structure().Size() < e.structure.Size() {
		return false, Equation{}, nil
	}
	var newEq Equation
	appliedFirstForward, newEq, err := equation.applyFirstRule(e, FORWARD)
	if err != nil {
		return appliedFirstForward, newEq, fmt.Errorf("error applying first rule forward: %v", err)
	}
	if appliedFirstForward {
		return appliedFirstForward, newEq, nil
	}

	appliedFirstBackwards, newEq, err := equation.applyFirstRule(e, BACKWARDS)
	if err != nil {
		return appliedFirstBackwards, newEq, fmt.Errorf("error applying first rule backwards: %v", err)
	}
	if appliedFirstBackwards {
		return appliedFirstBackwards, newEq, nil
	}
	appliedThirdForward, newEq, err := equation.applyThirdRule(e, FORWARD)
	if err != nil {
		return appliedThirdForward, newEq, fmt.Errorf("error applying third rule backwards: %v", err)
	}
	if appliedThirdForward {
		return appliedThirdForward, newEq, nil
	}
	appliedThirdBackwards, newEq, err := equation.applyThirdRule(e, BACKWARDS)
	if err != nil {
		return appliedThirdBackwards, newEq, fmt.Errorf("error applying third rule backwards: %v", err)
	}
	return appliedThirdBackwards, newEq, nil
}

func (equation *Equation) applyFirstRule(e Equation, mode int) (bool, Equation, error) {
	var i, j int
	var err error
	var lSym, rSym, lSymE, rSymE symbol.Symbol
	err = findFirstDifferentFirstRuleWLeftAndLeft(equation, &e, &i, &lSym, &lSymE, mode)
	if err != nil {
		return false, Equation{}, fmt.Errorf("error finding first different: %v", err)
	}
	if i != 0 {
		err = findFirstDifferentFirstRuleWRightAndRight(equation, &e, &j, &rSym, &rSymE, mode)
		if err != nil {
			return false, Equation{}, fmt.Errorf("error finding first different: %v", err)
		}
		if j == 0 {
			return false, Equation{}, nil
		} else {
			if i > j {
				if i != e.LeftPart.Length {
					return false, Equation{}, nil
				}

				var leftSymbolsFirst, leftSymbolsSecond []symbol.Symbol

				if mode == FORWARD {
					leftSymbolsFirst = make([]symbol.Symbol, len(e.RightPart.Symbols))
					copy(leftSymbolsFirst, e.RightPart.Symbols)
					leftSymbolsSecond = make([]symbol.Symbol, len(equation.LeftPart.Symbols[i:]))
					copy(leftSymbolsSecond, equation.LeftPart.Symbols[i:])
				} else {
					leftSymbolsFirst = make([]symbol.Symbol, len(equation.LeftPart.Symbols[:equation.LeftPart.Length-i]))
					copy(leftSymbolsFirst, equation.LeftPart.Symbols[:equation.LeftPart.Length-i])
					var leftSymbolsSecond = make([]symbol.Symbol, len(e.RightPart.Symbols))
					copy(leftSymbolsSecond, e.RightPart.Symbols)
				}

				var rightSymbols = make([]symbol.Symbol, len(equation.RightPart.Symbols))
				copy(rightSymbols, equation.RightPart.Symbols)
				newEq := NewEquation(append(leftSymbolsFirst, leftSymbolsSecond...), rightSymbols)
				newEq.Reduce()
				return false, newEq, nil
			} else {
				if j != e.RightPart.Length {
					return false, Equation{}, nil
				}

				var rightSymbolsFirst, rightSymbolsSecond []symbol.Symbol

				if mode == FORWARD {
					rightSymbolsFirst = make([]symbol.Symbol, len(e.LeftPart.Symbols))
					copy(rightSymbolsFirst, e.LeftPart.Symbols)
					rightSymbolsSecond = make([]symbol.Symbol, len(equation.RightPart.Symbols[j:]))
					copy(rightSymbolsSecond, equation.RightPart.Symbols[j:])
				} else {
					rightSymbolsFirst = make([]symbol.Symbol, len(equation.RightPart.Symbols[:equation.RightPart.Length-j]))
					copy(rightSymbolsFirst, equation.RightPart.Symbols[:equation.RightPart.Length-j])
					rightSymbolsSecond = make([]symbol.Symbol, len(e.LeftPart.Symbols))
					copy(rightSymbolsSecond, e.LeftPart.Symbols)
				}

				var leftSymbols = make([]symbol.Symbol, len(equation.LeftPart.Symbols))
				copy(leftSymbols, equation.LeftPart.Symbols)
				newEq := NewEquation(leftSymbols, append(rightSymbolsFirst, rightSymbolsSecond...))
				newEq.Reduce()
				return false, newEq, nil
			}
		}
	}
	err = findFirstDifferentFirstRuleWLeftAndRight(equation, &e, &i, &lSym, &rSymE, mode)
	if err != nil {
		return false, Equation{}, fmt.Errorf("error finding first different: %v", err)
	}

	if i != 0 {
		err = findFirstDifferentFirstRuleWRightAndLeft(equation, &e, &j, &rSym, &lSymE, mode)
		if err != nil {
			return false, Equation{}, fmt.Errorf("error finding first different: %v", err)
		}
		if j == 0 {
			return false, Equation{}, nil
		} else {
			if i > j {
				if i != e.RightPart.Length {
					return false, Equation{}, nil
				}

				var leftSymbolsFirst, leftSymbolsSecond []symbol.Symbol

				if mode == FORWARD {
					leftSymbolsFirst = make([]symbol.Symbol, len(e.LeftPart.Symbols))
					copy(leftSymbolsFirst, e.LeftPart.Symbols)
					leftSymbolsSecond = make([]symbol.Symbol, len(equation.LeftPart.Symbols[i:]))
					copy(leftSymbolsSecond, equation.LeftPart.Symbols[i:])
				} else {
					leftSymbolsFirst = make([]symbol.Symbol, len(equation.LeftPart.Symbols[:equation.LeftPart.Length-i]))
					copy(leftSymbolsFirst, equation.LeftPart.Symbols[:equation.LeftPart.Length-i])
					var leftSymbolsSecond = make([]symbol.Symbol, len(e.LeftPart.Symbols))
					copy(leftSymbolsSecond, e.LeftPart.Symbols)
				}

				var rightSymbols = make([]symbol.Symbol, len(equation.RightPart.Symbols))
				copy(rightSymbols, equation.RightPart.Symbols)
				newEq := NewEquation(append(leftSymbolsFirst, leftSymbolsSecond...), rightSymbols)
				newEq.Reduce()
				return true, newEq, nil
			} else {
				if j != e.LeftPart.Length {
					return false, Equation{}, nil
				}

				var rightSymbolsFirst, rightSymbolsSecond []symbol.Symbol

				if mode == FORWARD {
					rightSymbolsFirst = make([]symbol.Symbol, len(e.RightPart.Symbols))
					copy(rightSymbolsFirst, e.RightPart.Symbols)
					rightSymbolsSecond = make([]symbol.Symbol, len(equation.RightPart.Symbols[j:]))
					copy(rightSymbolsSecond, equation.RightPart.Symbols[j:])
				} else {
					rightSymbolsFirst = make([]symbol.Symbol, len(equation.RightPart.Symbols[:equation.RightPart.Length-j]))
					copy(rightSymbolsFirst, equation.RightPart.Symbols[:equation.RightPart.Length-j])
					rightSymbolsSecond = make([]symbol.Symbol, len(e.RightPart.Symbols))
					copy(rightSymbolsSecond, e.RightPart.Symbols)
				}

				var leftSymbols = make([]symbol.Symbol, len(equation.LeftPart.Symbols))
				copy(leftSymbols, equation.LeftPart.Symbols)
				newEq := NewEquation(leftSymbols, append(rightSymbolsFirst, rightSymbolsSecond...))
				newEq.Reduce()
				return true, newEq, nil
			}
		}
	}
	return false, Equation{}, nil
}

func findFirstDifferentFirstRuleWRightAndRight(equation, e *Equation, i *int, sym, symE *symbol.Symbol, mode int) error {
	return findFirstDifferentFirstRule(equation, e, i, sym, symE, mode, RIGHT, RIGHT)
}

func findFirstDifferentFirstRuleWRightAndLeft(equation, e *Equation, i *int, sym, symE *symbol.Symbol, mode int) error {
	return findFirstDifferentFirstRule(equation, e, i, sym, symE, mode, RIGHT, LEFT)
}

func findFirstDifferentFirstRuleWLeftAndRight(equation, e *Equation, i *int, sym, symE *symbol.Symbol, mode int) error {
	return findFirstDifferentFirstRule(equation, e, i, sym, symE, mode, LEFT, RIGHT)
}

func findFirstDifferentFirstRuleWLeftAndLeft(equation, e *Equation, i *int, sym, symE *symbol.Symbol, mode int) error {
	return findFirstDifferentFirstRule(equation, e, i, sym, symE, mode, LEFT, LEFT)
}

func findFirstDifferentFirstRule(equation, e *Equation, i *int, sym, symE *symbol.Symbol, mode int, side, sideE int) error {
	var err error
	var first, second *EqPart
	if side == RIGHT {
		first = &equation.RightPart
	} else {
		first = &equation.LeftPart
	}
	if sideE == RIGHT {
		second = &e.RightPart
	} else {
		second = &e.LeftPart
	}
	for *i < first.Length && *i < second.Length {
		*sym, err = first.GetSymbolMode(*i, mode)
		if err != nil {
			return fmt.Errorf("error getting symbol: %v", err)
		}
		*symE, err = second.GetSymbolMode(*i, mode)
		if err != nil {
			return fmt.Errorf("error getting symbol: %v", err)
		}
		if *sym == *symE {
			*i++
		} else {
			break
		}
	}

	return nil
}

func (equation *Equation) applyThirdRule(e Equation, mode int) (bool, Equation, error) {
	var err error
	var found bool
	var i, j, i1, j1 int
	var lSym, rSym, lSymE, rSymE symbol.Symbol
	found, err = findFirstDifferentThirdRule(equation, &i, &j, &lSym, &rSym, mode)
	if err != nil {
		return false, Equation{}, fmt.Errorf("error finding different vars: %v", err)
	}

	if !found {
		return false, Equation{}, nil
	}

	found, err = findFirstDifferentThirdRule(&e, &i1, &j1, &lSymE, &rSymE, mode)
	if err != nil {
		return false, Equation{}, fmt.Errorf("error finding different vars: %v", err)
	}

	if !found {
		return false, Equation{}, nil
	}

	var k int
	var firstIndex, thirdIndex int

	if checkThirdRuleEqualSides(lSym, rSym, lSymE, rSymE) {
		if mode == FORWARD {
			i -= i1
		} else {
			i += i1
		}
		for ; k+i < equation.LeftPart.Length && k < e.LeftPart.Length; k++ {
			lSym, err = equation.LeftPart.GetSymbolMode(k+i, mode)
			if err != nil {
				return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
			}
			lSymE, err = e.LeftPart.GetSymbolMode(k, mode)
			if err != nil {
				return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
			}
			if lSym != lSymE {
				break
			}
		}
		if k == e.LeftPart.Length {
			if mode == FORWARD {
				firstIndex = i
				thirdIndex = i + k
			} else {
				firstIndex = equation.LeftPart.Length - i - k
				thirdIndex = equation.LeftPart.Length - i
			}
			newEq := createNewEqThirdRuleWLeftAndRight(*equation, e, firstIndex, thirdIndex)
			newEq.Reduce()
			return true, newEq, nil
		} else {
			if mode == FORWARD {
				j -= j1
			} else {
				j += j1
			}
			k = 0
			for ; k+j < equation.RightPart.Length && k < e.RightPart.Length; k++ {
				rSym, err = equation.RightPart.GetSymbolMode(k+j, mode)
				if err != nil {
					return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
				}
				rSymE, err = e.RightPart.GetSymbolMode(k, mode)
				if err != nil {
					return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
				}
				if rSym != rSymE {
					break
				}
			}
			if k == e.RightPart.Length {
				if mode == FORWARD {
					firstIndex = j
					thirdIndex = j + k
				} else {
					firstIndex = equation.RightPart.Length - j - k
					thirdIndex = equation.RightPart.Length - j
				}
				newEq := createNewEqThirdRuleWRightAndLeft(*equation, e, firstIndex, thirdIndex)
				return true, newEq, nil
			}
		}
	} else if checkThirdRuleDifferentSides(lSym, rSym, lSymE, rSymE) {
		if mode == FORWARD {
			i -= i1
		} else {
			i += i1
		}
		for ; k+i < equation.RightPart.Length && k < e.RightPart.Length; k++ {
			lSym, err = equation.LeftPart.GetSymbolMode(k+i, mode)
			if err != nil {
				return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
			}
			rSymE, err = e.RightPart.GetSymbolMode(k, mode)
			if err != nil {
				return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
			}
			if lSym != rSymE {
				break
			}
		}
		if k == e.LeftPart.Length {
			if mode == FORWARD {
				firstIndex = i
				thirdIndex = i + k
			} else {
				firstIndex = equation.LeftPart.Length - i - k
				thirdIndex = equation.LeftPart.Length - i
			}
			newEq := createNewEqThirdRuleWLeftAndLeft(*equation, e, firstIndex, thirdIndex)
			return true, newEq, nil
		} else {
			if mode == FORWARD {
				i -= j1
			} else {
				i += j1
			}
			k = 0
			for ; k+j < equation.RightPart.Length && k < e.LeftPart.Length; k++ {
				rSym, err = equation.RightPart.GetSymbolMode(k+j, mode)
				if err != nil {
					return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
				}
				lSymE, err = e.LeftPart.GetSymbolMode(k, mode)
				if err != nil {
					return false, Equation{}, fmt.Errorf("error getting symbol: %v", err)
				}
				if rSym != lSymE {
					break
				}
			}
			if k == e.RightPart.Length {
				if mode == FORWARD {
					firstIndex = j
					thirdIndex = j + k
				} else {
					firstIndex = equation.RightPart.Length - j - k
					thirdIndex = equation.RightPart.Length - j
				}

				newEq := createNewEqThirdRuleWRightAndRight(*equation, e, firstIndex, thirdIndex)
				return true, newEq, nil
			}
		}
	}

	return false, Equation{}, nil
}

func findFirstDifferentThirdRule(e *Equation, i, j *int, lSym, rSym *symbol.Symbol, mode int) (bool, error) {
	var err error
	var foundLeft, foundRight bool
	for *i < e.LeftPart.Length && *j < e.RightPart.Length {
		*lSym, err = e.LeftPart.GetSymbolMode(*i, mode)
		if err != nil {
			return foundRight && foundLeft, fmt.Errorf("error getting symbol: %v", err)
		}
		*rSym, err = e.RightPart.GetSymbolMode(*j, mode)
		if err != nil {
			return foundRight && foundLeft, fmt.Errorf("error getting symbol: %v", err)
		}
		if foundLeft && foundRight {
			if *lSym == *rSym {
				foundRight = false
				foundLeft = false
				*i++
				*j++
				continue
			} else {
				break
			}
		}
		if !foundLeft {
			if symbol.IsVar(*lSym) {
				foundLeft = true
			} else {
				*i++
			}
		}

		if !foundRight {
			if symbol.IsVar(*rSym) {
				foundRight = true
			} else {
				*j++
			}
		}
	}

	return foundRight && foundLeft, nil
}

func createNewEqThirdRuleWLeftAndLeft(equation, e Equation, firstIndex, thirdIndex int) Equation {
	return createNewEqThirdRule(equation, e, firstIndex, thirdIndex, LEFT, LEFT)
}

func createNewEqThirdRuleWLeftAndRight(equation, e Equation, firstIndex, thirdIndex int) Equation {
	return createNewEqThirdRule(equation, e, firstIndex, thirdIndex, LEFT, RIGHT)
}

func createNewEqThirdRuleWRightAndLeft(equation, e Equation, firstIndex, thirdIndex int) Equation {
	return createNewEqThirdRule(equation, e, firstIndex, thirdIndex, RIGHT, LEFT)
}

func createNewEqThirdRuleWRightAndRight(equation, e Equation, firstIndex, thirdIndex int) Equation {
	return createNewEqThirdRule(equation, e, firstIndex, thirdIndex, RIGHT, RIGHT)
}

func createNewEqThirdRule(equation, e Equation, firstIndex, thirdIndex int, side, sideE int) Equation {
	var newEq Equation
	var equationToSplit, equationOriginal []symbol.Symbol
	if side == RIGHT {
		equationToSplit = equation.RightPart.Symbols
		equationOriginal = equation.LeftPart.Symbols
	} else {
		equationToSplit = equation.LeftPart.Symbols
		equationOriginal = equation.RightPart.Symbols
	}

	var symbolsFirst = make([]symbol.Symbol, len(equationToSplit[:firstIndex]))
	copy(symbolsFirst, equationToSplit[:firstIndex])
	var symbolsSecond []symbol.Symbol
	if sideE == RIGHT {
		symbolsSecond = make([]symbol.Symbol, len(e.RightPart.Symbols))
		copy(symbolsSecond, e.RightPart.Symbols)
	} else {
		symbolsSecond = make([]symbol.Symbol, len(e.LeftPart.Symbols))
		copy(symbolsSecond, e.LeftPart.Symbols)
	}

	var symbolsThird = make([]symbol.Symbol, len(equationToSplit[thirdIndex:]))
	copy(symbolsThird, equationToSplit[thirdIndex:])

	newS := append(symbolsFirst, append(symbolsSecond, symbolsThird...)...)

	var originalCopy = make([]symbol.Symbol, len(equationOriginal))
	copy(originalCopy, equationOriginal)
	if side == RIGHT {
		newEq = NewEquation(originalCopy, newS)
	} else {
		newEq = NewEquation(newS, originalCopy)
	}
	newEq.Reduce()
	return newEq
}

func checkThirdRuleEqualSides(equationLeft symbol.Symbol, equationRight symbol.Symbol, eLeft symbol.Symbol, eRight symbol.Symbol) bool {
	return equationLeft == eLeft && equationRight == eRight
}

func checkThirdRuleDifferentSides(equationLeft symbol.Symbol, equationRight symbol.Symbol, eLeft symbol.Symbol, eRight symbol.Symbol) bool {
	return equationLeft == eRight && equationRight == eLeft
}
