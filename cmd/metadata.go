package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Update Plex metadata",
	Long: `This tool updates your media's Plex metadata with information fetched from
various online sources. Much like the built-in agents in Plex itself, but with a
lot more flexibility and control!`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
