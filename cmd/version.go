package cmd

import (
	"fmt"

	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//This version will be set by a goreleaser ldflag
var version = "0.0.0.dev"
var versionString = `Version: %s
Build by elhmn
Support osscameroon here https://opencollective.com/osscameroon
`

//NewVersionCommand output program version
func NewVersionCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Show ckp version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := versionCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func versionCommand(conf config.Config) error {
	fmt.Fprintf(conf.OutWriter, versionString, version)
	return nil
}
