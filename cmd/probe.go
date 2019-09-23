package cmd

import (
	"fmt"
	"github.com/jyggen/plex-tools/plex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var probeCmd = &cobra.Command{
	Use:  "probe",
	Args: cobra.NoArgs,
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

		switch viper.GetString("format") {
		case "ascii":
			probe.Ascii(os.Stdout)
		case "html":
			if err := probe.Html(os.Stdout); err != nil {
				return err
			}
		default:
			return fmt.Errorf("\"%s\" is not a supported output format", viper.GetString("format"))
		}

		return nil
	},
}

func init() {
	probeCmd.Flags().String("format", "ascii", "output format")
	probeCmd.Flags().String("library", "", "Plex library key")
	probeCmd.Flags().String("server", "", "Plex server name")
	probeCmd.Flags().String("token", "", "Plex access token")
	rootCmd.AddCommand(probeCmd)
}
