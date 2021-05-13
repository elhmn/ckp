package cmd

import (
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewStoreCommand stores everything that written after --code or --solution flag
func NewStoreCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "store",
		Short: "`ckp store` will store your solution or code",
		Long: `ckp store will store your solution or code in your solution repository

	example: ckp store code 'echo je suis con'
	Will store 'echo je suis con' as a code asset in your solution repository


	example: ckp store solution 'https://opensource.code'
	Will store 'https://opensource.code' as a solution asset in your solution repository
`,
	}

	//Add commands
	command.AddCommand(NewStoreCodeCommand(conf))
	command.AddCommand(NewStoreSolutionCommand(conf))

	//Set flags
	command.PersistentFlags().StringP("comment", "m", "", `ckp store -m <comment>`)
	command.PersistentFlags().StringP("path", "p", "", `ckp store -p <path_to_your_code_or_solution>`)

	return command
}
