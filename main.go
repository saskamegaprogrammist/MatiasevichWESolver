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

func parseFLags() (bool, string, string, int, bool, string, bool, bool, bool) {
	fullGraph := flag.Bool("full_graph", false, "print full graph")
	fullSystem := flag.Bool("full_system", false, "solve full system")
	inputFile := flag.String("input_file", "", "input filename")
	inputDir := flag.String("input_directory", "", "input directory")
	cycleRange := flag.Int("cycle_range", 0, "cycle depth")
	makePng := flag.Bool("png", false, "create graph png")
	outputDir := flag.String("output_directory", ".", "output directory")
	splitByEquidecomposability := flag.Bool("use_eq_split", false, "split equation into system")
	lengthAnalysis := flag.Bool("use_length_analysis", false, "use length analysis")
	flag.Parse()
	return *fullGraph, *inputFile, *inputDir, *cycleRange, *makePng, *outputDir, *splitByEquidecomposability, *fullSystem, *lengthAnalysis
}

func process(inputSource *os.File, fullGraph bool, makePng bool, cycleRange int, outputDir string, splitByEq bool, fullSystem bool, lengthAnalysis bool) {
	var err error
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
	err = scanInput(scanner)
	if err != nil {
		return
	}
	equation := scanner.Text()

	var eqSolver solver.Solver
	solveOptions := solver.SolveOptions{
		LengthAnalysis:             lengthAnalysis,
		SplitByEquidecomposability: splitByEq,
		CycleRange:                 cycleRange,
		FullGraph:                  fullGraph,
		AlgorithmMode:              algorithmType,
		FullSystem:                 fullSystem,
	}
	printOptions := solver.PrintOptions{
		Png:       makePng,
		OutputDir: outputDir,
	}
	err = eqSolver.Init(constantsAlph, varsAlph, equation, printOptions, solveOptions)
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
	fullGraph, inputFilename, inputDirName, cycleRange, makePng, outputDir, splitByEquidecomposability, fullSystem, lengthAnalysis := parseFLags()

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

		for _, file := range files {
			inputFile, err := os.Open(fmt.Sprintf("%s%c%s", inputDirName, os.PathSeparator, file.Name()))
			if err != nil {
				logger.Errorf("error opening input file: %v", err)
			}
			process(inputFile, fullGraph, makePng, cycleRange, outputDir, splitByEquidecomposability, fullSystem, lengthAnalysis)
		}
	} else if inputFilename != "" {
		inputFile, err := os.Open(inputFilename)
		if err != nil {
			logger.Errorf("error opening input file: %v", err)
		}
		process(inputFile, fullGraph, makePng, cycleRange, outputDir, splitByEquidecomposability, fullSystem, lengthAnalysis)
	} else {
		process(os.Stdin, fullGraph, makePng, cycleRange, outputDir, splitByEquidecomposability, fullSystem, lengthAnalysis)
	}
}
