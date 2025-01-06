package shortuuid

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/bits"
	"unicode/utf8"
	"unsafe"
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

	tx = 0b10000000
	t2 = 0b11000000
	t3 = 0b11100000
	t4 = 0b11110000

	maskx = 0b00111111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	runeErrorByte0 = t3 | (utf8.RuneError >> 12)
	runeErrorByte1 = tx | (utf8.RuneError>>6)&maskx
	runeErrorByte2 = tx | utf8.RuneError&maskx
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
	buf := make([]byte, defaultEncLen)
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
	return string(buf)
}

func (e encoder) encode(num uint128) string {
	var r, ind uint64
	i := e.alphabet.encLen - 1
	buf := make([]byte, e.alphabet.encLen*e.alphabet.maxBytes)
	curByteInd := len(buf) - 1
	l := uint64(e.alphabet.len)
	d, n := maxPow(l)

	for num.Hi > 0 || num.Lo > 0 {
		num, r = num.quoRem64(d)
		for j := 0; j < n && i >= 0; j++ {
			r, ind = r/l, r%l
			curByteInd -= placeRuneEndingAt(buf, e.alphabet.chars[ind], curByteInd)
			i--
		}
	}
	for ; i >= 0; i-- {
		curByteInd -= placeRuneEndingAt(buf, e.alphabet.chars[0], curByteInd)
	}
	buf = buf[curByteInd+1:]
	return unsafe.String(unsafe.SliceData(buf), len(buf)) // same as in strings.Builder
}

func placeRuneEndingAt(p []byte, r rune, ind int) int {
	switch i := uint32(r); {
	case i <= rune1Max:
		p[ind] = byte(r)
		return 1
	case i <= rune2Max:
		p[ind] = tx | byte(r)&maskx
		p[ind-1] = t2 | byte(r>>6)
		return 2
	case i < surrogateMin, surrogateMax < i && i <= rune3Max:
		p[ind] = tx | byte(r)&maskx
		p[ind-1] = tx | byte(r>>6)&maskx
		p[ind-2] = t3 | byte(r>>12)
		return 3
	case i > rune3Max && i <= utf8.MaxRune:
		p[ind] = tx | byte(r)&maskx
		p[ind-1] = tx | byte(r>>6)&maskx
		p[ind-2] = tx | byte(r>>12)&maskx
		p[ind-3] = t4 | byte(r>>18)
		return 4
	default:
		p[ind] = runeErrorByte2
		p[ind-1] = runeErrorByte1
		p[ind-2] = runeErrorByte0
		return 3
	}
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
