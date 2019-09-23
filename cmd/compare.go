package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:  "compare",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
