package equation

type EquationsByLength []Equation

func (ebl EquationsByLength) Len() int { return len(ebl) }
func (ebl EquationsByLength) Less(i, j int) bool {
	return ebl[i].structure.Size() < ebl[j].structure.Size()
}
func (ebl EquationsByLength) Swap(i, j int) { ebl[i], ebl[j] = ebl[j], ebl[i] }
