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
// based on Base57.
var DefaultEncoder = b57Encoder{}

// Well-known namespace IDs from RFC 9562, Section 6.6. The standard library
// uuid package does not export these, so they are spelled out here.
var (
	// NameSpaceDNS is the UUID DNS namespace.
	NameSpaceDNS = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	// NameSpaceURL is the UUID URL namespace.
	NameSpaceURL = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")

	// NameSpaceOID is the UUID OID namespace.
	NameSpaceOID = uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")

	// NameSpaceX500 is the UUID X500 namespace.
	NameSpaceX500 = uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// Encoder is an interface for encoding/decoding UUIDs to strings.
type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

// NewEncoder returns an Encoder over the alphabet abc. The alphabet is sorted
// and deduplicated first, so any permutation of the same characters yields the
// same encoding.
//
// Combined with a UUID of the caller's choosing, this covers every
// version/alphabet pairing without a constructor for each one:
//
//	enc := shortuuid.NewEncoder(abc)
//	enc.Encode(uuid.New())
//	enc.Encode(uuid.NewV7())
//	enc.Encode(shortuuid.UUIDv5(shortuuid.NameSpaceDNS, "example.com"))
//
// Panics if abc (after removing duplicates) has fewer than 2 characters.
func NewEncoder(abc string) Encoder {
	a := newAlphabet(abc)
	if string(a.chars) == DefaultAlphabet {
		return DefaultEncoder
	}
	return encoder{a}
}

// New returns a new UUIDv4, encoded with base57.
func New() string {
	return DefaultEncoder.Encode(uuid.New())
}

// NewV5 returns the UUIDv5 of name within namespace, encoded with base57.
//
// The namespace is given explicitly, unlike NewWithNamespace which guesses it
// from the name. Use NewEncoder to encode with a different alphabet:
//
//	shortuuid.NewEncoder(abc).Encode(shortuuid.UUIDv5(ns, name))
func NewV5(namespace uuid.UUID, name string) string {
	return DefaultEncoder.Encode(UUIDv5(namespace, name))
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
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
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.New()
	case hasPrefixCaseInsensitive(name, "https://"):
		u = UUIDv5(NameSpaceURL, name)
	case hasPrefixCaseInsensitive(name, "http://"):
		u = UUIDv5(NameSpaceURL, name)
	default:
		u = UUIDv5(NameSpaceDNS, name)
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
	return NewWithEncoder(NewEncoder(abc))
}

func hasPrefixCaseInsensitive(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// UUIDv5 returns the version 5 (SHA-1, name-based) UUID of name within
// namespace, as defined in RFC 9562, Section 5.5.
//
// Pass the bare name: an OID is "1.2.840.113549", not
// "urn:oid:1.2.840.113549".
func UUIDv5(namespace uuid.UUID, name string) (u uuid.UUID) {
	h := sha1.New()
	h.Write(namespace[:])
	h.Write(unsafe.Slice(unsafe.StringData(name), len(name)))
	s := h.Sum(make([]byte, 0, sha1.Size))
	copy(u[:], s)
	u[6] = (u[6] & 0x0f) | uint8((5&0xf)<<4)
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 4122 variant
	return u
}
