package cli

import (
	"log"

	"github.com/spf13/cobra"
)

func NewDevCmd() *cobra.Command {
	return &cobra.Command{
		Use: "dev",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("dev")
			return nil
		},
	}
}
