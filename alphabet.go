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
const DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// alphabet represents a character set for base-N encoding. It stores the
// sorted, deduplicated characters along with precomputed values for efficient
// encoding and decoding.
type alphabet struct {
	chars    []rune // sorted, deduplicated characters
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

	// mb carries the multibyte lookup tables behind one pointer, keeping the
	// struct the hot paths copy small. It is nil when maxBytes == 1.
	mb *mbTables
}

type mbTables struct {
	// packed[i] holds the UTF-8 encoding of chars[i]: the encoded bytes occupy
	// the top of the low 32-bit little-endian word (so a 4-byte store ending at
	// the write position lands them right-aligned), and bits 32-39 hold the
	// byte length.
	packed []uint64

	// runeIdx maps c - minRune to the index of c in chars, with 255 marking
	// runes outside the alphabet. Only set when the rune range is small
	// enough; decode falls back to binary search otherwise.
	runeIdx []byte
	minRune rune
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
		encLen:   uint8(math.Ceil(128 / math.Log2(float64(len(abc))))),
		maxBytes: uint8(utf8.RuneLen(abc[len(abc)-1])),
	}
	a.maxDivisor, a.maxDigits = maxPow(uint64(len(abc)))
	if a.maxBytes == 1 {
		a.reverse = new([256]byte)
		for i := range a.reverse {
			a.reverse[i] = 255
		}
		for i, c := range a.chars {
			a.reverse[c] = byte(i)
		}
	} else {
		mb := &mbTables{packed: make([]uint64, len(a.chars))}
		for i, c := range a.chars {
			var tmp [4]byte
			sz := utf8.EncodeRune(tmp[:], c)
			v := uint64(sz) << 32
			for j := range sz {
				v |= uint64(tmp[j]) << (8 * (4 - sz + j))
			}
			mb.packed[i] = v
		}
		mb.minRune = a.chars[0]
		if spread := a.chars[len(a.chars)-1] - mb.minRune + 1; spread <= 4096 && len(a.chars) <= 255 {
			mb.runeIdx = make([]byte, spread)
			for i := range mb.runeIdx {
				mb.runeIdx[i] = 255
			}
			for i, c := range a.chars {
				mb.runeIdx[c-mb.minRune] = byte(i)
			}
		}
		a.mb = mb
	}
	return a
}

// isDefault reports whether a holds exactly DefaultAlphabet. Both are sorted
// and deduplicated, so comparing elementwise is enough. The comparison stays
// in place because converting chars to a string allocates, and costs more
// than building the alphabet did.
func (a *alphabet) isDefault() bool {
	if len(a.chars) != len(DefaultAlphabet) {
		return false
	}
	for i, c := range a.chars {
		if c != rune(DefaultAlphabet[i]) {
			return false
		}
	}
	return true
}

// Index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (a *alphabet) Index(t rune) (int, error) {
	i, j := 0, len(a.chars)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a.chars[h] < t {
			i = h + 1
		} else {
			j = h
		}
	}
	if i >= len(a.chars) || a.chars[i] != t {
		return 0, notInAlphabet(t)
	}
	return i, nil
}
