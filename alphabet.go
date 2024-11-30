package shortuuid

import (
	"fmt"
	"math"
)

// DefaultAlphabet is the default alphabet used.
const DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type alphabet struct {
	chars  []rune
	len    int64
	encLen int64
}

// Remove duplicates and sort it to ensure reproducibility.
func newAlphabet(s string) alphabet {
	abc := []rune(s)
	// sortRunes can be replaced with slices.Sort if upgraded to go 1.18+
	// (use of generics avoids using reflect, and reduces allocations)
	sortRunes(abc)
	// dedupe can be replaced with slices.Compact() if upgraded to go 1.21+
	abc = dedupe(abc)

	if len(abc) < 2 {
		panic("encoding alphabet must be at least two characters")
	}

	a := alphabet{
		chars:  abc,
		len:    int64(len(abc)),
		encLen: int64(math.Ceil(128 / math.Log2(float64(len(abc))))),
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

// dedupe replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// dedupe modifies the contents of the slice s and returns the modified slice,
// which may have a smaller length.
// dedupe zeroes the elements between the new length and the original length.
func dedupe(s []rune) []rune {
	if len(s) < 2 {
		return s
	}
	for k := 1; k < len(s); k++ {
		if s[k] == s[k-1] {
			s2 := s[k:]
			for k2 := 1; k2 < len(s2); k2++ {
				if s2[k2] != s2[k2-1] {
					s[k] = s2[k2]
					k++
				}
			}

			for k2 := k; k2 < len(s); k2++ { // zero/nil out the obsolete elements, for GC
				s[k2] = 0
			}
			return s[:k]
		}
	}
	return s
}
