package shortuuid

import (
	"fmt"
	"strings"
)

// NewV5FromName shows a simple case-sensitive variant of the deprecated
// [NewWithNamespace] heuristic.
func ExampleNewV5_fromName() {
	NewV5FromName := func(name string) string {
		if name == "" {
			return NewV4()
		}

		if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") {
			return NewV5(NameSpaceURL, name)
		}

		if after, found := strings.CutPrefix(name, "urn:oid:"); found {
			return NewV5(NameSpaceOID, after)
		}

		if after, found := strings.CutPrefix(name, "x500:"); found {
			return NewV5(NameSpaceX500, after)
		}

		return NewV5(NameSpaceDNS, name)
	}

	fmt.Println(NewV5FromName("http://example.com"))
	fmt.Println(NewV5FromName("urn:oid:1.2.840.113549"))
	// Output:
	// T35fvrnVz6SMSdh9y5hs8c
	// HVizdopCKiLaGoTrVJrg9r
}
