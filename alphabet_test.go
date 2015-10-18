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

func TestAlphabetIndexZero(t *testing.T) {
	abc := newAlphabet("abc")
	if abc.Index("z") != 0 {
		t.Errorf("expected index 0, got %d", abc.Index("z"))
	}
}
