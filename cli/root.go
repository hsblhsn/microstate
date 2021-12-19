package cli

import (
	"fmt"
	"os"

	"github.com/hsblhsn/microstate/state"
	"github.com/spf13/cobra"
)

const FileName = state.DefaultFileName

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "microstate",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Use --help to see available commands.")
			os.Exit(1)
			return nil
		},
	}
	publish := &cobra.Command{
		Use: "publish",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Use --help to see available commands.")
			os.Exit(1)
			return nil
		},
	}
	init := NewInitCmd()
	rollback := NewRollbackCmd()
	status := NewStatusCmd()
	dev := NewDevCmd()
	alpha := NewAlphaCmd()
	beta := NewBetaCmd()
	rc := NewRCCmd()
	ga := NewGACmd()
	eol := NewEOLCmd()
	unsupported := NewUnsupportedCmd()
	publish.AddCommand(dev, alpha, beta, rc, ga, eol, unsupported)
	cmd.AddCommand(init, status, publish, rollback)
	return cmd
}
