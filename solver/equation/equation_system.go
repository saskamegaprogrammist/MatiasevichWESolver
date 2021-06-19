package equation

import (
	"fmt"
)

type EqSystem struct {
	Equations []Equation
	Size      int
}

func (eqSystem *EqSystem) PrintInfo() {
	for i, eq := range eqSystem.Equations {
		fmt.Printf("equation %d:\n", i)
		eq.Print()
	}
}

func (eqSystem *EqSystem) Print() {
	fmt.Print(eqSystem.String())
}

func (eqSystem *EqSystem) String() string {
	var result string
	for _, eq := range eqSystem.Equations {
		result += fmt.Sprintf("%v\n", eq.String())
	}
	return result
}

func SystemFromEquation(eq Equation) EqSystem {
	return EqSystem{
		Equations: []Equation{eq},
		Size:      1,
	}
}

func SystemFromEquations(eqs []Equation) EqSystem {
	return EqSystem{
		Equations: eqs,
		Size:      len(eqs),
	}
}

func (eqSystem *EqSystem) CheckSameness(system *EqSystem) bool {
	if eqSystem.Size != system.Size {
		return false
	}
	var copySystemEqs = make([]Equation, len(system.Equations))
	var copySystemEqsHelp []Equation
	copy(copySystemEqs, system.Equations)
	for _, eq := range eqSystem.Equations {
		for j, sysEq := range copySystemEqs {
			if !eq.CheckSameness(&sysEq) {
				copySystemEqsHelp = append(copySystemEqsHelp, sysEq)
			} else {
				copySystemEqsHelp = append(copySystemEqsHelp, copySystemEqs[j+1:]...)
				break
			}
		}
		if len(copySystemEqs) == len(copySystemEqsHelp) {
			return false
		} else {
			copySystemEqs = copySystemEqsHelp
			copySystemEqsHelp = nil // clear array
		}
	}
	return true
}

func (eqSystem *EqSystem) CheckInequality() bool {
	for _, eq := range eqSystem.Equations {
		unequal, _ := eq.CheckInequality()
		if unequal {
			return true
		}
	}
	return false
}

func (eqSystem *EqSystem) CheckEquality() bool {
	for _, eq := range eqSystem.Equations {
		if !eq.CheckEquality() {
			return false
		}
	}
	return true
}
