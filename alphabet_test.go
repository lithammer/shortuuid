package shortuuid

import (
	"sort"
	"testing"
)

func TestDedupe(t *testing.T) {
	tests := []struct {
		in, out []rune
	}{
		{[]rune("01010101010101"), []rune("01")},
		{[]rune("abcabcfoo"), []rune("abcfo")},
	}

	for _, test := range tests {
		sort.Slice(test.in, func(i, j int) bool { return test.in[i] < test.in[j] })
		in := dedupe(test.in)
		if string(in) != string(test.out) {
			t.Errorf("expected %q, got %q", string(test.out), string(in))
		}
	}
}

func TestAlphabetIndex(t *testing.T) {
	abc := newAlphabet(DefaultAlphabet)
	idx, err := abc.Index('z')
	if err != nil {
		t.Errorf("expected index 56, got an error trying to get it %v", err)
	}
	if idx != 56 {
		t.Errorf("expected index 56, got %d", idx)
	}
}

func TestAlphabetIndexZero(t *testing.T) {
	abc := newAlphabet(DefaultAlphabet)
	idx, err := abc.Index('2')
	if err != nil {
		t.Errorf("expected index 0, got an error trying to get it %v", err)
	}
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}
}

func TestAlphabetIndexError(t *testing.T) {
	abc := newAlphabet(DefaultAlphabet)
	idx, err := abc.Index('l')
	if err == nil {
		t.Errorf("expected an error, got a valid index %d", idx)
	}
}

func BenchmarkAlphabetIndex(b *testing.B) {
	abc := newAlphabet(DefaultAlphabet)
	for i := 0; i < b.N; i++ {
		for _, ch := range DefaultAlphabet {
			_, _ = abc.Index(ch)
		}
	}
}
