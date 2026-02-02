# shortuuid

[![Build Status](https://github.com/lithammer/shortuuid/workflows/CI/badge.svg)](https://github.com/lithammer/shortuuid/actions)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/lithammer/shortuuid/v4)

A Go library that generates concise, unambiguous, URL-safe UUIDs. Based on and
compatible with the Python library
[`shortuuid`](https://github.com/skorokithakis/shortuuid).

Often, one needs to use non-sequential IDs in places where users will see them,
but the IDs must be as concise and easy to use as possible. shortuuid solves
this problem by generating UUIDs using
[google/uuid](https://github.com/google/uuid) and then translating them to
base57 using lowercase and uppercase letters and digits, and removing
similar-looking characters such as l, 1, I, O and 0.

## Usage

```go
package main

import (
	"fmt"

	"github.com/lithammer/shortuuid/v4"
)

func main() {
	u := shortuuid.NewV4()
	fmt.Println(u) // KwSysDpxcBU9FNhGkn2dCf
}
```

To use UUID v5 (instead of the default v4), use `NewV5(namespace, name)`:

```go
shortuuid.NewV5(shortuuid.NameSpaceDNS, "example.com/")
shortuuid.NewV5(shortuuid.NameSpaceURL, "http://example.com")
shortuuid.NewV5(shortuuid.NameSpaceOID, "1.2.840.113549")
shortuuid.NewV5(shortuuid.NameSpaceX500, "CN=example,O=org")
```

<details>
<summary>Migrating from NewWithNamespace (deprecated)</summary>

`NewWithNamespace(name)` is deprecated but still available. It uses URL/DNS
heuristics based on the name prefix, supports `urn:oid:` and `x500:` prefixes,
and falls back to v4 when the name is empty.

```go
shortuuid.NewWithNamespace("http://example.com")
shortuuid.NewWithNamespace("urn:oid:1.2.840.113549")
shortuuid.NewWithNamespace("x500:CN=example,O=org")
```

See the [NewV5 (FromName) example](https://pkg.go.dev/github.com/lithammer/shortuuid/v4#example-NewV5-FromName)
for a migration guide.

</details>

It's possible to use a custom alphabet as well (at least 2
characters long).  
It will automatically sort and remove duplicates from your alphabet to ensure consistency

```go
alphabet := "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy="
shortuuid.NewV4WithAlphabet(alphabet) // iZsai==fWebXd5rLRWFB=u
```

For UUID v5 with a custom alphabet, provide an explicit namespace:

```go
shortuuid.NewV5WithAlphabet(shortuuid.NameSpaceDNS, "example.com/", alphabet)
```

Bring your own encoder! For example, base58 is popular among bitcoin.

```go
package main

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}

func main() {
	enc := base58Encoder{}
	fmt.Println(shortuuid.NewV4WithEncoder(enc)) // 6R7VqaQHbzC1xwA5UueGe6
}
```

## License

MIT
