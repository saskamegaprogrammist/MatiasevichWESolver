package solver

import (
	"fmt"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"time"
)

func (solver *Solver) solveEquationAsSystem(eq equation.Equation) (time.Duration, error) {
	timeStart := time.Now()
	err := solver.setWriter(eq) // setting filename with equation
	if err != nil {
		return 0, fmt.Errorf("error setting writer: %v", err)
	}
	tree := Node{
		number: "0",
		value:  equation.NewSingleEquation(eq),
	}
	err = solver.dotWriter.StartDOTDescription()
	if err != nil {
		return 0, fmt.Errorf("error writing DOT description: %v", err)
	}
	err = solver.solveWithSystem(&tree)
	if err != nil {
		return 0, fmt.Errorf("error solving equation: %v", err)
	}
	measuredTime := time.Since(timeStart)
	err = solver.dotWriter.EndDOTDescription()
	if err != nil {
		return measuredTime, fmt.Errorf("error writing DOT description: %v", err)
	}
	return measuredTime, nil
}

func (solver *Solver) solveWithSystem(node *Node) error {
	var err error
	err = solver.dotWriter.WriteNode(node)
	if err != nil {
		return fmt.Errorf("error writing node: %v", err)
	}
	if !solver.solveOptions.FullGraph && solver.hasSolution {
		return nil
	}
	if len(node.number) > solver.solveOptions.CycleRange {
		solver.cycled = true
		return nil
	}
	hasBeen, tr := checkHasBeen(node)
	if hasBeen {
		err = solver.dotWriter.WriteDottedEdge(node, tr)
		if err != nil {
			return fmt.Errorf("error writing dotted edge: %v", err)
		}
		//fmt.Println("HAS BEEN")
		//fmt.Println(node.number)
		return nil
	}
	var newEquations []equation.Equation
	for _, eq := range node.value.GetEquations() {
		if !checkFirstSymbols(&eq) || checkThirdRuleRight(&eq) && !checkThirdRuleLeft(&eq) ||
			!checkThirdRuleRight(&eq) && checkThirdRuleLeft(&eq) {
			solver.createFalseNode(node, REGULAR_FALSE)
			err = solver.dotWriter.WriteInfoNode(node.infoChild)
			if err != nil {
				return fmt.Errorf("error writing info node: %v", err)
			}
			err = solver.dotWriter.WriteInfoEdgeWithLabel(node, (node.infoChild).(*FalseNode))
			if err != nil {
				return fmt.Errorf("error writing info edge: %v", err)
			}
			return nil
		}
		system := eq.SplitByEquidecomposability()
		if system.Size != 1 {
			for _, neq := range system.Equations {
				if !neq.CheckEquality() {
					newEquations = append(newEquations, neq)
				}
			}
		} else {
			if !eq.CheckEquality() {
				newEquations = append(newEquations, eq)
			}
		}
	}
	if len(newEquations) == 0 {
		solver.createTrueNode(node)
		err = solver.dotWriter.WriteInfoNode(node.infoChild)
		if err != nil {
			return fmt.Errorf("error writing info node: %v", err)
		}
		err = solver.dotWriter.WriteInfoEdge(node, node.infoChild)
		if err != nil {
			return fmt.Errorf("error writing info edge: %v", err)
		}
		return nil
	} else if len(newEquations) > node.value.Size() {
		child := NewNodeWEquationsSystem(equation.Substitution{},
			"x"+node.number, node, equation.NewConjunctionSystemFromEquations(newEquations))
		err = solver.dotWriter.WriteNode(&child)
		if err != nil {
			return fmt.Errorf("error writing node: %v", err)
		}
		err = solver.dotWriter.WriteLabelEdgeBold(node, &child)
		if err != nil {
			return fmt.Errorf("error writing splitting edge: %v", err)
		}
		node = &child
	}
	firstEq := newEquations[0]
	if solver.algorithmType == FINITE {
		if checkSecondRuleLeftFinite(&firstEq) {
			substitution := equation.NewSubstitution(firstEq.RightPart.Symbols[0], []symbol.Symbol{firstEq.LeftPart.Symbols[0]})

			var substNewEquationsFirst []equation.Equation
			var newEq equation.Equation
			for _, neq := range newEquations {
				newEq = neq.Substitute(substitution)
				substNewEquationsFirst = append(substNewEquationsFirst, newEq)
			}

			child := NewNodeWEquationsSystem(equation.Substitution{},
				"a"+node.number, node, equation.NewConjunctionSystemFromEquations(substNewEquationsFirst))
			node.children = []*Node{&child}
			err = solver.dotWriter.WriteLabelEdge(node, &child, &firstEq.RightPart.Symbols[0], substitution.RightPart())
			if err != nil {
				return fmt.Errorf("error writing label edge: %v", err)
			}
		}

		if checkSecondRuleRightFinite(&firstEq) {
			substitution := equation.NewSubstitution(firstEq.LeftPart.Symbols[0], []symbol.Symbol{firstEq.RightPart.Symbols[0]})

			var substNewEquationsFirst []equation.Equation
			var newEq equation.Equation
			for _, neq := range newEquations {
				newEq = neq.Substitute(substitution)
				substNewEquationsFirst = append(substNewEquationsFirst, newEq)
			}

			child := NewNodeWEquationsSystem(equation.Substitution{},
				"b"+node.number, node, equation.NewConjunctionSystemFromEquations(substNewEquationsFirst))

			node.children = []*Node{&child}
			err = solver.dotWriter.WriteLabelEdge(node, &child, &firstEq.LeftPart.Symbols[0], substitution.RightPart())
			if err != nil {
				return fmt.Errorf("error writing label edge: %v", err)
			}
		}
	}
	if checkSecondRuleLeft(&firstEq) {
		var newEq equation.Equation
		substFirst := equation.NewSubstitution(firstEq.RightPart.Symbols[0], []symbol.Symbol{symbol.Empty()})

		var substNewEquationsFirst []equation.Equation
		for _, neq := range newEquations {
			newEq = neq.Substitute(substFirst)
			substNewEquationsFirst = append(substNewEquationsFirst, newEq)
		}

		firstChild := NewNodeWEquationsSystem(equation.Substitution{},
			node.number+"1", node, equation.NewConjunctionSystemFromEquations(substNewEquationsFirst))

		substSecond := equation.NewSubstitution(firstEq.RightPart.Symbols[0], []symbol.Symbol{firstEq.LeftPart.Symbols[0], firstEq.RightPart.Symbols[0]})

		var substNewEquationsSecond []equation.Equation
		for _, neq := range newEquations {
			newEq = neq.Substitute(substSecond)
			substNewEquationsSecond = append(substNewEquationsSecond, newEq)
		}

		secondChild := NewNodeWEquationsSystem(equation.Substitution{},
			node.number+"2", node, equation.NewConjunctionSystemFromEquations(substNewEquationsSecond))

		node.children = []*Node{&firstChild, &secondChild}
		err = solver.dotWriter.WriteLabelEdge(node, &firstChild, &firstEq.RightPart.Symbols[0], substFirst.RightPart())
		if err != nil {
			return fmt.Errorf("error writing label edge: %v", err)
		}
		err = solver.dotWriter.WriteLabelEdge(node, &secondChild, &firstEq.RightPart.Symbols[0], substSecond.RightPart())
		if err != nil {
			return fmt.Errorf("error writing label edge: %v", err)
		}
	}
	if checkSecondRuleRight(&firstEq) {
		var newEq equation.Equation
		substFirst := equation.NewSubstitution(firstEq.LeftPart.Symbols[0], []symbol.Symbol{symbol.Empty()})

		var substNewEquationsFirst []equation.Equation
		for _, neq := range newEquations {
			newEq = neq.Substitute(substFirst)
			substNewEquationsFirst = append(substNewEquationsFirst, newEq)
		}

		firstChild := NewNodeWEquationsSystem(equation.Substitution{},
			node.number+"3", node, equation.NewConjunctionSystemFromEquations(substNewEquationsFirst))

		substSecond := equation.NewSubstitution(firstEq.LeftPart.Symbols[0], []symbol.Symbol{firstEq.RightPart.Symbols[0], firstEq.LeftPart.Symbols[0]})

		var substNewEquationsSecond []equation.Equation
		for _, neq := range newEquations {
			newEq = neq.Substitute(substSecond)
			substNewEquationsSecond = append(substNewEquationsSecond, newEq)
		}

		secondChild := NewNodeWEquationsSystem(equation.Substitution{},
			node.number+"4", node, equation.NewConjunctionSystemFromEquations(substNewEquationsSecond))

		node.children = []*Node{&firstChild, &secondChild}

		err = solver.dotWriter.WriteLabelEdge(node, &firstChild, &firstEq.LeftPart.Symbols[0], substFirst.RightPart())
		if err != nil {
			return fmt.Errorf("error writing label edge: %v", err)
		}
		err = solver.dotWriter.WriteLabelEdge(node, &secondChild, &firstEq.LeftPart.Symbols[0], substSecond.RightPart())
		if err != nil {
			return fmt.Errorf("error writing label edge: %v", err)
		}
	}
	for _, child := range node.children {
		err = solver.solveWithSystem(child)
		if err != nil {
			return fmt.Errorf("error solving for child: %v", err)
		}
	}
	return nil
}
