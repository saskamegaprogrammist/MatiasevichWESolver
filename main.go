package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/logger"
	matlog "github.com/saskamegaprogrammist/MatiasevichWESolver/logger"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver"
	"os"
)

func handleScannerError(scanner *bufio.Scanner) error {
	var scannerErr error
	if err := scanner.Err(); err != nil {
		scannerErr = fmt.Errorf("scanner error: %v", err)
		logger.Errorf(scannerErr.Error())
	}
	return scannerErr
}

func scanInput(scanner *bufio.Scanner) error {
	scanner.Scan()
	return handleScannerError(scanner)
}

func parseFLags() bool {
	fullGraph := flag.Bool("full_graph", false, "print full graph")
	flag.Parse()
	return *fullGraph
}

func main() {
	var err error
	matlog.LoggerSetup()
	fullGraph := parseFLags()
	scanner := bufio.NewScanner(os.Stdin)
	err = handleScannerError(scanner)
	if err != nil {
		return
	}
	err = scanInput(scanner)
	if err != nil {
		return
	}
	algorithmType := scanner.Text()
	err = scanInput(scanner)
	if err != nil {
		return
	}
	constantsAlph := scanner.Text()
	err = scanInput(scanner)
	if err != nil {
		return
	}
	varsAlph := scanner.Text()
	err = scanInput(scanner)
	if err != nil {
		return
	}
	equation := scanner.Text()

	var solver solver.Solver
	err = solver.Init(algorithmType, constantsAlph, varsAlph, equation, fullGraph)
	if err != nil {
		logger.Errorf("error initializing solver: %v", err)
		return
	}
	hasSolution, err := solver.Solve()
	if err != nil {
		logger.Errorf("error writing graph: %v", err)
	}
	fmt.Print(hasSolution)
}
