# shortuuid

[![Build Status](https://github.com/lithammer/shortuuid/workflows/CI/badge.svg)](https://github.com/lithammer/shortuuid/actions)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/lithammer/shortuuid/v5)

A Go library that generates concise, unambiguous, URL-safe UUIDs. Based on and
compatible with the Python library
[`shortuuid`](https://github.com/skorokithakis/shortuuid).

Often, one needs to use non-sequential IDs in places where users will see them,
but the IDs must be as concise and easy to use as possible. shortuuid solves
this problem by generating UUIDs using the standard library
[`uuid`](https://pkg.go.dev/uuid) package (Go 1.27+) and then translating them to
base57 using lowercase and uppercase letters and digits, and removing
similar-looking characters such as l, 1, I, O and 0.

## Usage

```go
package main

import (
	"fmt"

	"github.com/lithammer/shortuuid/v5"
)

func main() {
	u := shortuuid.New()
	fmt.Println(u) // KwSysDpxcBU9FNhGkn2dCf
}
```

`New` returns a v4 (random) UUID and is the right default. Name a version
explicitly when you need one:

```go
shortuuid.NewV4() // random, same as New
shortuuid.NewV7() // time-ordered
```

v7 UUIDs lead with a millisecond timestamp, and base57 preserves that order, so
v7 shortuuids sort lexicographically in creation order.

For v5 (derived from a name rather than random) pass a namespace and a name:

```go
shortuuid.NewV5(shortuuid.NameSpaceDNS, "example.com")          // exu3DTbj2ncsn9tLdLWspw
shortuuid.NewV5(shortuuid.NameSpaceURL, "http://example.com")   // T35fvrnVz6SMSdh9y5hs8c
shortuuid.NewV5(shortuuid.NameSpaceOID, "1.2.840.113549")       // HVizdopCKiLaGoTrVJrg9r
shortuuid.NewV5(shortuuid.NameSpaceX500, "CN=example,O=org")
```

Pass the bare name — an OID is `"1.2.840.113549"`, not `"urn:oid:1.2.840.113549"`.

A custom alphabet (at least 2 characters long) needs its own encoder. Build it
once and reuse it — the alphabet is sorted and deduplicated up front, which
costs more than encoding does:

```go
var enc = shortuuid.NewEncoder("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=")

enc.Encode(uuid.New())   // iZsai==fWebXd5rLRWFB=u
```

The encoder takes any UUID, so one encoder covers every version:

```go
enc.Encode(uuid.NewV7())
enc.Encode(shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com")) // dwt3CSai2mbrm9sKcKVrov
```

<details>
<summary>Migrating from v4</summary>

`New` is unchanged, so code calling it only needs its import path bumped to
`/v5`.

`NewWithEncoder` and `NewWithAlphabet` still work, but are deprecated in favour
of `NewEncoder`. `NewWithEncoder(enc)` is `enc.Encode(uuid.New())` at the same
cost, and encoding directly also reaches the other UUID versions.
`NewWithAlphabet` re-sorts the alphabet on every call and offers no way to avoid
it, so a reused encoder is roughly 2.5x faster:

```go
shortuuid.NewWithAlphabet(abc)        // deprecated

var enc = shortuuid.NewEncoder(abc)   // hoisted once
enc.Encode(uuid.New())
```

`NewWithNamespace(name)` still works and its IDs are unchanged, but it is
deprecated. It guesses the namespace from the name — `NameSpaceURL` for
`http://` and `https://` prefixes, matched case-insensitively, `NameSpaceDNS`
for everything else — which mirrors the Python `shortuuid` library. Names that
are neither hostnames nor URLs, OIDs and X.500 DNs included, are hashed under
`NameSpaceDNS`.

Prefer `NewV5` with the namespace you actually mean:

```go
shortuuid.NewWithNamespace("http://example.com")                // deprecated
shortuuid.NewV5(shortuuid.NameSpaceURL, "http://example.com")   // same ID, explicit
```

To keep the guessing behaviour, it is four lines:

```go
func fromName(name string) string {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return shortuuid.NewV5(shortuuid.NameSpaceURL, name)
	}
	return shortuuid.NewV5(shortuuid.NameSpaceDNS, name)
}
```

shortuuid now requires Go 1.27, and has no dependencies.

</details>

Bring your own encoder! For example, base58 is popular among bitcoin.

```go
package main

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lithammer/shortuuid/v5"
	"uuid"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	b := base58.Decode(s)
	if len(b) != 16 {
		return uuid.UUID{}, fmt.Errorf("invalid UUID (got %d bytes)", len(b))
	}
	return uuid.UUID(b), nil
}

func main() {
	enc := base58Encoder{}
	fmt.Println(shortuuid.NewWithEncoder(enc)) // 6R7VqaQHbzC1xwA5UueGe6
}
```

## License

MIT
