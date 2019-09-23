package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock all metadata fields",
	Long: `This tool unlocks all metadata fields in Plex so Plex and its agents are able to
modify them.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
}
