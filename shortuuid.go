package shortuuid

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// UUID returns a new UUID v4.
func UUID() string {
	return New().UUID("")
}

// ShortUUID represents a short UUID encoder.
type ShortUUID struct {
	// alphabet is the character set to construct the UUID from.
	alphabet alphabet

	uuid string
}

// New returns a new (short) UUID with the default alphabet.
func New() *ShortUUID {
	return &ShortUUID{alphabet: newAlphabet(DefaultAlphabet)}
}

// NewWithAlphabet returns a new (short) UUID using alphabet.
// The alphabet will be sorted and if it contains duplicate characters, they
// will be removed.
func NewWithAlphabet(alphabet string) (*ShortUUID, error) {
	if alphabet == "" {
		return nil, fmt.Errorf("alphabet must not be empty")
	}

	return &ShortUUID{alphabet: newAlphabet(alphabet)}, nil
}

// UUID returns a new (short) UUID. If name is non-empty, the namespace
// matching the name will be used to generate a UUID.
func (su *ShortUUID) UUID(name string) string {
	var u uuid.UUID

	switch {
	case name == "":
		u = uuid.NewV4()
	case strings.HasPrefix(name, "http"):
		u = uuid.NewV5(uuid.NamespaceURL, name)
	default:
		u = uuid.NewV5(uuid.NamespaceDNS, name)
	}

	su.uuid = su.Encode(u)
	return su.uuid
}

// Encode encodes uuid.UUID into a string using the least significant bits
// (LSB) first according to the alphabet. if the most significant bits (MSB)
// are 0, the string might be shorter.
func (su *ShortUUID) Encode(u uuid.UUID) string {
	var num big.Int
	num.SetString(strings.Replace(u.String(), "-", "", 4), 16)

	// Calculate encoded length.
	factor := math.Log(float64(25)) / math.Log(float64(su.alphabet.Length()))
	length := math.Ceil(factor * float64(len(u.Bytes())))

	return su.numToString(&num, int(length))
}

// Decode decodes a string according to the alphabet into a uuid.UUID. If s is
// too short, its most significant bits (MSB) will be padded with 0 (zero).
func (su *ShortUUID) Decode(u string) (uuid.UUID, error) {
	return uuid.FromString(su.stringToNum(u))
}

// String returns a (short) UUID v4, and will generate one if necessary. To
// generate a new string, use UUID().
// Implements the fmt.Stringer interface.
func (su *ShortUUID) String() string {
	if su.uuid != "" {
		return su.uuid
	}

	return su.UUID("")
}

func (su *ShortUUID) numToString(number *big.Int, padToLen int) string {
	var (
		out   string
		digit *big.Int
	)

	for number.Uint64() > 0 {
		number, digit = new(big.Int).DivMod(number, big.NewInt(su.alphabet.Length()), new(big.Int))
		out += su.alphabet.chars[digit.Int64()]
	}

	if padToLen > 0 {
		remainder := math.Max(float64(padToLen-len(out)), 0)
		out = out + strings.Repeat(su.alphabet.chars[0], int(remainder))
	}

	return out
}

// stringToNum converts a string a number using the given alpabet.
func (su *ShortUUID) stringToNum(s string) string {
	n := big.NewInt(0)

	for i := len(s) - 1; i >= 0; i-- {
		n.Mul(n, big.NewInt(su.alphabet.Length()))
		n.Add(n, big.NewInt(su.alphabet.Index(string(s[i]))))
	}

	x := fmt.Sprintf("%x", n)

	// Pad the most significant bit (MSG) with 0 (zero) if the string is too short.
	if len(x) < 32 {
		x = strings.Repeat("0", 32-len(x)) + x
	}

	return fmt.Sprintf("%s-%s-%s-%s-%s", x[0:8], x[8:12], x[12:16], x[16:20], x[20:32])
}
