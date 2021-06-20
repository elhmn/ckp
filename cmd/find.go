package cmd

import (
	"fmt"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
)

//NewFindCommand display a prompt for you to enter the code or solution you are looking for
func NewFindCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "find",
		Aliases: []string{"f"},
		Short:   "find your code snippets and solutions",
		Long: `find your code snippets and solutions

	example: ckp find
	Will display a prompt for you to enter the code or solution you are looking for
	`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := findCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().Bool("from-history", false, `list code and solution records from history`)

	return command
}

func findCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}
	flags := cmd.Flags()
	fromHistory, err := flags.GetBool("from-history")
	if err != nil {
		return fmt.Errorf("could not parse `fromHistory` flag: %s", err)
	}

	//Get the store file path
	var storeFilePath string
	if !fromHistory {
		storeFilePath, err = config.GetStoreFilePath(conf)
		if err != nil {
			return fmt.Errorf("failed to get the store file path: %s", err)
		}
	} else {
		storeFilePath, err = config.GetHistoryFilePath(conf)
		if err != nil {
			return fmt.Errorf("failed to get the history store file path: %s", err)
		}
	}

	storeData, _, err := store.LoadStore(storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to laod store: %s", err)
	}

	scripts := storeData.Scripts

	_, result, err := conf.Printers.SelectScriptEntry(scripts)
	if err != nil {
		return fmt.Errorf("prompt failed %v", err)
	}

	fmt.Fprintf(conf.OutWriter, "\n%s", result)
	return nil
}
