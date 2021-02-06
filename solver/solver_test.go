package solver

import (
	"fmt"
	"testing"
)

var test1InitErrorMessage = "error matching alphabet type: invalid algorithm type: Invalid"

func Test_Init_Error_1(t *testing.T) {
	var solver Solver
	err := solver.Init("Invalid", "", "", "")
	if err == nil {
		t.Errorf("Test_Init_Error_1 failed: error shouldn\\'t be nil")
	} else {
		if err.Error() != test1InitErrorMessage {
			t.Errorf("Test_Init_Error_1 failed: wrong error message")
		}
	}
}

var test2InitErrorMessage = "error parsing constants: invalid constants alphabet: a,c"

func Test_Init_Error_2(t *testing.T) {
	var solver Solver
	err := solver.Init("Finite", "a,c", "", "")
	if err == nil {
		t.Errorf("Test_Init_Error_2: error shouldn\\'t be nil")
	} else {
		if err.Error() != test2InitErrorMessage {
			t.Errorf("Test_Init_Error_2 failed: wrong error message")
		}
	}
}

var test3InitErrorMessage = "error parsing constants: empty constant in alphabet: {a,c,,s}"

func Test_Init_Error_3(t *testing.T) {
	var solver Solver
	err := solver.Init("Finite", "{a,c,,s}", "", "")
	if err == nil {
		t.Errorf("Test_Init_Error_3: error shouldn\\'t be nil")
	} else {
		if err.Error() != test3InitErrorMessage {
			t.Errorf("Test_Init_Error_3 failed: wrong error message")
		}
	}
}

var test4InitErrorMessage = "error parsing vars: invalid constants alphabet: b"

func Test_Init_Error_4(t *testing.T) {
	var solver Solver
	err := solver.Init("Finite", "{a,c,s}", "b", "")
	if err == nil {
		t.Errorf("Test_Init_Error_4: error shouldn\\'t be nil")
	} else {
		if err.Error() != test4InitErrorMessage {
			t.Errorf("Test_Init_Error_4 failed: wrong error message")
		}
	}
}

var test5InitErrorMessage = "error parsing vars: empty constant in alphabet: {a,,s}"

func Test_Init_Error_5(t *testing.T) {
	var solver Solver
	err := solver.Init("Finite", "{b,n}", "{a,,s}", "")
	if err == nil {
		t.Errorf("Test_Init_Error_5: error shouldn\\'t be nil")
	} else {
		if err.Error() != test5InitErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_Init_Error_5 failed: wrong error message")
		}
	}
}

var test6InitErrorMessage = "error parsing equation: invalid equation: ab"

func Test_Init_Error_6(t *testing.T) {
	var solver Solver
	err := solver.Init("Standart", "{b,n}", "{a,s}", "ab")
	if err == nil {
		t.Errorf("Test_Init_Error_6 error shouldn\\'t be nil")
	} else {
		if err.Error() != test6InitErrorMessage {
			fmt.Println(err.Error())
			t.Errorf("Test_Init_Error_6 failed: wrong error message")
		}
	}
}

var trueStr = "TRUE"
var falseStr = "FALSE"
var cycledStr = "CYCLED"

func Test_Solve_1(t *testing.T) {
	var solver Solver
	err := solver.Init("Standart", "{a}", "{u,v}", "uav=vau")
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_1 error should be nil")
	} else {
		result, _ := solver.Solve()
		if result != trueStr {
			t.Errorf("Test_Solve_1 result should be: %s, but got: %s", trueStr, result)
		}
	}
}

func Test_Solve_2(t *testing.T) {
	var solver Solver
	err := solver.Init("Standart", "{a,b}", "{u}", "uua=buu")
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_2 error should be nil")
	} else {
		result, _ := solver.Solve()
		if result != cycledStr {
			t.Errorf("Test_Solve_2 result should be: %s, but got: %s", cycledStr, result)
		}
	}
}

func Test_Solve_3(t *testing.T) {
	var solver Solver
	err := solver.Init("Standart", "{}", "{u,v,z}", "uuvv=zz")
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_3 error should be nil")
	} else {
		result, _ := solver.Solve()
		if result != trueStr {
			t.Errorf("Test_Solve_3 result should be: %s, but got: %s", trueStr, result)
		}
	}
}

func Test_Solve_4(t *testing.T) {
	var solver Solver
	err := solver.Init("Standart", "{a}", "{u}", "au=u")
	if err != nil {
		fmt.Printf("error initializing solver: %v \n", err)
		t.Errorf("Test_Solve_4 error should be nil")
	} else {
		result, _ := solver.Solve()
		if result != falseStr {
			t.Errorf("Test_Solve_4 result should be: %s, but got: %s", falseStr, result)
		}
	}
}
