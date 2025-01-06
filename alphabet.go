package shortuuid

import (
	"fmt"
	"math"
	"slices"
	"unicode/utf8"
)

// DefaultAlphabet is the default alphabet used.
const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type alphabet struct {
	chars    []rune
	len      int64
	encLen   uint8
	maxBytes uint8
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
		chars:    abc,
		len:      int64(len(abc)),
		encLen:   uint8(math.Ceil(128 / math.Log2(float64(len(abc))))),
		maxBytes: uint8(utf8.RuneLen(abc[len(abc)-1])),
	}

	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}

// Index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (a *alphabet) Index(t rune) (int64, error) {
	i, j := 0, int(a.len)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a.chars[h] < t {
			i = h + 1
		} else {
			j = h
		}
	}
	if i >= int(a.len) || a.chars[i] != t {
		return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
	}
	return int64(i), nil
}
