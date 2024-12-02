package shortuuid

import (
	"strings"

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
		u = uuid.NewSHA1(uuid.NameSpaceURL, []byte(name))
	case hasPrefixCaseInsensitive(name, "http://"):
		u = uuid.NewSHA1(uuid.NameSpaceURL, []byte(name))
	default:
		u = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))
	}

	return DefaultEncoder.Encode(u)
}

// NewWithAlphabet returns a new UUIDv4, encoded with base57 using the
// alternative alphabet abc.
func NewWithAlphabet(abc string) string {
	enc := encoder{newAlphabet(abc)}
	return enc.Encode(uuid.New())
}

func hasPrefixCaseInsensitive(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}
