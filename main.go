package main

import (
	"github.com/jyggen/plex-tools/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
