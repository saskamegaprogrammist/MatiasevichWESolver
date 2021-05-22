package equation

import "fmt"

const (
	OPENBR  = "{"
	CLOSEBR = "}"
	COMMA   = ","
	SPACE   = " "
)

type Alphabet struct {
	words         []string
	size          int
	maxWordLength int
}

func NewAlphabet(words []string, size int, maxWordLength int) (Alphabet, error) {
	if len(words) != size {
		return Alphabet{}, fmt.Errorf("wrong size: %d != %d", len(words), size)
	}
	return Alphabet{
		words:         words,
		size:          size,
		maxWordLength: maxWordLength,
	}, nil
}

func (alphabet *Alphabet) SetMaxWordLength(len int) {
	alphabet.maxWordLength = len
}

func (alphabet *Alphabet) AddWord(word string) {
	alphabet.words = append(alphabet.words, word)
	alphabet.size++
}

func (alphabet *Alphabet) Has(word string) bool {
	for _, w := range alphabet.words {
		if w == word {
			return true
		}
	}
	return false
}

func (alphabet *Alphabet) At(index int) (string, error) {
	if index >= alphabet.size || index < 0 {
		return "", fmt.Errorf("invalid index: %d", index)
	}
	return alphabet.words[index], nil
}
