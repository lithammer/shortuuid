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

type encoder struct {
	// alphabet is the character set to construct the UUID from.
	alphabet alphabet
}

const (
	defaultBase    = 57
	defaultEncLen  = 22
	defaultNDigits = 10
	defaultDivisor = 362033331456891249 // 57^10
)

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
	if e.alphabet.len == defaultBase && e.alphabet.maxBytes == 1 {
		return e.defaultEncode(num)
	}
	return e.encode(num)
}

func (e encoder) defaultEncode(num uint128) string { // compiler optimizes a lot of divisions by constant
	var i int
	var r uint64
	var buf [defaultEncLen]byte
	for i = defaultEncLen - 1; num.Hi > 0 || num.Lo > 0; {
		num, r = num.quoRem64(defaultDivisor)
		for j := 0; j < defaultNDigits && i >= 0; j++ {
			buf[i] = byte(e.alphabet.chars[r%defaultBase])
			r /= defaultBase
			i--
		}
	}
	for ; i >= 0; i-- {
		buf[i] = byte(e.alphabet.chars[0])
	}
	return string(buf[:])
}

func (e encoder) encode(num uint128) string {
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
			lastPlaced -= utf8.EncodeRune(buf[lastPlaced-utf8.RuneLen(c):], c)
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
