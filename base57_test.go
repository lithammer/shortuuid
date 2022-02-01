package shortuuid

import (
	"testing"
)

func TestReverse(t *testing.T) {
	a := []rune("abc123")
	reverse(a)
	if string(a) != "321cba" {
		t.Errorf("expected string to be %q, got %q", "321cba", string(a))
	}
}
