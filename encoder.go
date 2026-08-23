package shortuuid

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"unicode/utf8"
	"unsafe"
	"uuid"
)

// encoder encodes and decodes UUIDs over any alphabet, multibyte included.
type encoder struct {
	alphabet alphabet
}

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (e encoder) Encode(u uuid.UUID) string {
	num := uint128{
		Lo: binary.BigEndian.Uint64(u[8:]),
		Hi: binary.BigEndian.Uint64(u[:8]),
	}
	if e.alphabet.maxBytes == 1 {
		return e.encodeBytes(num)
	}
	return e.encodeRunes(num)
}

// encodeBytes encodes over a single-byte alphabet, where every digit is one
// byte wide, so the output is exactly encLen bytes and the write position is
// the digit position.
func (e encoder) encodeBytes(num uint128) string {
	chars := e.alphabet.chars
	l := uint64(len(chars))
	d, n := e.alphabet.maxDivisor, e.alphabet.maxDigits
	buf := make([]byte, e.alphabet.encLen)
	i := len(buf) - 1
	var r uint64

	for num.Hi > 0 {
		num, r = num.quoRem64(d)
		for j := 0; j < n && i >= 0; j++ {
			buf[i] = byte(chars[r%l])
			r /= l
			i--
		}
	}
	for r = num.Lo; r > 0 && i >= 0; i-- {
		buf[i] = byte(chars[r%l])
		r /= l
	}
	for ; i >= 0; i-- {
		buf[i] = byte(chars[0])
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf)) // same as in strings.Builder
}

// encodeRunes encodes over a multibyte alphabet. Runes have varying widths,
// so digits are written back to front at lastPlaced while i counts the digits
// still to place; each digit is one 4-byte store of its precomputed UTF-8
// encoding, and the buffer is sized for the widest case and trimmed.
func (e encoder) encodeRunes(num uint128) string {
	packed := e.alphabet.mb.packed
	l := uint64(len(packed))
	d, n := e.alphabet.maxDivisor, e.alphabet.maxDigits
	var r uint64
	i := int(e.alphabet.encLen - 1)
	// 3 front-pad bytes so the 4-byte stores below stay in bounds when the
	// leftmost character sits at the start of the encoded output.
	buf := make([]byte, 3+int64(e.alphabet.encLen)*int64(e.alphabet.maxBytes))
	lastPlaced := len(buf)

	for num.Hi > 0 {
		num, r = num.quoRem64(d)
		for j := 0; j < n && i >= 0; j++ {
			v := packed[r%l]
			r /= l
			binary.LittleEndian.PutUint32(buf[lastPlaced-4:], uint32(v))
			lastPlaced -= int(v >> 32)
			i--
		}
	}
	for r = num.Lo; r > 0 && i >= 0; i-- {
		v := packed[r%l]
		r /= l
		binary.LittleEndian.PutUint32(buf[lastPlaced-4:], uint32(v))
		lastPlaced -= int(v >> 32)
	}
	v0 := packed[0]
	for ; i >= 0; i-- {
		binary.LittleEndian.PutUint32(buf[lastPlaced-4:], uint32(v0))
		lastPlaced -= int(v0 >> 32)
	}
	buf = buf[lastPlaced:]
	return unsafe.String(unsafe.SliceData(buf), len(buf)) // same as in strings.Builder
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
//
// Digits accumulate in a uint64 group that folds into the 128-bit result once
// per maxDigits characters, so the 128-bit arithmetic runs a handful of times
// instead of once per character.
func (e encoder) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	if e.alphabet.maxBytes == 1 {
		n, err = e.decodeBytes(s)
	} else {
		n, err = e.decodeRunes(s)
	}
	if err != nil {
		return u, err
	}
	binary.BigEndian.PutUint64(u[:8], n.Hi)
	binary.BigEndian.PutUint64(u[8:], n.Lo)
	return u, nil
}

