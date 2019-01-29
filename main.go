package main

import (
	"git-ghost/cmd"
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	// RootCmd prints errors if exists
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
