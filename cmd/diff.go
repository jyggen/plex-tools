package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Keep track of changes to your media",
	Long:  `This tool pulls down metadata about your media from Plex Media Server and compares it to the previous run in order to detect and keep track of changes.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
