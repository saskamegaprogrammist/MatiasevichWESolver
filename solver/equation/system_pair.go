package equation

type EquationsSystemPair struct {
	value      EquationsSystem
	simplified EquationsSystem
}

func (esp *EquationsSystemPair) Value() EquationsSystem {
	return esp.value
}

func (esp *EquationsSystemPair) Simplified() EquationsSystem {
	return esp.simplified
}

func NewEquationsSystemPair(value EquationsSystem, simplified EquationsSystem) EquationsSystemPair {
	return EquationsSystemPair{
		value:      value,
		simplified: simplified,
	}
}
