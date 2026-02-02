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

var (
	// DefaultEncoder is the default encoder used when generating new UUIDs, and is
	// based on Base57.
	DefaultEncoder = b57Encoder{}

	// NameSpaceDNS is the UUID DNS namespace.
	NameSpaceDNS = uuid.NameSpaceDNS

	// NameSpaceURL is the UUID URL namespace.
	NameSpaceURL = uuid.NameSpaceURL

	// NameSpaceOID is the UUID OID namespace.
	NameSpaceOID = uuid.NameSpaceOID

	// NameSpaceX500 is the UUID X500 namespace.
	NameSpaceX500 = uuid.NameSpaceX500
)

// Encoder is an interface for encoding/decoding UUIDs to strings.
type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

// NewV4 returns a new UUIDv4, encoded with base57.
func NewV4() string {
	return DefaultEncoder.Encode(uuid.New())
}

// NewV4WithEncoder returns a new UUIDv4, encoded with enc.
func NewV4WithEncoder(enc Encoder) string {
	return enc.Encode(uuid.New())
}

// NewV4WithAlphabet returns a new UUIDv4, encoded using the alternative
// alphabet abc.
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
// The alphabet will be automatically sorted and deduplicated to ensure
// consistency.
func NewV4WithAlphabet(abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(uuid.New())
}

// NewV5 returns a new UUIDv5, encoded with base57.
// The provided namespace is used directly, with no heuristic.
func NewV5(namespace uuid.UUID, name string) string {
	return DefaultEncoder.Encode(hashedUUID(namespace, name))
}

// NewV5WithAlphabet returns a new UUIDv5, encoded using the alternative
// alphabet abc. The provided namespace is used directly, with no heuristic.
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
// The alphabet will be automatically sorted and deduplicated to ensure
// consistency.
func NewV5WithAlphabet(namespace uuid.UUID, name, abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(hashedUUID(namespace, name))
}

// New returns a new UUIDv4, encoded with base57.
//
// Deprecated: Use NewV4.
func New() string {
	return NewV4()
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
//
// Deprecated: Use NewV4WithEncoder.
func NewWithEncoder(enc Encoder) string {
	return NewV4WithEncoder(enc)
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
//
// Deprecated: Use NewV5(namespace, name) for explicit namespaces or NewV4 for
// random v4. This function keeps URL/DNS heuristics and the OID/X500 prefixes.
func NewWithNamespace(name string) string {
	if name == "" {
		return NewV4()
	}
	return DefaultEncoder.Encode(uuidV5FromName(name))
}

// NewWithAlphabet returns a new UUIDv4, encoded using the alternative
// alphabet abc.
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
// The alphabet will be automatically sorted and deduplicated to ensure
// consistency.
//
// Deprecated: Use NewV4WithAlphabet.
func NewWithAlphabet(abc string) string {
	return NewV4WithAlphabet(abc)
}

func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

func cutPrefixFold(s, prefix string) (string, bool) {
	if hasPrefixFold(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}

func uuidV5FromName(name string) uuid.UUID {
	if hasPrefixFold(name, "https://") || hasPrefixFold(name, "http://") {
		return hashedUUID(uuid.NameSpaceURL, name)
	}

	if after, found := cutPrefixFold(name, "urn:oid:"); found {
		return hashedUUID(uuid.NameSpaceOID, after)
	}

	if after, found := cutPrefixFold(name, "x500:"); found {
		return hashedUUID(uuid.NameSpaceX500, after)
	}

	return hashedUUID(uuid.NameSpaceDNS, name)
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
