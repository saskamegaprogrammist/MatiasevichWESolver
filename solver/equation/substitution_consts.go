package equation

const STANDARD = 1
const SPLITTING = 2
const REDUCING = 3
const APPLYING = 4
const SIMPLIFYING = 5

var sTypesMap = map[int]string{SPLITTING: "splitting", REDUCING: "reducing", APPLYING: "applying", SIMPLIFYING: "simplifying"}
