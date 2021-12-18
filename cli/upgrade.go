package cli

/*
func NewUpgradeCmd() *cobra.Command {
	type Opts struct {
		Kind           string
		IncMajor       bool
		IncMinor       bool
		IncPatch       bool
		WithService    string
		WithoutService string
	}
	opts := new(Opts)
	cmd := &cobra.Command{
		Use: "upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, err := state.NewReleaseKind(opts.Kind)
			if err != nil {
				return err
			}
			next := store.Latest(kind).Copy()
			version, err := semver.NewVersion(next.Tag)
			if err != nil {
				return err
			}
			var nextVersion semver.Version
			switch {
			case opts.IncMajor:
				nextVersion = version.IncMajor()
			case opts.IncMinor:
				nextVersion = version.IncMinor()
			case opts.IncPatch:
				nextVersion = version.IncPatch()
			}
			if opts.WithService != "" {
				if err := addServices(next.Versions, opts.WithService); err != nil {
					return err
				}
			}
			if opts.WithoutService != "" {
				if err := removeServices(next.Versions, opts.WithoutService); err != nil {
					return err
				}
			}
			release, err := state.NewRelease(kind, nextVersion.String(), next.Versions)
			if err != nil {
				return err
			}
			if err := states.CreateRelease(release); err != nil {
				return err
			}
			if err := states.Export(FileName); err != nil {
				return err
			}
			return nil
		},
	}
	fl := cmd.Flags()
	fl.BoolVarP(&opts.IncMajor, "major", "", false, "")
	fl.BoolVarP(&opts.IncMinor, "minor", "", false, "")
	fl.BoolVarP(&opts.IncPatch, "patch", "", false, "")
	fl.StringVarP(&opts.Kind, "kind", "", "", "")
	fl.StringVarP(&opts.WithService, "with-service", "", "", "")
	fl.StringVarP(&opts.WithoutService, "without-service", "", "", "")
	return cmd
}

func addServices(m state.VersionMap, s string) error {
	serviceStr := strings.Split(s, ",")
	for _, v := range serviceStr {
		parts := strings.Split(v, "@")
		if len(parts) != 2 {
			return errors.New("cli: could not parse service string")
		}
		m.Set(parts[0], parts[1])
	}
	return nil
}

func removeServices(m state.VersionMap, s string) error {
	serviceStr := strings.Split(s, ",")
	for _, v := range serviceStr {
		parts := strings.Split(v, "@")
		m.Remove(parts[0])
	}
	return nil
}
*/
