package cmd

import (
	"fmt"

	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewUpdateCommand will update your binary to the latest release
func NewUpdateCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "Update the binary to the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func updateCommand(conf config.Config) error {
	fmt.Fprintf(conf.OutWriter, "update")
	return nil
}
