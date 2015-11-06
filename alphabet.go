package shortuuid

import (
	"sort"
	"strings"
)

// DefaultAlphabet is the default alphabet used.
const DefaultAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

type alphabet struct {
	chars []string
	len   int64
}

func newAlphabet(s string) alphabet {
	// Remove duplicates and sort it to ensure reproducability.
	abc := dedupe(strings.Split(s, ""))
	sort.Strings(abc)
	return alphabet{
		chars: abc,
		len:   int64(len(abc)),
	}
}

func (a *alphabet) Length() int64 {
	return a.len
}

// Index returns the index of the first instance of t in the alphabet, or 0 if t is not present.
func (a *alphabet) Index(t string) int64 {
	for i, char := range a.chars {
		if char == t {
			return int64(i)
		}
	}
	return 0
}

// dudupe removes duplicate characters from s.
func dedupe(s []string) []string {
	var out []string
	m := make(map[string]bool)

	for _, char := range s {
		if _, ok := m[char]; !ok {
			m[char] = true
			out = append(out, char)
		}
	}

	return out
}
