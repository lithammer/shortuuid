package shortuuid

import (
	"crypto/sha1"
	"strings"
	"unsafe"

	"github.com/google/uuid"
)

// DefaultEncoder is the default encoder uses when generating new UUIDs, and is
// based on Base57.
var DefaultEncoder = &encoder{newAlphabet(DefaultAlphabet)}

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

// NewWithAlphabet returns a new UUIDv4, encoded with base57 using the
// alternative alphabet abc.
func NewWithAlphabet(abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(uuid.New())
}

// NewV7 returns a new UUIDv7, encoded with base57.
// UUIDv7 is time-based and provides monotonicity, making it ideal for certain database indexing scenarios.
func NewV7() string {
	return DefaultEncoder.Encode(uuid.Must(uuid.NewV7()))
}

// NewV7WithEncoder returns a new UUIDv7, encoded with enc.
func NewV7WithEncoder(enc Encoder) string {
	return enc.Encode(uuid.Must(uuid.NewV7()))
}

// NewV7WithNamespace returns a new UUIDv5 (or v7 if name is empty), encoded with base57.
func NewV7WithNamespace(name string) string {
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.Must(uuid.NewV7())
	case hasPrefixCaseInsensitive(name, "https://"):
		u = hashedUUID(uuid.NameSpaceURL, name)
	case hasPrefixCaseInsensitive(name, "http://"):
		u = hashedUUID(uuid.NameSpaceURL, name)
	default:
		u = hashedUUID(uuid.NameSpaceDNS, name)
	}

	return DefaultEncoder.Encode(u)
}

// NewV7WithAlphabet returns a new UUIDv7, encoded with base57 using the
// alternative alphabet abc.
func NewV7WithAlphabet(abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(uuid.Must(uuid.NewV7()))
}

func hasPrefixCaseInsensitive(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

func hashedUUID(space uuid.UUID, data string) (u uuid.UUID) {
	h := sha1.New()
	h.Write(space[:])                                         //nolint:errcheck
	h.Write(unsafe.Slice(unsafe.StringData(data), len(data))) //nolint:errcheck
	s := h.Sum(make([]byte, 0, sha1.Size))
	copy(u[:], s)
	u[6] = (u[6] & 0x0f) | uint8((5&0xf)<<4)
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 4122 variant
	return u
}
