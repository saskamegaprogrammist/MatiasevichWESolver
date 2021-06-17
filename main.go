package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/logger"
	matlog "github.com/saskamegaprogrammist/MatiasevichWESolver/logger"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver"
	"os"
	"sort"
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

func parseFLags() (bool, string, string, int, bool, string, bool, bool, bool, bool, bool, bool, bool, bool) {
	fullGraph := flag.Bool("full_graph", false, "print full graph")
	fullSystem := flag.Bool("full_system", false, "solve full system")
	inputFile := flag.String("input_file", "", "input filename")
	inputDir := flag.String("input_directory", "", "input directory")
	cycleRange := flag.Int("cycle_range", 0, "cycle depth")
	makePng := flag.Bool("png", false, "create graph png")
	makeDot := flag.Bool("dot", false, "create dot description")
	outputDir := flag.String("output_directory", ".", "output directory")
	splitByEquidecomposability := flag.Bool("use_eq_split", false, "split equation into system")
	lengthAnalysis := flag.Bool("use_length_analysis", false, "use length analysis")
	simplification := flag.Bool("use_simplification", false, "use simplification")
	applying := flag.Bool("use_applying", false, "use applying")
	defaultName := flag.Bool("default_name", false, "use default filename")
	solveSystem := flag.Bool("solve_system", false, "solve equations system")

	flag.Parse()
	return *fullGraph, *inputFile, *inputDir, *cycleRange, *makePng,
		*outputDir, *makeDot, *splitByEquidecomposability, *fullSystem,
		*lengthAnalysis, *simplification, *defaultName, *solveSystem, *applying
}

func process(inputSource *os.File, fullGraph bool, makePng bool, makeDot bool,
	cycleRange int, outputDir string, splitByEq bool, fullSystem bool,
	lengthAnalysis bool, simplification bool, defaultName bool, solveSystem bool, applying bool) {
	var err error
	var equations = make([]string, 0)
	var equation string
	scanner := bufio.NewScanner(inputSource)
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
	var eqSolver solver.Solver
	solveOptions := solver.SolveOptions{
		LengthAnalysis:             lengthAnalysis,
		SplitByEquidecomposability: splitByEq,
		CycleRange:                 cycleRange,
		FullGraph:                  fullGraph,
		AlgorithmMode:              algorithmType,
		FullSystem:                 fullSystem,
		NeedsSimplification:        simplification,
		ApplyEquations:             applying,
	}
	printOptions := solver.PrintOptions{
		Dot:            makeDot,
		Png:            makePng,
		OutputDir:      outputDir,
		UseDefaultName: defaultName,
		InputFilename:  inputSource.Name(),
	}
	if solveSystem {
		for {
			err = scanInput(scanner)
			if err != nil {
				return
			}
			equation = scanner.Text()
			if equation == "" {
				break
			}
			equations = append(equations, equation)
		}
		err = eqSolver.InitWithSystem(constantsAlph, varsAlph, equations, printOptions, solveOptions)
	} else {
		err = scanInput(scanner)
		if err != nil {
			return
		}
		equation = scanner.Text()
		err = eqSolver.Init(constantsAlph, varsAlph, equation, printOptions, solveOptions)
	}

	if err != nil {
		logger.Errorf("error initializing solver: %v", err)
		return
	}
	hasSolution, measuredTime, err := eqSolver.Solve()
	if err != nil {
		logger.Errorf("error writing graph: %v", err)
	}
	fmt.Printf("took time: %v \ngot solution: %s \n\n", measuredTime, hasSolution)
}

func main() {
	matlog.LoggerSetup()
	fullGraph, inputFilename, inputDirName, cycleRange, makePng, outputDir, makeDot,
		splitByEquidecomposability, fullSystem, lengthAnalysis, simplification, defaultName, solveSystem, applying := parseFLags()

	if inputDirName != "" {
		inputDir, err := os.Open(inputDirName)
		if err != nil {
			logger.Errorf("error opening directory: %v", err)
		}
		files, err := inputDir.Readdir(-1)
		inputDir.Close()
		if err != nil {
			logger.Errorf("error closing directory: %v", err)
		}

		sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() }) //sorting files by name
		var fullPath string
		for _, file := range files {
			fullPath = inputDirName + string(os.PathSeparator) + file.Name()
			inputFile, err := os.Open(fullPath)
			if err != nil {
				logger.Errorf("error opening input file: %v", err)
			}
			process(inputFile, fullGraph, makePng, makeDot,
				cycleRange, outputDir, splitByEquidecomposability,
				fullSystem, lengthAnalysis, simplification, defaultName, solveSystem, applying)
		}
	} else if inputFilename != "" {
		inputFile, err := os.Open(inputFilename)
		if err != nil {
			logger.Errorf("error opening input file: %v", err)
		}
		process(inputFile, fullGraph, makePng, makeDot,
			cycleRange, outputDir, splitByEquidecomposability,
			fullSystem, lengthAnalysis, simplification, defaultName, solveSystem, applying)
	} else {
		process(os.Stdin, fullGraph, makePng, makeDot,
			cycleRange, outputDir, splitByEquidecomposability,
			fullSystem, lengthAnalysis, simplification, defaultName, solveSystem, applying)
	}
}