// decodeBytes decodes over a single-byte alphabet, indexing the reverse table
// with raw bytes. Any byte of a multibyte character maps to 255 in the table,
// so such input fails the same way any character outside the alphabet does.
func (e encoder) decodeBytes(s string) (n uint128, err error) {
	l := uint64(len(e.alphabet.chars))
	reverse := e.alphabet.reverse
	var group uint64
	var digits int
	for i := range len(s) {
		ind := reverse[s[i]]
		if ind == 255 {
			r, _ := utf8.DecodeRuneInString(s[i:])
			return n, notInAlphabet(r)
		}
		group = group*l + uint64(ind)
		if digits++; digits == e.alphabet.maxDigits {
			if n, err = n.mulAdd64(e.alphabet.maxDivisor, group); err != nil {
				return n, err
			}
			group, digits = 0, 0
		}
	}
	return n.mulAddDigits(l, group, digits)
}

// decodeRunes decodes over a multibyte alphabet, looking runes up in the
// rune-index table when the alphabet was compact enough to build one, and by
// binary search otherwise.
func (e encoder) decodeRunes(s string) (n uint128, err error) {
	l := uint64(len(e.alphabet.chars))
	var ind int
	var group uint64
	var digits int
	if t := e.alphabet.mb.runeIdx; t != nil {
		minRune := e.alphabet.mb.minRune
		for _, c := range s {
			d := uint32(c - minRune)
			if d >= uint32(len(t)) || t[d] == 255 {
				return n, notInAlphabet(c)
			}
			group = group*l + uint64(t[d])
			if digits++; digits == e.alphabet.maxDigits {
				if n, err = n.mulAdd64(e.alphabet.maxDivisor, group); err != nil {
					return n, err
				}
				group, digits = 0, 0
			}
		}
		return n.mulAddDigits(l, group, digits)
	}
	for _, c := range s {
		ind, err = e.alphabet.Index(c)
		if err != nil {
			return n, err
		}
		group = group*l + uint64(ind)
		if digits++; digits == e.alphabet.maxDigits {
			if n, err = n.mulAdd64(e.alphabet.maxDivisor, group); err != nil {
				return n, err
			}
			group, digits = 0, 0
		}
	}
	return n.mulAddDigits(l, group, digits)
}

// errOutOfRange reports a decoded value too large for the 128 bits a UUID has.
var errOutOfRange = errors.New("number is out of range (need a 128-bit value)")

// b57MaxU64Divisor is 57^10, the largest power of 57 that fits in a uint64.
const b57MaxU64Divisor = 362033331456891249

// b57Encoder is an optimized encoder for the default base57 alphabet. Encode
// is specialized so the compiler turns every division by 57 into cheaper
// multiplications; Decode gains nothing from that and delegates to the
// generic encoder.
type b57Encoder struct{}

// genericB57 is the generic encoder over DefaultAlphabet, carrying the
// reverse table b57Encoder.Decode needs.
var genericB57 = encoder{newAlphabet(DefaultAlphabet)}

func (e b57Encoder) Encode(u uuid.UUID) string {
	num := uint128{
		Lo: binary.BigEndian.Uint64(u[8:]),
		Hi: binary.BigEndian.Uint64(u[:8]),
	}
	var r uint64
	var buf [22]byte
	// The 22 digits split 10+10+2: quoRem64 pulls out two full uint64 groups
	// and num.Lo keeps the top pair, since 2^128/57^20 < 57^2. Within a group,
	// eight chained divisions leave r < 57^2, so the last two writes need no
	// shift.
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
	return unsafe.String(unsafe.SliceData(buf[:]), 22) // same as in strings.Builder
}

func (e b57Encoder) Decode(s string) (uuid.UUID, error) {
	return genericB57.Decode(s)
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
		return uint128{}, errOutOfRange
	}
	return uint128{lo, hi}, nil
}

// mulAddDigits folds a partial group of digits into u, scaling u by
// base^digits. A full group scales by the precomputed maxDivisor instead;
// this recomputes the power because a partial group occurs at most once per
// decode.
func (u uint128) mulAddDigits(base, group uint64, digits int) (uint128, error) {
	pow := uint64(1)
	for range digits {
		pow *= base
	}
	return u.mulAdd64(pow, group)
}
