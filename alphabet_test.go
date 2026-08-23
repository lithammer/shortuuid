package shortuuid

import (
	"math/big"
	"testing"
)

func TestAlphabetIndex(t *testing.T) {
	abc := newAlphabet(DefaultAlphabet)
	tests := []struct {
		char rune
		want int
		err  bool
	}{
		{'2', 0, false},
		{'z', 56, false},
		{'l', 0, true},
	}
	for _, test := range tests {
		idx, err := abc.Index(test.char)
		if test.err {
			if err == nil {
				t.Errorf("Index(%q) returned %d, want an error", test.char, idx)
			}
			continue
		}
		if err != nil {
			t.Errorf("Index(%q) returned %v, want %d", test.char, err, test.want)
		} else if idx != test.want {
			t.Errorf("Index(%q) = %d, want %d", test.char, idx, test.want)
		}
	}
}

func TestEncLenExact(t *testing.T) {
	// encLen comes from float math, and Encode trusts it as the exact digit
	// count for any 128-bit value: one too small would silently drop the most
	// significant digits. Pin every base to the integer ground truth, the
	// smallest n with base^n >= 2^128. The power-of-two bases are the ones
	// that land exactly on the ceil boundary.
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	for base := 2; base <= 4096; base++ {
		chars := make([]rune, base)
		for i := range chars {
			chars[i] = rune(0x4E00 + i)
		}
		got := int(newAlphabet(string(chars)).encLen)

		want := 0
		for v := big.NewInt(1); v.Cmp(limit) < 0; want++ {
			v.Mul(v, big.NewInt(int64(base)))
		}
		if got != want {
			t.Errorf("encLen for base %d = %d, want %d", base, got, want)
		}
	}
}
