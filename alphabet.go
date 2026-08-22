package shortuuid

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"unicode/utf8"
)

// errNotInAlphabet is wrapped with the offending character by the decoders.
var errNotInAlphabet = errors.New("character is not part of the alphabet")

// notInAlphabet names the character that could not be decoded. Both decoders
// report it, so the format lives here rather than at each of their call sites.
func notInAlphabet(c rune) error {
	return fmt.Errorf("%w: %q", errNotInAlphabet, c)
}

// DefaultAlphabet is the default alphabet used for base57 encoding.
// It excludes similar-looking characters (0, 1, I, O, l) to avoid confusion.
const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// alphabet represents a character set for base-N encoding. It stores the
// sorted, deduplicated characters along with precomputed values for efficient
// encoding and decoding.
type alphabet struct {
	chars    []rune // sorted, deduplicated characters
	len      int64  // number of characters in the alphabet
	encLen   uint8  // maximum encoded length for a 128-bit value
	maxBytes uint8  // maximum UTF-8 bytes needed for any character

	// A 128-bit value is encoded and decoded in uint64-sized groups of digits:
	// maxDigits is how many base-len digits fit in a uint64, and maxDivisor is
	// len^maxDigits, the place value that shifts one full group in or out.
	maxDigits  int
	maxDivisor uint64

	// reverse maps a byte to its index in chars, with 255 marking bytes outside
	// the alphabet. It is nil when maxBytes > 1; a single-byte alphabet holds
	// only ASCII, so its indexes stay well below the marker.
	reverse *[256]byte
}

// maxPow calculates the maximum power of b that fits in a uint64, returning
// both the value (d = b^n) and the exponent n.
func maxPow(b uint64) (d uint64, n int) {
	d, n = b, 1
	for m := math.MaxUint64 / b; d <= m; {
		d *= b
		n++
	}
	return
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

	a := alphabet{
		chars:    abc,
		len:      int64(len(abc)),
		encLen:   uint8(math.Ceil(128 / math.Log2(float64(len(abc))))),
		maxBytes: uint8(utf8.RuneLen(abc[len(abc)-1])),
	}
	a.maxDivisor, a.maxDigits = maxPow(uint64(a.len))
	if a.maxBytes == 1 {
		a.reverse = new([256]byte)
		for i := range a.reverse {
			a.reverse[i] = 255
		}
		for i, c := range a.chars {
			a.reverse[c] = byte(i)
		}
	}
	return a
}

func (a *alphabet) Length() int64 {
	return a.len
}

// isDefault reports whether a holds exactly DefaultAlphabet. It compares bytes
// in place because converting chars to a string to compare it allocates, and
// costs more than building the alphabet did.
//
// maxBytes of 1 means every rune is single byte, chars being sorted, so the
// truncation to byte is safe only under that guard.
func (a *alphabet) isDefault() bool {
	if a.maxBytes != 1 || int(a.len) != len(DefaultAlphabet) {
		return false
	}
	for i, c := range a.chars {
		if byte(c) != DefaultAlphabet[i] {
			return false
		}
	}
	return true
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
		return 0, notInAlphabet(t)
	}
	return int64(i), nil
}
