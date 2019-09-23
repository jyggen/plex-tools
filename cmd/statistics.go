package cmd

import (
	"github.com/jyggen/plex-tools/plex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Display statistics about a library",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlag("format", cmd.Flags().Lookup("format")); err != nil {
			return err
		}

		if err := viper.BindPFlag("library", cmd.Flags().Lookup("library")); err != nil {
			return err
		}

		if err := viper.BindPFlag("server", cmd.Flags().Lookup("server")); err != nil {
			return err
		}

		if err := viper.BindPFlag("token", cmd.Flags().Lookup("token")); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := plex.New(viper.GetString("token"))

		if err != nil {
			return err
		}

		serverName := viper.GetString("server")

		var server *plex.Server

		if serverName == "" {
			server, err = p.PromptForServer()
		} else {
			server, err = p.GetServerByName(serverName)
		}

		if err != nil {
			return err
		}

		p.UseServer(server)

		libraryKey := viper.GetString("library")

		if libraryKey == "" {
			libraryKey, err = p.PromptForLibraryKey()

			if err != nil {
				return err
			}
		}

		probe, err := p.Probe(libraryKey)

		if err != nil {
			return err
		}

		probe.Statistics().Ascii(os.Stdout)

		return nil
	},
}

func init() {
	statisticsCmd.Flags().String("format", "ascii", "output format")
	statisticsCmd.Flags().String("library", "", "Plex library key")
	statisticsCmd.Flags().String("server", "", "Plex server name")
	statisticsCmd.Flags().String("token", "", "Plex access token")
	rootCmd.AddCommand(statisticsCmd)
}
