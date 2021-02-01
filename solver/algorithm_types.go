package solver

import "fmt"

const (
	FINITE   = 1
	INFINITE = 2
)

var typesMap = map[string]int64{
	"Finite":   FINITE,
	"Standart": INFINITE,
}

func matchAlgorithmType(algorithmType string) (int64, error) {
	intType := typesMap[algorithmType]
	if intType == 0 {
		return intType, fmt.Errorf("invalid algorithm type: %s", algorithmType)
	}
	return intType, nil
}
