package shortuuid

import (
	"github.com/google/uuid"
)

// DefaultEncoder is the default encoder uses when generating new UUIDs, and is
// based on Base57.
var DefaultEncoder = &base57{newAlphabet(DefaultAlphabet)}

// Encoder is an interface for encoding/decoding UUIDs to strings.
type Encoder interface {
	Encode(uuid.UUID) string
	Decode(string) (uuid.UUID, error)
}

// New returns a new UUIDv4, encoded with base57.
func New() string {
	rv, err := NewTyped(UUID_v4)
	if err != nil {
		panic(err)
	}
	return rv
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
func NewWithEncoder(enc Encoder) string {
	rv, err := NewTypedWithEncoder(UUID_v4, enc)
	if err != nil {
		panic(err)
	}
	return rv
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
func NewWithNamespace(name string) string {
	rv, err := NewTypedWithNamespace(UUID_v4, UUID_v5, name)
	if err != nil {
		panic(err)
	}

	return rv
}

// NewWithAlphabet returns a new UUIDv4, encoded with base57 using the
// alternative alphabet abc.
func NewWithAlphabet(abc string) string {
	enc := base57{newAlphabet(abc)}
	rv, err := NewTypedWithAlphabet(UUID_v4, abc, enc)
	if err != nil {
		panic(err)
	}

	return rv
}
