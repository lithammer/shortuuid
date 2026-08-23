package shortuuid_test

import (
	"fmt"

	"github.com/lithammer/shortuuid/v5"
)

// Encoding a v5 UUID over base 16 reproduces the UUID's own hex digits, which
// is a useful thing to check an alphabet against.
func ExampleNewEncoder() {
	enc := shortuuid.NewEncoder("0123456789abcdef")

	u := shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com")
	fmt.Println(u)
	fmt.Println(enc.Encode(u))

	// Output:
	// cfbff0d1-9375-5685-968c-48ce8b15ae17
	// cfbff0d193755685968c48ce8b15ae17
}
