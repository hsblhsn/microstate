package cli

import (
	"fmt"

	"github.com/hsblhsn/microstate/state"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	var (
		store = state.NewState()
	)
	return &cobra.Command{
		Use: "status",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Import(FileName); err != nil {
				return eris.Wrap(err, "cli: could not import state file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			for i := state.ReleaseKindDev; i < state.ReleaseKindUnsupported; i++ {
				latest := store.Latest(i)
				if len(latest.Versions) != 0 {
					fmt.Printf("%-6s :: %s\n", latest.Kind.String(), latest.Tag)
				}
			}
			return nil
		},
	}
}
