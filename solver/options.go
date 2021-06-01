package solver

type SolveOptions struct {
	LengthAnalysis             bool
	SplitByEquidecomposability bool
	CycleRange                 int
	FullGraph                  bool
	FullSystem                 bool
	AlgorithmMode              string
	SaveLettersSubstitutions   bool
}

type PrintOptions struct {
	Png       bool
	Dot       bool
	OutputDir string
}
