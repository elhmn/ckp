package cmd

import (
	"fmt"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
)

//NewRunCommand create new cobra command for the run command
func NewRunCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "run [code_id]",
		Short: "runs your code entries from the store",
		Long: `runs your code entries from the store

		example: ckp run
		Will prompt an interactive UI that will allow you to search and run
		a code entry

		example: ckp run <entry_id>
		Will run the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().Bool("from-history", false, `list code and solution records from history`)

	return command
}

func runCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	var entryID string
	if len(args) >= 1 {
		entryID = args[0]
	}

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

	_, storeData, _, err := loadStore(storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	index, err := getScriptEntryIndex(conf, storeData.Scripts, entryID, store.EntryTypeCode)
	if err != nil {
		return fmt.Errorf("failed to get script `%s` entry index: %s", entryID, err)
	}

	if err := runCodeEntry(conf, storeData.Scripts, index); err != nil {
		return fmt.Errorf("failed to run code entry: %s", err)
	}
	return nil
}

func runCodeEntry(conf config.Config, scripts []store.Script, index int) error {
	script := scripts[index]
	if script.Code.Content == "" {
		return fmt.Errorf("might not be a code entry, nothing to run")
	}

	return conf.Exec.RunInteractive("bash", "-c", script.Code.Content)
}
