package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
