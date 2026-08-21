package shortuuid_test

import (
	"fmt"
	"uuid"

	"github.com/lithammer/shortuuid/v5"
)

// A custom alphabet needs its own encoder. Build it once — sorting and
// deduplicating the alphabet costs more than the encoding — then hand it any
// UUID.
func ExampleNewEncoder() {
	enc := shortuuid.NewEncoder("0123456789abcdef")

	// Base 16 over 128 bits is just the UUID's own hex digits.
	u := shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com")
	fmt.Println(u)
	fmt.Println(enc.Encode(u))

	// The same encoder serves every version.
	fmt.Println(len(enc.Encode(uuid.New())), len(enc.Encode(uuid.NewV7())))

	// Output:
	// cfbff0d1-9375-5685-968c-48ce8b15ae17
	// cfbff0d193755685968c48ce8b15ae17
	// 32 32
}
