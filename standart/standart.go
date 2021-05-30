package standart

import "github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"

func Min(first int, second int) int {
	if first < second {
		return first
	} else {
		return second
	}
}

func Max(first int, second int) int {
	if first > second {
		return first
	} else {
		return second
	}
}

func CopySymbolIntMap(originalMap *map[symbol.Symbol]int, destMap *map[symbol.Symbol]int) {
	for k, v := range *originalMap {
		(*destMap)[k] = v
	}
}

func CopySymbolBoolMap(originalMap *map[symbol.Symbol]bool, destMap *map[symbol.Symbol]bool) {
	for k, v := range *originalMap {
		(*destMap)[k] = v
	}
}

func CopyIntBoolMap(originalMap *map[int]bool, destMap *map[int]bool) {
	for k, v := range *originalMap {
		(*destMap)[k] = v
	}
}

func SymbolArrayFromBoolMap(symbolMap map[symbol.Symbol]bool) []symbol.Symbol {
	var symbolArray []symbol.Symbol
	for k, _ := range symbolMap {
		symbolArray = append(symbolArray, k)
	}
	return symbolArray
}

func SymbolArrayFromIntMap(symbolMap map[symbol.Symbol]int) []symbol.Symbol {
	var symbolArray []symbol.Symbol
	for k, _ := range symbolMap {
		symbolArray = append(symbolArray, k)
	}
	return symbolArray
}

func MergeMaps(symbolMap *map[symbol.Symbol]bool, graphMap map[symbol.Symbol]int) {
	for k, _ := range graphMap {
		(*symbolMap)[k] = true
	}
}
