package shortuuid

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type uint128 struct {
	Lo, Hi uint64
}

func (u uint128) quoRem64(v uint64) (q uint128, r uint64) {
	if u.Hi < v {
		q.Lo, r = bits.Div64(u.Hi, u.Lo, v)
	} else {
		q.Hi, r = bits.Div64(0, u.Hi, v)
		q.Lo, r = bits.Div64(r, u.Lo, v)
	}
	return
}

func (u uint128) mul64(v uint64) (uint128, error) {
	hi, lo := bits.Mul64(u.Lo, v)
	p0, p1 := bits.Mul64(u.Hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		return uint128{}, fmt.Errorf("number is out of range (need a 128-bit value)")
	}
	return uint128{lo, hi}, nil
}

func (u uint128) add64(v uint64) (uint128, error) {
	lo, carry := bits.Add64(u.Lo, v, 0)
	hi, carry := bits.Add64(u.Hi, 0, carry)
	if carry != 0 {
		return uint128{}, fmt.Errorf("number is out of range (need a 128-bit value)")
	}
	return uint128{lo, hi}, nil
}

func (u uint128) putBytes(b []byte) {
	binary.BigEndian.PutUint64(b[:8], u.Hi)
	binary.BigEndian.PutUint64(b[8:], u.Lo)
}
