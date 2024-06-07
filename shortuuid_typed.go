package shortuuid

import (
	"fmt"
	"github.com/google/uuid"
)

type UnderlyingType int

const (
	UUID_v1 UnderlyingType = 0
	UUID_v3 UnderlyingType = iota
	UUID_v4 UnderlyingType = iota
	UUID_v5 UnderlyingType = iota
	UUID_v6 UnderlyingType = iota
	UUID_v7 UnderlyingType = iota
)

func ut2uuid(ut UnderlyingType) (uuid.UUID, error) {
	switch ut {
	case UUID_v1:
		return uuid.NewUUID()
	case UUID_v4:
		return uuid.New(), nil
	case UUID_v6:
		return uuid.NewV6()
	case UUID_v7:
		return uuid.NewV7()
	default:
		panic("unknown underlying type")
	}
}

// NewTyped returns a new id (based on ut type), encoded with DefaultEncoder.
func NewTyped(ut UnderlyingType) (string, error) {
	rv, err := ut2uuid(ut)
	if err != nil {
		return "", err
	}

	return DefaultEncoder.Encode(rv), nil
}

// NewTypedWithEncoder returns a new id (based on ut type), encoded with enc.
func NewTypedWithEncoder(ut UnderlyingType, enc Encoder) (string, error) {
	rv, err := ut2uuid(ut)
	if err != nil {
		return "", err
	}

	return enc.Encode(rv), nil
}

// NewTypedWithNamespace returns a new id (based on ut type and name)
//
// when name is empty id will be based on emptyNameUt,
// otherwise id  will be based on ut
func NewTypedWithNamespace(emptyNameUt UnderlyingType, ut UnderlyingType, name string) (string, error) {
	nameLen := len(name)

	if nameLen == 0 {
		rv, err := ut2uuid(emptyNameUt)
		if err != nil {
			return "", err
		}

		return DefaultEncoder.Encode(rv), nil
	}

	ns := (func(name string, nameLen int) uuid.UUID {
		//returns namespace by name prefix (case-insensitive compare)
		var ch uint8
		if nameLen >= len("http://") {
			for {
				idx := 0
				ch = name[idx]
				if ch != 'h' && ch != 'H' {
					break
				}

				idx++
				ch = name[idx]
				if ch != 't' && ch != 'T' {
					break
				}

				idx++
				ch = name[idx]
				if ch != 't' && ch != 'T' {
					break
				}

				idx++
				ch = name[idx]
				if ch != 'p' && ch != 'P' {
					break
				}

				idx++
				ch = name[idx]
				if ch != ':' {
					//maybe httpS ?
					if nameLen >= len("https://") && (ch == 's' || ch == 'S') {
						idx++
						ch = name[idx]
						if ch != ':' {
							break
						}
					} else {
						break
					}
				}

				idx++
				ch = name[idx]
				if ch != '/' {
					break
				}

				idx++
				ch = name[idx]
				if ch != '/' {
					break
				}

				// yeah! name starts with "http://" or "https://"
				return uuid.NameSpaceURL
			}
		}

		// default is DNS (backward compatibility)
		return uuid.NameSpaceDNS
	})(name, nameLen)

	switch ut {
	case UUID_v5:
		return DefaultEncoder.Encode(uuid.NewSHA1(ns, []byte(name))), nil
	case UUID_v3:
		return DefaultEncoder.Encode(uuid.NewMD5(ns, []byte(name))), nil
	default:
		return "", fmt.Errorf("unsupported underlying type [%v] for non-empty name", ut)
	}
}

// NewTypedWithAlphabet returns a new id, encoded with enc using the
// alternative alphabet abc.
func NewTypedWithAlphabet(ut UnderlyingType, abc string, enc Encoder) (string, error) {
	rv, err := ut2uuid(ut)
	if err != nil {
		return "", err
	}

	return enc.Encode(rv), nil
}
