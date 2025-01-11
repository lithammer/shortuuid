package shortuuid

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"slices"
	"unicode/utf8"
	"unsafe"

	"github.com/google/uuid"
)

const (
	DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// DefaultEncoder is the default BaseNEncoder uses when generating new UUIDs, and is
// based on Base57.
var DefaultEncoder = b57Encoder{}

type BaseNEncoder []rune

// NewEncoder creates new BaseNEncoder with given alphabet
// Removes duplicates and sorts it to ensure reproducibility.
func NewEncoder(alphabet string) (BaseNEncoder, error) {
	e := BaseNEncoder(alphabet)
	slices.Sort(e)
	e = slices.Compact(e)
	if len(e) < 2 {
		return nil, fmt.Errorf("encoding alphabet must be at least two characters")
	}
	return e, nil
}

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (e BaseNEncoder) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r, ind uint64
	encLen := int(math.Ceil(128 / math.Log2(float64(len(e)))))
	maxBytes := utf8.RuneLen(e[len(e)-1])
	i := encLen - 1
	buf := make([]byte, encLen*maxBytes)
	lastPlaced := len(buf)
	l := uint64(len(e))
	d, n := maxPow(l)

	for i >= 0 {
		num, r = num.quoRem64(d)
		for j := 0; j < n && i >= 0; j++ {
			r, ind = r/l, r%l
			c := e[ind]
			if maxBytes == 1 {
				buf[i] = byte(c)
				lastPlaced--
			} else {
				lastPlaced -= utf8.EncodeRune(buf[lastPlaced-utf8.RuneLen(c):], c)
			}
			i--
		}
	}
	buf = buf[lastPlaced:]
	return unsafe.String(unsafe.SliceData(buf), len(buf)) // same as in strings.Builder
}

func (e BaseNEncoder) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	var index uint64

	for _, char := range s {
		index, err = bSearch(e, char)
		if err != nil {
			return
		}
		n, err = n.mulAdd64(uint64(len(e)), index)
		if err != nil {
			return
		}
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return
}

type b57Encoder struct{}

func (e b57Encoder) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r uint64
	var buf [22]byte
	num, r = num.quoRem64(362033331456891249)
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
	num, r = num.quoRem64(362033331456891249)
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
	var n64, ind1, ind2 uint64

	for i := 0; i < 10; i++ {
		if s[i] > 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i])
		}
		if s[i+10] > 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i+10])
		}
		ind1, ind2 = uint64(reverseB57[s[i]]), uint64(reverseB57[s[i+10]])
		if ind1 == 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i])
		}
		if ind2 == 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i+10])
		}
		n.Lo = n.Lo*57 + ind1
		n64 = n64*57 + ind2
	}
	n, err = n.mulAdd64(362033331456891249, n64)
	if err != nil {
		return
	}

	n64 = 0
	for i := 0; i < 2; i++ {
		if s[i+20] > 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i+20])
		}
		ind1 = uint64(reverseB57[s[i+20]])
		if ind1 == 255 {
			return u, fmt.Errorf("element '%v' is not part of the alphabet", s[i+20])
		}
		n64 = n64*57 + ind1
	}
	n, err = n.mulAdd64(57*57, n64)
	if err != nil {
		return
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return
}

func maxPow(b uint64) (d uint64, n int) {
	d, n = b, 1
	for m := math.MaxUint64 / b; d <= m; n++ {
		d *= b
	}
	return
}

// bSearch returns the index of the first instance of char in the alphabet, or an
// error if char is not present.
func bSearch(alphabet []rune, char rune) (uint64, error) {
	i, j := 0, len(alphabet)
	for i < j {
		h := int(uint(i+j) >> 1)
		if alphabet[h] < char {
			i = h + 1
		} else {
			j = h
		}
	}
	if i >= len(alphabet) || alphabet[i] != char {
		return 0, fmt.Errorf("element '%v' is not part of the alphabet", char)
	}
	return uint64(i), nil
}

type uint128 struct {
	Lo, Hi uint64
}

func (u uint128) quoRem64(v uint64) (q uint128, r uint64) {
	q.Hi, r = bits.Div64(0, u.Hi, v)
	q.Lo, r = bits.Div64(r, u.Lo, v)
	return
}

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
