package solver

import "fmt"

const (
	FINITE   = 1
	STANDARD = 2
)

var typesMap = map[string]int64{
	"Finite":   FINITE,
	"Standard": STANDARD,
}

func matchAlgorithmType(algorithmType string) (int64, error) {
	intType := typesMap[algorithmType]
	if intType == 0 {
		return intType, fmt.Errorf("invalid algorithm type: %s", algorithmType)
	}
	return intType, nil
}
