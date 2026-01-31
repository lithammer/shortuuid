// Package shortuuid provides a library for generating concise, unambiguous,
// URL-safe UUIDs. It generates UUIDs using github.com/google/uuid and then
// translates them to base57 using a custom alphabet that removes similar-looking
// characters (l, 1, I, O, 0).
//
// The package is compatible with the Python library shortuuid and provides
// both a default encoder (base57) and support for custom alphabets and encoders.
package shortuuid

import (
	"crypto/sha1"
	"strings"
	"unsafe"

	"github.com/google/uuid"
)

// DefaultEncoder is the default encoder used when generating new UUIDs, and is
// based on Base57.
var DefaultEncoder = b57Encoder{}

// Encoder is an interface for encoding/decoding UUIDs to strings.
type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

// New returns a new UUIDv4, encoded with base57.
func New() string {
	return DefaultEncoder.Encode(uuid.New())
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
func NewWithEncoder(enc Encoder) string {
	return enc.Encode(uuid.New())
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
func NewWithNamespace(name string) string {
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.New()
	case hasPrefixCaseInsensitive(name, "https://"):
		u = hashedUUID(uuid.NameSpaceURL, name)
	case hasPrefixCaseInsensitive(name, "http://"):
		u = hashedUUID(uuid.NameSpaceURL, name)
	default:
		u = hashedUUID(uuid.NameSpaceDNS, name)
	}

	return DefaultEncoder.Encode(u)
}

// NewWithAlphabet returns a new UUIDv4, encoded using the alternative
// alphabet abc.
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
// The alphabet will be automatically sorted and deduplicated to ensure
// consistency.
func NewWithAlphabet(abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(uuid.New())
}

func hasPrefixCaseInsensitive(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

func hashedUUID(space uuid.UUID, data string) (u uuid.UUID) {
	h := sha1.New()
	h.Write(space[:])
	h.Write(unsafe.Slice(unsafe.StringData(data), len(data)))
	s := h.Sum(make([]byte, 0, sha1.Size))
	copy(u[:], s)
	u[6] = (u[6] & 0x0f) | uint8((5&0xf)<<4)
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 4122 variant
	return u
}
