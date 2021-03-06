package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rotisserie/eris"
)

var (
	ErrNoRelease = eris.New("state: no releases")
)

// State holds all the release operations.
type State struct {
	Releases []*Release `json:"releases,omitempty"`
}

// NewState returns a new and empty state.
func NewState() *State {
	return &State{
		Releases: make([]*Release, 0),
	}
}

// CreateRelease creates a new release from the given data.
// It prepends the release to the state.
func (s *State) CreateRelease(r *Release) error {
	if r == nil {
		return eris.New("state: release is nil")
	}
	if err := r.Validate(); err != nil {
		return err
	}
	r.CreatedAt = time.Now()
	{
		if len(s.Releases) != 0 {
			r.PreviousBlockHash = s.Releases[0].BlockHash
		}
		hash, err := r.Hash()
		if err != nil {
			return err
		}
		r.BlockHash = hash
	}
	s.Releases = append([]*Release{r}, s.Releases...)
	return nil
}

// Clean removes all the releases of given kind from the state.
func (s *State) Clean(kind ReleaseKind) {
	for i, v := range s.Releases {
		if v.Kind.Is(kind) {
			s.Releases = append(s.Releases[:i], s.Releases[i+1:]...)
		}
	}
}

// Promote promotes the latest release of the given kind to the next kind.
func (s *State) Promote(from ReleaseKind) error {
	f := s.Latest(from)
	t, err := f.Promote()
	if err != nil {
		return eris.Wrap(err, "cli: could not promote")
	}
	if err := s.CreateRelease(t); err != nil {
		return eris.Wrap(err, "cli: could not promote")
	}
	return nil
}

// PromoteTo promotes the latest release of the previous kind to the given kind.
func (s *State) PromoteTo(to ReleaseKind) error {
	from, err := to.Prev()
	if err != nil {
		return err
	}
	return s.Promote(from)
}

// Export exports the state to the given filepath.
func (s *State) Export(filepath string) error {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, b, os.ModePerm)
}

// Import imports the state from the given filepath.
func (s *State) Import(filepath string) error {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if err := s.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates the state.
// It checks for block hashes and matches the previous block hashes.
func (s *State) Validate() error {
	var previousBlock Hash
	for i, v := range s.Releases {
		if err := v.Validate(); err != nil {
			return err
		}
		{
			hash, err := v.Hash()
			if err != nil {
				return err
			}
			if !v.BlockHash.Match(hash) {
				return eris.Errorf("state: block %s is corrupted. calculated hash %s", v.BlockHash.Short(), hash.String())
			}
		}
		if !previousBlock.IsEmpty() && !previousBlock.Match(v.BlockHash) {
			return eris.New("state: release hash does not match with previous release")
		}
		previousBlock = v.PreviousBlockHash
		// only last block can have a nil previous block hash
		if previousBlock.IsEmpty() && i != len(s.Releases)-1 {
			return eris.New("state: missing previous block hash")
		}
	}
	return nil
}

// Rollback removes the latest release from the state.
// It does not care about release kind.
// It just pops the latest release from the state stack.
func (s *State) Rollback() {
	if len(s.Releases) == 0 {
		return
	}
	s.Releases = s.Releases[1:]
}

// Head returns the latest release from the state stack.
// It does not care about release kind.
// It returns a shallow copy of the release.
func (s *State) Head() (*Release, error) {
	if len(s.Releases) == 0 {
		return nil, ErrNoRelease
	}
	return s.Releases[0].Copy(), nil
}

// Latest returns the latest release of the given kind.
func (s *State) Latest(kind ReleaseKind) *Release {
	blank := &Release{
		Kind: kind,
		Tag:  "v0.0.0",
	}
	for _, v := range s.Releases {
		if v.Kind.Is(kind) {
			return v.Copy()
		}
	}
	return blank
}

// GetRelease returns a shallow copy of the release of the given hash.
func (s *State) GetRelease(hash Hash) (*Release, error) {
	for _, v := range s.Releases {
		if v.BlockHash.String() == hash.String() || v.BlockHash.Short() == hash.Short() {
			return v.Copy(), nil
		}
	}
	return nil, eris.Errorf("state: release with hash %s not found", hash.String())
}
