package state

import (
	"strconv"
)

// ReleaseKind is a enum (iota) type for release states.
type ReleaseKind int

const (
	ReleaseKindDev ReleaseKind = iota + 1
	ReleaseKindAlpha
	ReleaseKindBeta
	ReleaseKindRC
	ReleaseKindGA
	ReleaseKindEOL
	ReleaseKindUnsupported
)

// NewReleaseKind returns a new release kind from the given type.
// If the given input is invalid, it returns error.
func NewReleaseKind(s int) (ReleaseKind, error) {
	k := ReleaseKind(s)
	if !k.IsValid() {
		return 0, ErrReleaseKindInvalid
	}
	return k, nil
}

// NewReleaseKindFromString returns a new release kind from the given type.
// If the given input is invalid, it returns error.
func NewReleaseKindFromString(s string) (ReleaseKind, error) {
	k := ReleaseKind(0)
	switch s {
	case "dev":
		k = ReleaseKindDev
	case "alpha":
		k = ReleaseKindAlpha
	case "beta":
		k = ReleaseKindBeta
	case "rc":
		k = ReleaseKindRC
	case "ga":
		k = ReleaseKindGA
	case "eol":
		k = ReleaseKindEOL
	case "unsupported":
		k = ReleaseKindUnsupported
	default:
		return 0, ErrReleaseKindInvalid
	}
	if !k.IsValid() {
		return 0, ErrReleaseKindInvalid
	}
	return k, nil
}

// IsValid returns true if the release kind is valid.
func (k ReleaseKind) IsValid() bool {
	switch k {
	case
		ReleaseKindDev,
		ReleaseKindAlpha,
		ReleaseKindBeta,
		ReleaseKindRC,
		ReleaseKindGA,
		ReleaseKindEOL,
		ReleaseKindUnsupported:
		return true
	default:
		return false
	}
}

// Next returns the next release kind.
// If the next release kind is invalid, it returns error.
func (k ReleaseKind) Next() (ReleaseKind, error) {
	result := k + 1
	if result.IsValid() {
		return result, nil
	}
	return 0, ErrReleaseKindInvalid
}

// Prev returns the previous release kind.
// If the previous release kind is invalid, it returns error.
func (k ReleaseKind) Prev() (ReleaseKind, error) {
	result := k - 1
	if result.IsValid() {
		return result, nil
	}
	return 0, ErrReleaseKindInvalid
}

func (k ReleaseKind) String() string {
	switch k {
	case ReleaseKindDev:
		return "dev"
	case ReleaseKindAlpha:
		return "alpha"
	case ReleaseKindBeta:
		return "beta"
	case ReleaseKindRC:
		return "rc"
	case ReleaseKindGA:
		return "ga"
	case ReleaseKindEOL:
		return "eol"
	case ReleaseKindUnsupported:
		return "unsupported"
	default:
		panic(ErrReleaseKindInvalid)
	}
}

// Is returns true if the release kind is equal to the given release kind.
func (k ReleaseKind) Is(s ReleaseKind) bool {
	return k == s
}

// MarshalJSON implements the json.Marshaler interface.
func (k *ReleaseKind) MarshalJSON() ([]byte, error) {
	if !k.IsValid() {
		return nil, ErrReleaseKindInvalid
	}
	return []byte(strconv.Quote(k.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (k *ReleaseKind) UnmarshalJSON(b []byte) error {
	val, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	n, err := NewReleaseKindFromString(val)
	if err != nil {
		return err
	}
	*k = n
	return nil
}
