package equation

const CONJUNCTION int = 1
const DISJUNCTION int = 2
const SINGLE_EQUATION int = 3

var typesMap = map[int]string{CONJUNCTION: "CONJUNCTION", DISJUNCTION: "DISJUNCTION", SINGLE_EQUATION: "EQUATION"}

const EQ_TYPE_SIMPLE int = 1
const EQ_TYPE_W_EMPTY int = 2
