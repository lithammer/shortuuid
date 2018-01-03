package shortuuid

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"fmt"
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
	str, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Unable to create UUIDv4: %s", err))
	}
	return DefaultEncoder.Encode(str)
}

// NewWithEncoder returns a new UUIDv4, encoded with enc.
func NewWithEncoder(enc Encoder) string {
	str, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Unable to create UUIDv4: %s", err))
	}
	return enc.Encode(str)
}

// NewWithNamespace returns a new UUIDv5 (or v4 if name is empty), encoded with base57.
func NewWithNamespace(name string) string {
	var u uuid.UUID
	var err error

	switch {
	case name == "":
		u, err = uuid.NewV4()
		if err != nil {
			panic(fmt.Sprintf("Unable to create UUIDv4: %s", err))
		}
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
	str, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Unable to create UUIDv4: %s", err))
	}
	return enc.Encode(str)
}
