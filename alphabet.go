package shortuuid

import (
	"fmt"
	"math"
	"slices"
	"unicode/utf8"
)

// DefaultAlphabet is the default alphabet used for base57 encoding.
// It excludes similar-looking characters (0, 1, I, O, l) to avoid confusion.
const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// alphabet represents a character set for base-N encoding. It stores the
// sorted, deduplicated characters along with precomputed values for efficient
// encoding and decoding.
type alphabet struct {
	chars    []rune  // sorted, deduplicated characters
	len      int64   // number of characters in the alphabet
	encLen   uint8   // maximum encoded length for a 128-bit value
	maxBytes uint8   // maximum UTF-8 bytes needed for any character
}

// newAlphabet creates a new alphabet from the given string. Removes
// duplicates and sorts the characters to ensure reproducibility.
//
// Panics if the alphabet (after removing duplicates) has fewer than 2
// characters. An alphabet must have at least 2 characters to be usable for
// base-N encoding.
func newAlphabet(s string) alphabet {
	abc := []rune(s)
	slices.Sort(abc)
	abc = slices.Compact(abc)

	if len(abc) < 2 {
		panic("encoding alphabet must be at least two characters")
	}

	return alphabet{
		chars:    abc,
		len:      int64(len(abc)),
		encLen:   uint8(math.Ceil(128 / math.Log2(float64(len(abc))))),
		maxBytes: uint8(utf8.RuneLen(abc[len(abc)-1])),
	}
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
