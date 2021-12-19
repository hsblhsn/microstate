package cli

import (
	"fmt"

	"github.com/hsblhsn/microstate/state"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func NewRollbackCmd() *cobra.Command {
	var (
		store  = state.NewState()
		logger = NewLogger()
	)
	return &cobra.Command{
		Use: "rollback",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Import(FileName); err != nil {
				return eris.Wrap(err, "cli: could not import state file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store.Rollback()
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Export(FileName); err != nil {
				return eris.Wrap(err, "cli: could not export state file")
			}
			top, err := store.Head()
			if err != nil {
				return eris.Wrap(err, "cli: could not get head release")
			}
			logger.OK(fmt.Sprintf("rolled back to %s", top.String()))
			return nil
		},
	}
}
