# shortuuid

[![Build Status](https://img.shields.io/travis/renstrom/shortuuid.svg?style=flat-square)](https://travis-ci.org/renstrom/shortuuid)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/renstrom/shortuuid)

A Go library that generates concise, unambiguous, URL-safe UUIDs. Based on and compatible with the Python library [`shortuuid`](https://github.com/stochastic-technologies/shortuuid).

Often, one needs to use non-sequential IDs in places where users will see them, but the IDs must be as concise and easy to use as possible. shortuuid solves this problem by generating UUIDs using [satori/go.uuid](https://github.com/satori/go.uuid) and then translating them to base57 using lowercase and uppercase letters and digits, and removing similar-looking characters such as l, 1, I, O and 0.

## Usage

```go
package main

import (
    "fmt"

    "github.com/renstrom/shortuuid"
)

func main() {
    id := shortuuid.UUID()  // "ajLWxEodc6CmQLHADuKVwD"

    u := shortuuid.New()
    fmt.Printf("%s", u)     // Cekw67uyMpBGZLRP2HFVbe
    u.UUID("")              // Generate a new UUID
    fmt.Printf("%s", u)     // 4pUYNRFHTG3YVgThPZvCgC
}
```

To use UUID v5 (instead of the default v4), pass a namespace (DNS or URL) to the `.UUID(name string)` call:

```go
shortuuid.New().UUID("http://example.com")
```

## License

MIT
