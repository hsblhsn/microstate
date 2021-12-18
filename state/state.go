package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rotisserie/eris"
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
	s.Releases = append([]*Release{r}, s.Releases...)
	return nil
}

// Latest returns the latest release of the given kind.
func (s *State) Latest(kind ReleaseKind) *Release {
	blank := &Release{
		Kind: kind,
		Tag:  "v0.0.0",
	}
	for _, v := range s.Releases {
		if v.Kind.Is(kind) {
			return v
		}
	}
	return blank
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
	return json.Unmarshal(b, &s)
}
