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

	ckpCommand.AddCommand(NewInitCommand(config))
	ckpCommand.AddCommand(NewResetCommand(config))
	ckpCommand.AddCommand(NewAddCommand(config))
	ckpCommand.AddCommand(NewListCommand(config))
	ckpCommand.AddCommand(NewFindCommand(config))
	ckpCommand.AddCommand(NewPushCommand(config))
	ckpCommand.AddCommand(NewPullCommand(config))
	ckpCommand.AddCommand(NewRmCommand(config))
	ckpCommand.AddCommand(NewEditCommand(config))
	ckpCommand.AddCommand(NewRunCommand(config))
	return ckpCommand
}
