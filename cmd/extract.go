package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract .rar files",
	Long: `This tool scans the input path and extracts all .rar files found within. In
order to prevent software that's monitoring the filesystem from picking up files
that are currently being extracted, the .rar files are extracted to a temporary
path and then the files within are moved back into the source path once
extracted.'`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
