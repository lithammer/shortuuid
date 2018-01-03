package shortuuid

import (
	"strings"

	uuid "github.com/satori/go.uuid"
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
	return DefaultEncoder.Encode(uuid.Must(uuid.NewV4()))
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
func NewWithEncoder(enc Encoder) string {
	return enc.Encode(uuid.Must(uuid.NewV4()))
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
func NewWithNamespace(name string) string {
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.Must(uuid.NewV4())
	case strings.HasPrefix(name, "http"):
		u = uuid.NewV5(uuid.NamespaceURL, name)
	default:
		u = uuid.NewV5(uuid.NamespaceDNS, name)
	}

	return DefaultEncoder.Encode(u)
}

// NewWithAlphabet returns a new UUIDv4, encoded with base57 using the
// alternative alphabet abc.
func NewWithAlphabet(abc string) string {
	enc := base57{newAlphabet(abc)}
	return enc.Encode(uuid.Must(uuid.NewV4()))
}
