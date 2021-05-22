package solver

type SolveOptions struct {
	LengthAnalysis             bool
	SplitByEquidecomposability bool
	CycleRange                 int
	FullGraph                  bool
	FullSystem                 bool
	AlgorithmMode              string
}

type PrintOptions struct {
	Png       bool
	OutputDir string
}
