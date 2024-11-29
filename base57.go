package shortuuid

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"math/bits"
	"strings"
)

type base57 struct {
	// alphabet is the character set to construct the UUID from.
	alphabet alphabet
}

const (
	strLen      = 22
	alphabetLen = 57
)

// Encode encodes uuid.UUID into a string using the most significant bits (MSB)
// first according to the alphabet.
func (b base57) Encode(u uuid.UUID) string {
	num := uint128{
		binary.BigEndian.Uint64(u[8:]),
		binary.BigEndian.Uint64(u[:8]),
	}
	var outIndexes [strLen]uint64

	for i := strLen - 1; num.Hi > 0 || num.Lo > 0; i-- {
		num, outIndexes[i] = num.quoRem64(alphabetLen)
	}

	var sb strings.Builder
	sb.Grow(strLen)
	for i := 0; i < strLen; i++ {
		sb.WriteRune(b.alphabet.chars[outIndexes[i]])
	}
	return sb.String()
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (b base57) Decode(s string) (u uuid.UUID, err error) {
	var n uint128
	var index int64

	for _, char := range s {
		index, err = b.alphabet.Index(char)
		if err != nil {
			return
		}
		n, err = n.mulAdd64(alphabetLen, uint64(index))
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
