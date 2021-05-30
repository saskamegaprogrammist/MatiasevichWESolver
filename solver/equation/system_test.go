package equation

import (
	"fmt"
	"testing"
)

func TestEquationsSystem_Equals_Equations(t *testing.T) {
	var err error
	var eq1, eq2, eq3 Equation
	err = eq1.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}
	err = eq2.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}

	err = eq3.Init("y = b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}

	var eqSys1, eqSys2, eqSys3 EquationsSystem

	eqSys1 = NewSingleEquation(eq1)
	eqSys2 = NewSingleEquation(eq2)
	eqSys3 = NewSingleEquation(eq3)

	var same bool
	same = eqSys1.Equals(eqSys2)
	if !same {
		t.Errorf("TestEquationsSystem_Equals failed: eqSys1 and eqSys2 should be the same")
		return
	}

	same = eqSys1.Equals(eqSys3)
	if same {
		t.Errorf("TestEquationsSystem_Equals failed: eqSys1 and eqSys3 should not be the same")
		return
	}
}

func TestEquationsSystem_Equals_Disjunction(t *testing.T) {
	var err error
	var eq1, eq2, eq3, eq4 Equation
	err = eq1.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}
	err = eq2.Init("x a b = b y", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}

	err = eq3.Init("y = b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}

	err = eq4.Init("x = b", &constAlphNew, &varsAlphNew)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("TestEquationsSystem_Equals failed: error should be nil")
		return
	}

	var eqSys1, eqSys2, eqSys3, eqSys4 EquationsSystem

	eqSys1 = NewDisjunctionSystemFromEquations([]Equation{eq1, eq3})
	eqSys2 = NewDisjunctionSystemFromEquations([]Equation{eq2, eq4})

	eqSys3 = NewDisjunctionSystemFromEquations([]Equation{eq1})
	eqSys4 = NewDisjunctionSystemFromEquations([]Equation{eq2})

	var same bool
	same = eqSys1.Equals(eqSys2)
	if same {
		t.Errorf("TestEquationsSystem_Equals failed: eqSys1 and eqSys2 should not be the same")
		return
	}

	same = eqSys3.Equals(eqSys4)
	if !same {
		t.Errorf("TestEquationsSystem_Equals failed: eqSys3 and eqSys4 should be the same")
		return
	}
}
