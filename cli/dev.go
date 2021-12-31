package cli

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hsblhsn/microstate/state"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func NewDevCmd() *cobra.Command {
	var (
		store  = state.NewState()
		logger = NewLogger()
	)
	var (
		from     string
		fromKind string
		services []string
		major    bool
		minor    bool
		patch    bool
	)
	cmd := &cobra.Command{
		Use: "dev",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Import(FileName); err != nil {
				return eris.Wrap(err, "cli: could not import state file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			versionMap := state.NewVersionMap()
			for _, v := range services {
				v = strings.TrimSpace(v)
				parts := strings.SplitN(v, "@", 2)
				if len(parts) < 2 {
					return eris.New("invalid service: " + v)
				}
				versionMap.Set(parts[0], parts[1])
			}
			latestDev := store.Latest(state.ReleaseKindDev)
			latestDevVersion, err := semver.NewVersion(latestDev.Tag)
			if err != nil {
				return eris.Wrap(err, "cli: could not parse latest dev release version")
			}
			var nextTag string
			{
				var nextVersion semver.Version
				if major {
					nextVersion = latestDevVersion.IncMajor()
				} else if minor {
					nextVersion = latestDevVersion.IncMinor()
				} else if patch {
					nextVersion = latestDevVersion.IncPatch()
				} else {
					nextVersion = latestDevVersion.IncPatch()
				}
				nextVersion, err = nextVersion.SetPrerelease(state.ReleaseKindDev.String())
				if err != nil {
					return eris.Wrap(err, "cli: could not set prerelease info on version tag")
				}
				nextTag = "v" + nextVersion.String()
			}
			if from != "" || fromKind != "" {
				var fromRelease *state.Release
				if from != "" {
					hash, err := state.NewHash(from)
					if err != nil {
						return eris.Wrap(err, "cli: hash is not valid")
					}
					fromRelease, err = store.GetRelease(hash)
					if err != nil {
						return eris.Wrap(err, "cli: could not get release")
					}
				} else if fromKind != "" {
					kind, err := state.NewReleaseKindFromString(fromKind)
					if err != nil {
						return eris.Wrap(err, "cli: hash is not valid")
					}
					fromRelease = store.Latest(kind)
				} else {
					return eris.New("cli: --from or --from-kind must be specified")
				}
				// replace original versions with new ones
				for k, v := range versionMap {
					fromRelease.Versions[k] = v
				}
				fromRelease.Tag = nextTag
				fromRelease.Kind = state.ReleaseKindDev
				if err := store.CreateRelease(fromRelease); err != nil {
					return eris.Wrap(err, "cli: could not create release")
				}
			} else {
				r, err := state.NewRelease(state.ReleaseKindDev, nextTag, versionMap)
				if err != nil {
					return eris.Wrap(err, "cli: could not build release")
				}
				if err := store.CreateRelease(r); err != nil {
					return eris.Wrap(err, "cli: could not create release")
				}
			}
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Export(FileName); err != nil {
				return eris.Wrap(err, "cli: could not export state file")
			}
			logger.OK(fmt.Sprintf("dev release created: %s", store.Latest(state.ReleaseKindDev).Tag))
			fmt.Print(store.Latest(state.ReleaseKindDev).Tag)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fromKind, "from-kind", "k", "", "service name and version")
	cmd.Flags().StringVarP(&from, "from", "f", "", "Copy from service.")
	cmd.Flags().BoolVarP(&major, "major", "", false, "Major version upgrade")
	cmd.Flags().BoolVarP(&minor, "minor", "", false, "Minor version upgrade")
	cmd.Flags().BoolVarP(&patch, "patch", "", false, "Patch version upgrade")
	cmd.Flags().StringArrayVarP(&services, "service", "s", make([]string, 0),
		"Service name and version. It accepts array of values. (e.g. --service serviceA@v1.0 --service serviceB@v1.0)",
	)
	return cmd
}
