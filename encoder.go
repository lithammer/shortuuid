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

type encoder []rune

// DefaultEncoder is the default encoder uses when generating new UUIDs, and is
// based on Base57.
var DefaultEncoder = encoder(DefaultAlphabet)

// NewEncoder creates new encoder with given alphabet
// Remove duplicates and sort it to ensure reproducibility.
func NewEncoder(alphabet string) Encoder {
	e := encoder(alphabet)
	slices.Sort(e)
	e = slices.Compact(e)
	if len(e) < 2 {
		panic("encoding alphabet must be at least two characters")
	}
	return e
}

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (e encoder) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	if unsafe.SliceData(e) == unsafe.SliceData(DefaultEncoder) {
		return defaultEncode(num)
	}
	return e.encode(num)
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (e encoder) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	var index uint64
	l := uint64(len(e))
	for _, char := range s {
		index, err = e.index(char)
		if err != nil {
			return
		}
		n, err = n.mulAdd64(l, index)
		if err != nil {
			return
		}
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return
}

func defaultEncode(num uint128) string {
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
	buf[13], r = DefaultAlphabet[r%57], r/57
	buf[12] = DefaultAlphabet[r%57]
	num, r = num.quoRem64(362033331456891249)
	buf[11], r = DefaultAlphabet[r%57], r/57
	buf[10], r = DefaultAlphabet[r%57], r/57
	buf[9], r = DefaultAlphabet[r%57], r/57
	buf[8], r = DefaultAlphabet[r%57], r/57
	buf[7], r = DefaultAlphabet[r%57], r/57
	buf[6], r = DefaultAlphabet[r%57], r/57
	buf[5], r = DefaultAlphabet[r%57], r/57
	buf[4], r = DefaultAlphabet[r%57], r/57
	buf[3], r = DefaultAlphabet[r%57], r/57
	buf[2] = DefaultAlphabet[r%57]
	_, r = num.quoRem64(362033331456891249)
	buf[1], r = DefaultAlphabet[r%57], r/57
	buf[0] = DefaultAlphabet[r%57]
	return unsafe.String(unsafe.SliceData(buf[:]), 22)
}

func maxPow(b uint64) (d uint64, n int) {
	d, n = b, 1
	for m := math.MaxUint64 / b; d <= m; n++ {
		d *= b
	}
	return
}

func (e encoder) encode(num uint128) string {
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

// index returns the index of the first instance of t in the alphabet, or an
// error if t is not present.
func (e encoder) index(t rune) (uint64, error) {
	i, j := 0, len(e)
	for i < j {
		h := int(uint(i+j) >> 1)
		if e[h] < t {
			i = h + 1
		} else {
			j = h
		}
	}
	if i >= len(e) || e[i] != t {
		return 0, fmt.Errorf("element '%v' is not part of the alphabet", t)
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
