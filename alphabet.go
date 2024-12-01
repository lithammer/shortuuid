package shortuuid

import (
	"fmt"
	"math"
	"slices"
)

// DefaultAlphabet is the default alphabet used.
const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	rune1Max        = 1<<7 - 1
)

type alphabet struct {
	chars       []rune
	indexMap    map[rune]int
	len         int64
	encLen      int64
	singleBytes bool
}

// Remove duplicates and sort it to ensure reproducibility.
func newAlphabet(s string) alphabet {
	abc := []rune(s)
	slices.Sort(abc)
	abc = slices.Compact(abc)

	if len(abc) < 2 {
		panic("encoding alphabet must be at least two characters")
	}

	a := alphabet{
		chars:       abc,
		len:         int64(len(abc)),
		encLen:      int64(math.Ceil(128 / math.Log2(float64(len(abc))))),
		singleBytes: true,
		indexMap:    make(map[rune]int, len(abc)),
	}
	for i, c := range a.chars {
		if c > rune1Max {
			a.singleBytes = false
		}
		a.indexMap[c] = i
	}

	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}

// Index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (a *alphabet) Index(t rune) (int64, error) {
	if i, ok := a.indexMap[t]; ok {
		return int64(i), nil
	}
	return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
}
