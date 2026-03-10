package cmd

import (
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewSyncCommand creates a new cobra command for the sync command
func NewSyncCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "sync",
		Short: "will sync your data with external tools",
		Long: `will sync your data with external tools

	example: ckp sync alias
	Will sync your aliases in a ~/.ckp/ckp_aliases file and reference them in your main rc file
`,
	}

	//Add commands
	command.AddCommand(NewSyncAliasCommand(conf))

	return command
}
