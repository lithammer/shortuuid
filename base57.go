package shortuuid

import (
	"encoding/binary"
	"github.com/google/uuid"
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

	return b.numToString(num)
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (b base57) Decode(u string) (uuid.UUID, error) {
	buf, err := b.stringToNumBytes(u)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(buf)
}

// numToString converts a number a string using the given alphabet.
func (b *base57) numToString(number uint128) string {
	var digit uint64
	out := make([]rune, strLen)

	i := 0
	for number.Hi > 0 || number.Lo > 0 {
		number, digit = number.quoRem64(alphabetLen)
		out[i] = b.alphabet.chars[digit]
		i++
	}
	for i < strLen {
		out[i] = b.alphabet.chars[0]
		i++
	}

	reverse(out)

	return string(out)
}

// stringToNumBytes converts a string a number using the given alphabet.
func (b *base57) stringToNumBytes(s string) ([]byte, error) {
	var (
		n     uint128
		err   error
		index int64
	)

	for _, char := range s {
		n, err = n.mul64(uint64(b.alphabet.Length()))
		if err != nil {
			return nil, err
		}

		index, err = b.alphabet.Index(char)
		if err != nil {
			return nil, err
		}

		n, err = n.add64(uint64(index))
		if err != nil {
			return nil, err
		}
	}
	buf := make([]byte, 16)
	n.putBytes(buf)
	return buf, nil
}

// reverse reverses a inline.
func reverse(a []rune) {
	n := len(a)
	for i := 0; i < n/2; i++ {
		a[i], a[n-1-i] = a[n-1-i], a[i]
	}
}
