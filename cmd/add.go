package cmd

import (
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewAddCommand adds everything that written after --code or --solution flag
func NewAddCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "add",
		Short: "will store your code snippets or solutions",
		Long: `will store your code snippets or solutions

	example: ckp add code 'echo je suis con'
	Will store 'echo je suis con' as a code asset in your solution repository


	example: ckp add solution 'https://opensource.code'
	Will store 'https://opensource.code' as a solution asset in your solution repository
`,
	}

	//Add commands
	command.AddCommand(NewAddCodeCommand(conf))
	command.AddCommand(NewAddSolutionCommand(conf))
	command.AddCommand(NewAddHistoryCommand(conf))

	//Set flags
	command.PersistentFlags().StringP("comment", "m", "", `ckp add -m <comment>`)
	return command
}
