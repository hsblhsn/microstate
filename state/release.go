package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rotisserie/eris"
)

var (
	ErrReleaseTagInvalid   = eris.New("state: release tag is invalid")
	ErrReleaseKindInvalid  = eris.New("state: release kind is invalid")
	ErrReleaseKindIsNotDev = eris.New("state: to create a release, kind must have to be dev")
	ErrServiceMapInvalid   = eris.New("state: invalid service map. at least one active service is required to create a release")
)

type Hash string

func (h Hash) String() string {
	return string(h)
}

func (h Hash) Short() string {
	if len(h) < 10 {
		return ""
	}
	return string(h[:9])
}

func (h Hash) Match(s Hash) bool {
	return h == s
}

func (h Hash) IsEmpty() bool {
	return h == ""
}

// Release type hols release informtation.
type Release struct {
	Kind              ReleaseKind `json:"kind,omitempty"`
	Tag               string      `json:"tag,omitempty"`
	Versions          VersionMap  `json:"versions,omitempty"`
	CreatedAt         time.Time   `json:"created_at,omitempty"`
	BlockHash         Hash        `json:"block_hash,omitempty"`
	PreviousBlockHash Hash        `json:"previous_block_hash,omitempty"`
}

// NewRelease returns a new release from the given data.
func NewRelease(k ReleaseKind, tag string, v VersionMap) (*Release, error) {
	if !k.Is(ReleaseKindDev) {
		return nil, ErrReleaseKindIsNotDev
	}
	return &Release{
		Kind:     k,
		Tag:      tag,
		Versions: v,
	}, nil
}

// Validate returns an error if the release is invalid.
func (r Release) Validate() error {
	if !r.Kind.IsValid() {
		return ErrReleaseKindInvalid
	}
	v, err := semver.NewVersion(r.Tag)
	if err != nil {
		return eris.Wrapf(
			ErrReleaseTagInvalid,
			"state: could not parse version string %q from release tag: %v", r.Tag, err,
		)
	}
	if version := fmt.Sprintf("v%s", v); r.Tag != version {
		return eris.Wrapf(
			ErrReleaseTagInvalid,
			"state: release tag %q does not match to the parsed version %q", r.Tag, version,
		)
	}
	if len(r.Versions) == 0 {
		return ErrServiceMapInvalid
	}
	return nil
}

// String returns the string representation of the release.
func (r Release) String() string {
	return fmt.Sprintf("%s@%s", r.Kind, r.Tag)
}

// Copy returns a copy of the release.
// The copied release is safe to modify.
func (r Release) Copy() *Release {
	r.Versions = r.Versions.Copy()
	return &r
}

// Promote the release to the next release kind.
// It returns error if next kind is invalid.
// It returns the promoted copy of the release.
func (r Release) Promote() (*Release, error) {
	copied := r.Copy()
	next, err := copied.Kind.Next()
	if err != nil {
		return nil, err
	}
	copied.Kind = next
	version, err := semver.NewVersion(copied.Tag)
	if err != nil {
		return nil, err
	}
	kind := copied.Kind.String()
	if copied.Kind.Is(ReleaseKindGA) {
		kind = ""
	}
	{
		upgradedVersion, err := version.SetPrerelease(kind)
		if err != nil {
			return nil, err
		}
		upgradedVersion, err = upgradedVersion.SetMetadata("")
		if err != nil {
			return nil, err
		}
		version = &upgradedVersion
	}
	copied.Tag = fmt.Sprintf("v%s", version)
	return copied, nil
}

func (r Release) Hash() (Hash, error) {
	r.BlockHash = ""
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(b)
	return Hash(hex.EncodeToString(h.Sum(nil))), nil
}
