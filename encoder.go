package shortuuid

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"unicode/utf8"
	"unsafe"

	"github.com/google/uuid"
)

// encoder is a generic encoder that can encode/decode UUIDs using any alphabet.
// It provides full support for custom alphabets including multibyte UTF-8 characters.
type encoder struct {
	// alphabet is the character set to construct the UUID from.
	alphabet alphabet
}

// maxPow calculates the maximum power of b that fits in a uint64, returning
// both the value (d = b^n) and the exponent n. This is used during encoding
// to process the 128-bit UUID value in chunks that fit in 64-bit arithmetic.
func maxPow(b uint64) (d uint64, n int) {
	d, n = b, 1
	for m := math.MaxUint64 / b; d <= m; {
		d *= b
		n++
	}
	return
}

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (e encoder) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r, ind uint64
	i := int(e.alphabet.encLen - 1)
	buf := make([]byte, int64(e.alphabet.encLen)*int64(e.alphabet.maxBytes))
	lastPlaced := len(buf)
	l := uint64(e.alphabet.len)
	d, n := maxPow(l)

	for num.Hi > 0 || num.Lo > 0 {
		num, r = num.quoRem64(d)
		for j := 0; j < n && i >= 0; j++ {
			r, ind = r/l, r%l
			c := e.alphabet.chars[ind]
			if e.alphabet.maxBytes == 1 {
				buf[i] = byte(c)
				lastPlaced--
			} else {
				lastPlaced -= utf8.EncodeRune(buf[lastPlaced-utf8.RuneLen(c):], c)
			}
			i--
		}
	}
	firstRuneLen := utf8.RuneLen(e.alphabet.chars[0])
	for ; i >= 0; i-- {
		lastPlaced -= utf8.EncodeRune(buf[lastPlaced-firstRuneLen:], e.alphabet.chars[0])
	}
	buf = buf[lastPlaced:]
	return unsafe.String(unsafe.SliceData(buf), len(buf)) // same as in strings.Builder
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (e encoder) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	var index int64

	for _, char := range s {
		index, err = e.alphabet.Index(char)
		if err != nil {
			return
		}
		n, err = n.mulAdd64(uint64(e.alphabet.len), uint64(index))
		if err != nil {
			return
		}
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return
}

const (
	b57MaxU64Digits  = 10
	b57MaxU64Divisor = 362033331456891249 // 57^10
)

// b57Encoder is an optimized encoder for the default base57 alphabet.
// It uses a specialized implementation that's faster than the generic encoder
// for the common case of base57 encoding/decoding.
type b57Encoder struct{}

func (e b57Encoder) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r uint64
	var buf [22]byte
	num, r = num.quoRem64(b57MaxU64Divisor)
	buf[21], r = DefaultAlphabet[r%57], r/57
	buf[20], r = DefaultAlphabet[r%57], r/57
	buf[19], r = DefaultAlphabet[r%57], r/57
	buf[18], r = DefaultAlphabet[r%57], r/57
	buf[17], r = DefaultAlphabet[r%57], r/57
	buf[16], r = DefaultAlphabet[r%57], r/57
	buf[15], r = DefaultAlphabet[r%57], r/57
	buf[14], r = DefaultAlphabet[r%57], r/57
	buf[13] = DefaultAlphabet[r%57]
	buf[12] = DefaultAlphabet[r/57]
	num, r = num.quoRem64(b57MaxU64Divisor)
	buf[11], r = DefaultAlphabet[r%57], r/57
	buf[10], r = DefaultAlphabet[r%57], r/57
	buf[9], r = DefaultAlphabet[r%57], r/57
	buf[8], r = DefaultAlphabet[r%57], r/57
	buf[7], r = DefaultAlphabet[r%57], r/57
	buf[6], r = DefaultAlphabet[r%57], r/57
	buf[5], r = DefaultAlphabet[r%57], r/57
	buf[4], r = DefaultAlphabet[r%57], r/57
	buf[3] = DefaultAlphabet[r%57]
	buf[2] = DefaultAlphabet[r/57]
	buf[1] = DefaultAlphabet[num.Lo%57]
	buf[0] = DefaultAlphabet[num.Lo/57]
	return unsafe.String(unsafe.SliceData(buf[:]), 22)
}

func (e b57Encoder) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	var n64, ind, i uint64

	for _, c := range s {
		if c > 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", c)
		}
		ind = uint64(reverseB57[c])
		if ind == 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", c)
		}
		n64 = n64*57 + ind
		i++
		if i == b57MaxU64Digits {
			n, err = n.mulAdd64(b57MaxU64Divisor, n64)
			if err != nil {
				return
			}
			i = 0
			n64 = 0
		}
	}
	n, err = n.mulAdd64(57*57, n64)
	if err != nil {
		return
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return
}

// uint128 represents a 128-bit unsigned integer as two 64-bit words.
// Lo contains the least significant 64 bits, and Hi contains the most
// significant 64 bits.
type uint128 struct {
	Lo, Hi uint64
}

// quoRem64 divides u by v and returns the quotient q and remainder r.
// The division is performed using 128-bit arithmetic, handling the
// high and low 64-bit words separately.
func (u uint128) quoRem64(v uint64) (q uint128, r uint64) {
	q.Hi, r = bits.Div64(0, u.Hi, v)
	q.Lo, r = bits.Div64(r, u.Lo, v)
	return
}

// mulAdd64 multiplies u by m and adds a, returning the result.
// Returns an error if the result would exceed 128 bits.
// This is used during base-N decoding to accumulate the decoded value.
func (u uint128) mulAdd64(m uint64, a uint64) (uint128, error) {
	hi, lo := bits.Mul64(u.Lo, m)
	p0, p1 := bits.Mul64(u.Hi, m)
	lo, c0 := bits.Add64(lo, a, 0)
	hi, c1 := bits.Add64(hi, p1, c0)
	if p0 != 0 || c1 != 0 {
		return uint128{}, fmt.Errorf("number is out of range (need a 128-bit value)")
	}
	return uint128{lo, hi}, nil
}

// reverseB57 is a lookup table for fast base57 decoding. It maps ASCII byte
// values (0-255) to their corresponding index in the default alphabet.
// A value of 255 indicates that the byte is not part of the alphabet.
// The table is indexed by the byte value directly, allowing O(1) lookup
// during decoding.
var reverseB57 = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 0, 1, 2, 3, 4, 5,
	6, 7, 255, 255, 255, 255, 255, 255,
	255, 8, 9, 10, 11, 12, 13, 14,
	15, 255, 16, 17, 18, 19, 20, 255,
	21, 22, 23, 24, 25, 26, 27, 28,
	29, 30, 31, 255, 255, 255, 255, 255,
	255, 32, 33, 34, 35, 36, 37, 38,
	39, 40, 41, 42, 255, 43, 44, 45,
	46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
}
