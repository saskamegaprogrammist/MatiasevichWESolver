# Compilers course work

Matiasevich word equations solver

### flags:

- full_graph - 
*boolean* create full graph description

- input_file - 
*string* full path to input file with equation description

- input_directory - 
*string* full path to input directory with equation description files

- output_directory - 
*string* full path to output directory with graph description files

- png - 
*boolean* create graph png image

- cycle_range - 
*int* cycle depth

### run app:

` go run main.go -full_graph -input_directory=checked `

### run tests:

` cd solver `

`go test `


### Input format:

- Finite | Standard - *algorithm type*
- {} | {const(, const)*}  - *constants alphabet*
- {} | {var(, var)*} - *variables alphabet*
- u a v = v a u - *equation*

### Output format:

- l g l = A A Y - *equation*
- Standard - *algorithm type*
- took time: 307.88µs - *time took algorithm to run excluding png creation*
- got solution: TRUE - *answer, whether algorithm has solutions or not*