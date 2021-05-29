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
		Short: "run run your code entries from the store",
		Long: `run run your code entries from the store

		example: ckp run
		Will prompt an interactive UI that will allow you to search and run
		a code entry

		example: ckp run <entry_id>
		Will run the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			var entryID string
			if len(args) >= 1 {
				entryID = args[0]
			}

			if err := runCommand(conf, entryID); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func runCommand(conf config.Config, entryID string) error {
	_, storeData, _, err := loadStore(conf)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	index, err := getScriptEntryIndex(conf, storeData.Scripts, entryID)
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
