package shortuuid

import (
	"strings"
	"testing"
)

func TestDedupe(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"01010101010101", "01"},
		{"abcabcfoo", "abcfo"},
	}

	for _, test := range tests {
		in := strings.Join(dedupe(strings.Split(test.in, "")), "")
		if in != test.out {
			t.Errorf("expected %q, got %q", in, test.out)
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
