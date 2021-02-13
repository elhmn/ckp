package cmd

import (
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewCKPCommand returns a new ckp cobra command
func NewCKPCommand(config config.Config) *cobra.Command {
	var ckpCommand = &cobra.Command{
		Use:   "ckp",
		Short: "ckp saves and pull your bash scripts",
		Long:  `ckp is a tool that helps you save and fetch the bash scripts you use frequently`,
	}

	ckpCommand.AddCommand(newInitCommand(config))
	return ckpCommand
}
