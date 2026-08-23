// Package shortuuid provides a library for generating concise, unambiguous,
// URL-safe UUIDs. It generates UUIDs using the standard library uuid package
// and then translates them to base57 using a custom alphabet that removes
// similar-looking characters (l, 1, I, O, 0).
//
// The package is compatible with the Python library shortuuid and provides
// both a default encoder (base57) and support for custom alphabets and encoders.
package shortuuid

import (
	"crypto/sha1"
	"strings"
	"unsafe"
	"uuid"
)

// DefaultEncoder is the default encoder used when generating new UUIDs, and is
// based on base57.
var DefaultEncoder = b57Encoder{}

// Well-known namespace IDs from RFC 9562, Section 6.6. The standard library
// uuid package does not export these.
var (
	NameSpaceDNS  = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NameSpaceURL  = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NameSpaceOID  = uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NameSpaceX500 = uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// Encoder is an interface for encoding/decoding UUIDs to strings.
//
// Decode must invert Encode. The encoders this package provides write the
// UUID as a base-N number, most significant digit first, and accept shorter
// input on Decode by treating the missing leading digits as zeros.
type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

// NewEncoder returns an Encoder over the alphabet abc, which it sorts and
// deduplicates first, so the same characters in any order encode alike.
//
// Build one encoder and reuse it; it encodes a UUID of any version:
//
//	enc.Encode(uuid.NewV7())
//	enc.Encode(shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com"))
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
func NewEncoder(abc string) Encoder {
	a := newAlphabet(abc)
	if a.isDefault() {
		return DefaultEncoder
	}
	return encoder{a}
}

// New returns a new UUID, encoded with base57.
//
// Callers with no need for a particular UUID version should use New. It is
// currently equivalent to NewV4.
func New() string {
	return DefaultEncoder.Encode(uuid.New())
}

// NewV4 returns a new UUIDv4, encoded with base57. Version 4 UUIDs are 122
// bits of random data.
func NewV4() string {
	return DefaultEncoder.Encode(uuid.NewV4())
}

// NewV7 returns a new UUIDv7, encoded with base57. Version 7 UUIDs lead with a
// Unix millisecond timestamp, so they sort in creation order.
//
// The encoded strings sort in that order too, because every encoding is 22
// characters of most-significant-digit-first base57 over a sorted alphabet.
func NewV7() string {
	return DefaultEncoder.Encode(uuid.NewV7())
}

// NewV5 returns the UUIDv5 of name within namespace, encoded with base57.
// The same namespace and name always produce the same shortuuid.
func NewV5(namespace uuid.UUID, name string) string {
	return DefaultEncoder.Encode(UUIDv5(namespace, name))
}

// UUIDv5 returns the version 5 (SHA-1, name-based) UUID of name within
// namespace, as defined in RFC 9562, Section 5.5.
//
// Pass the bare name: an OID is "1.2.840.113549", not
// "urn:oid:1.2.840.113549".
func UUIDv5(namespace uuid.UUID, name string) (u uuid.UUID) {
	h := sha1.New()
	h.Write(namespace[:])
	// Writing name in place skips the []byte copy (sha1 has no WriteString);
	// safe because the hash neither mutates nor retains the slice.
	h.Write(unsafe.Slice(unsafe.StringData(name), len(name)))
	s := h.Sum(make([]byte, 0, sha1.Size))
	copy(u[:], s)
	u[6] = (u[6] & 0x0f) | 0x50 // version 5
	u[8] = (u[8] & 0x3f) | 0x80 // variant 10
	return u
}

// NewWithEncoder returns a new UUID, encoded with enc. The version is the one
// New returns.
//
// Deprecated: Call enc.Encode directly. It costs the same, and it encodes the
// other versions too:
//
//	enc.Encode(uuid.New())
//	enc.Encode(uuid.NewV7())
func NewWithEncoder(enc Encoder) string {
	return enc.Encode(uuid.New())
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
//
// The namespace is guessed from the name: NameSpaceURL for http and https
// prefixes, matched case-insensitively, and NameSpaceDNS for everything else.
// This mirrors the Python shortuuid library, so both produce the same ID for
// the same name. Names that are neither hostnames nor URLs are hashed under
// NameSpaceDNS too, OIDs and X.500 DNs included; reach those namespaces with
// NewV5 instead.
//
// Deprecated: Use NewV5 with an explicit namespace, or New for a random v4.
func NewWithNamespace(name string) string {
	switch {
	case name == "":
		return New()
	case hasPrefixFold(name, "https://"), hasPrefixFold(name, "http://"):
		return NewV5(NameSpaceURL, name)
	default:
		return NewV5(NameSpaceDNS, name)
	}
}

// NewWithAlphabet returns a new UUID, encoded using the alternative alphabet
// abc. The version is the one New returns.
//
// The alphabet is sorted and deduplicated first, so the same characters in
// any order encode alike. Panics if fewer than 2 distinct characters remain.
//
// Deprecated: That preparation costs more than the encoding and happens on
// every call, with no way to lift it out. Keep one encoder instead:
//
//	var enc = shortuuid.NewEncoder(abc)
//	enc.Encode(uuid.New())
func NewWithAlphabet(abc string) string {
	return NewEncoder(abc).Encode(uuid.New())
}

func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}
