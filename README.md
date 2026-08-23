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

`New` returns a v4 (random) UUID. Name a version explicitly when you want a
different one:

```go
shortuuid.NewV4() // random, same as New
shortuuid.NewV7() // time-ordered
```

v7 UUIDs lead with a millisecond timestamp, and base57 keeps that order, so
sorting v7 shortuuids as strings sorts them by creation time.

For v5 (derived from a name rather than random) pass a namespace and a name:

```go
shortuuid.NewV5(shortuuid.NameSpaceDNS, "example.com")          // exu3DTbj2ncsn9tLdLWspw
shortuuid.NewV5(shortuuid.NameSpaceURL, "http://example.com")   // T35fvrnVz6SMSdh9y5hs8c
shortuuid.NewV5(shortuuid.NameSpaceOID, "1.2.840.113549")       // HVizdopCKiLaGoTrVJrg9r
shortuuid.NewV5(shortuuid.NameSpaceX500, "CN=example,O=org")
```

Pass the bare name — an OID is `"1.2.840.113549"`, not `"urn:oid:1.2.840.113549"`.

Decoding turns a shortuuid back into a UUID:

```go
u, err := shortuuid.DefaultEncoder.Decode("KwSysDpxcBU9FNhGkn2dCf")
// 64d1355f-d052-4bd9-83f4-39b93fb1c01f
```

A custom alphabet (at least 2 distinct characters) needs its own encoder. Sorting
and deduplicating the alphabet happens up front, so build the encoder once and
reuse it:

```go
var enc = shortuuid.NewEncoder("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=")

enc.Encode(uuid.New())
```

The encoder takes any UUID, so one encoder covers every version:

```go
enc.Encode(uuid.NewV7())
enc.Encode(shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com")) // dwt3CSai2mbrm9sKcKVrov
```

Decode with the same encoder, since the alphabet has to match:

```go
u, err := enc.Decode("dwt3CSai2mbrm9sKcKVrov")
// cfbff0d1-9375-5685-968c-48ce8b15ae17
```

Bring your own encoder! For example, base58 is popular among cryptocurrencies
like Bitcoin.

```go
package main

import (
	"fmt"
	"uuid"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lithammer/shortuuid/v5"
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

<details>
<summary>Migrating from v4</summary>

`New` is unchanged, so code calling it only needs its import path bumped to
`/v5`.

`Encoder` now names the standard library `uuid.UUID` rather than
`github.com/google/uuid.UUID`, so a custom encoder stops satisfying the
interface until it imports `uuid` instead. Both types are `[16]byte`, so a
google/uuid value converts with `uuid.UUID(v)`.

`NewWithEncoder` and `NewWithAlphabet` still work, but are deprecated. Use
`NewEncoder` instead. `NewWithEncoder(enc)` is `enc.Encode(uuid.New())` at the
same cost, and encoding directly also reaches the other UUID versions.
`NewWithAlphabet` re-sorts the alphabet on every call with no way to avoid it,
so a reused encoder is more than twice as fast:

```go
shortuuid.NewWithAlphabet(abc)        // deprecated

var enc = shortuuid.NewEncoder(abc)   // hoisted once
enc.Encode(uuid.New())
```

`NewWithNamespace(name)` still works and its IDs are unchanged, but it is
deprecated. It guesses the namespace from the name: `NameSpaceURL` for
`http://` and `https://` prefixes, matched case-insensitively, and
`NameSpaceDNS` for everything else, which is what the Python `shortuuid`
library does. An OID or an X.500 DN hashes under `NameSpaceDNS` too.

Prefer `NewV5` with the namespace you mean:

```go
shortuuid.NewWithNamespace("http://example.com")                // deprecated
shortuuid.NewV5(shortuuid.NameSpaceURL, "http://example.com")   // same ID, explicit
```

Reproducing the guess takes a small helper:

```go
func fromName(name string) string {
	if name == "" {
		return shortuuid.New()
	}

	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return shortuuid.NewV5(shortuuid.NameSpaceURL, name)
	}
	return shortuuid.NewV5(shortuuid.NameSpaceDNS, name)
}
```

v5 requires Go 1.27 and has no dependencies.

</details>

## License

MIT
