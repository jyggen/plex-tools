package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var version string

var appName = filepath.Base(os.Args[0])
var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   fmt.Sprintf("%s is a collection of tools and utilities for Plex Media Server administrators", appName),
	Version: version,
	Args:    cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfgFile := cmd.Flag("config").Value.String()

		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		}

		if viper.ConfigFileUsed() == "" {
			home, err := homedir.Dir()

			if err != nil {
				return err
			}

			viper.SetConfigFile(filepath.Join(home, fmt.Sprintf(".%s.yaml", appName)))
		}

		_ = viper.ReadInConfig()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", fmt.Sprintf("config file (default \"$HOME/.%s.yaml\")", appName))
}

func Execute() error {
	return rootCmd.Execute()
}
