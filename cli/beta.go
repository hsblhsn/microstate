package cli

import (
	"github.com/hsblhsn/microstate/state"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func NewBetaCmd() *cobra.Command {
	var (
		store  = state.NewState()
		logger = NewLogger()
	)
	return &cobra.Command{
		Use: "beta",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Import(FileName); err != nil {
				return eris.Wrap(err, "cli: could not import state file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := store.PromoteTo(state.ReleaseKindBeta); err != nil {
				return eris.Wrap(err, "cli: could not promote")
			}
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Export(FileName); err != nil {
				return eris.Wrap(err, "cli: could not export state file")
			}
			logger.Promotion(store, state.ReleaseKindBeta)
			return nil
		},
	}
}
