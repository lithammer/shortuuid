package shortuuid

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/bits"
	"strings"
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
	if e.alphabet.singleBytes {
		return e.encodeSingleBytes(u)
	}
	return e.encode(u)
}

func (e encoder) encodeSingleBytes(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r uint64
	var i int
	var buf []byte
	if e.alphabet.len == defaultBase { // compiler optimizations using constants for default base
		buf = make([]byte, defaultEncLen)
		for i = defaultEncLen - 1; num.Hi > 0 || num.Lo > 0; {
			num, r = num.quoRem64(defaultDivisor)
			for j := 0; j < defaultNDigits && i >= 0; j++ {
				buf[i] = byte(e.alphabet.chars[r%defaultBase])
				r /= defaultBase
				i--
			}
		}
	} else {
		buf = make([]byte, e.alphabet.encLen)
		l := uint64(e.alphabet.len)
		d, n := maxPow(l)
		for i = int(e.alphabet.encLen - 1); num.Hi > 0 || num.Lo > 0; {
			num, r = num.quoRem64(d)
			for j := 0; j < n && i >= 0; j++ {
				buf[i] = byte(e.alphabet.chars[r%l])
				r /= l
				i--
			}
		}
	}
	for ; i >= 0; i-- {
		buf[i] = byte(e.alphabet.chars[0])
	}
	return string(buf[:])
}

func (e encoder) encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var r uint64
	var outIndexes []uint64
	if e.alphabet.len == defaultBase { // compiler optimizations using constants for default base
		outIndexes = make([]uint64, defaultEncLen) // avoids escaping to heap for base57 when used with constant
		for i := defaultEncLen - 1; num.Hi > 0 || num.Lo > 0; {
			num, r = num.quoRem64(defaultDivisor)
			for j := 0; j < defaultNDigits && i >= 0; j++ {
				outIndexes[i] = r % defaultBase
				r /= defaultBase
				i--
			}
		}
	} else {
		outIndexes = make([]uint64, e.alphabet.encLen)
		l := uint64(e.alphabet.len)
		d, n := maxPow(l)
		for i := int(e.alphabet.encLen - 1); num.Hi > 0 || num.Lo > 0; {
			num, r = num.quoRem64(d)
			for j := 0; j < n && i >= 0; j++ {
				outIndexes[i] = r % l
				r /= l
				i--
			}
		}
	}

	var sb strings.Builder
	sb.Grow(int(e.alphabet.encLen))
	for i := 0; i < int(e.alphabet.encLen); i++ {
		sb.WriteRune(e.alphabet.chars[outIndexes[i]])
	}
	return sb.String()
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
