package state_test

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/hsblhsn/microstate/state"
	"github.com/rotisserie/eris"
)

func TestRelease_Validate(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Release", func() {
		g.It("should fail on any other kind than dev", func() {
			r, err := state.NewRelease(state.ReleaseKindBeta, "v1.0.0", state.VersionMap{})
			g.Assert(err).Equal(state.ErrReleaseKindIsNotDev)
			g.Assert(r).IsNil()
		})
		g.It("should validate", func() {
			r, err := state.NewRelease(state.ReleaseKindDev, "v1.0.0", state.VersionMap{})
			g.Assert(err).IsNil()
			g.Assert(r.Validate()).IsNotNil()
		})
		g.It("should fail to validate on unknown kind", func() {
			r, err := state.NewRelease(state.ReleaseKindDev, "v1.0.0", state.VersionMap{})
			r.Kind = 99
			g.Assert(err).IsNil()
			g.Assert(r.Validate()).Equal(state.ErrReleaseKindInvalid)
		})
		g.It("should fail to validate on empty tag", func() {
			r, err := state.NewRelease(state.ReleaseKindDev, "", state.VersionMap{})
			g.Assert(err).IsNil()
			g.Assert(eris.Cause(r.Validate())).Equal(state.ErrReleaseTagInvalid)
		})
		g.It("should fail to validate on empty service map", func() {
			r, err := state.NewRelease(state.ReleaseKindDev, "v1.0.0", state.VersionMap{})
			g.Assert(err).IsNil()
			r.Versions = state.VersionMap{}
			g.Assert(r.Validate()).Equal(state.ErrServiceMapInvalid)
		})
	})
}

func TestRelease_Upgrade(t *testing.T) {
	g := goblin.Goblin(t)
	r, err := state.NewRelease(state.ReleaseKindDev, "v1.0.0-dev.2+feature-a", state.VersionMap{})
	g.Assert(err).IsNil()
	g.Assert(r.String()).Equal("dev@v1.0.0-dev.2+feature-a")
	g.Describe("Promote", func() {
		g.It("should get promoted to alpha", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindAlpha)
			g.Assert(upgraded.Tag).Equal("v1.0.0-alpha")
			g.Assert(upgraded.String()).Equal("alpha@v1.0.0-alpha")
			r = upgraded
		})
		g.It("should get promoted to beta", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindBeta)
			g.Assert(upgraded.Tag).Equal("v1.0.0-beta")
			g.Assert(upgraded.String()).Equal("beta@v1.0.0-beta")
			r = upgraded
		})
		g.It("should get promoted to rc", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindRC)
			g.Assert(upgraded.Tag).Equal("v1.0.0-rc")
			g.Assert(upgraded.String()).Equal("rc@v1.0.0-rc")
			r = upgraded
		})
		g.It("should get promoted to ga", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindGA)
			g.Assert(upgraded.Tag).Equal("v1.0.0")
			g.Assert(upgraded.String()).Equal("ga@v1.0.0")
			r = upgraded
		})
		g.It("should get promoted to eol", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindEOL)
			g.Assert(upgraded.Tag).Equal("v1.0.0-eol")
			g.Assert(upgraded.String()).Equal("eol@v1.0.0-eol")
			r = upgraded
		})
		g.It("should get promoted to unsupported", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNil()
			g.Assert(upgraded.Kind).Equal(state.ReleaseKindUnsupported)
			g.Assert(upgraded.Tag).Equal("v1.0.0-unsupported")
			g.Assert(upgraded.String()).Equal("unsupported@v1.0.0-unsupported")
			r = upgraded
		})
		g.It("should not get promoted from unsupported", func() {
			upgraded, err := r.Promote()
			g.Assert(err).IsNotNil()
			g.Assert(upgraded).IsNil()
		})
	})
}
