package cmd

import (
	"fmt"
	"os"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
)

//NewRmCommand create new cobra command for the rm command
func NewRmCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "rm [code_id | solution_id]",
		Short: "removes code or solution entries from the store",
		Long: `removes code or solution entries from the store

		example: ckp rm
		Will prompt an interactive UI that will allow you to search and delete
		a code or solution entry

		example: ckp rm <entry_id>
		Will remove the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := rmCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().Bool("from-history", false, `list code and solution records from history`)

	return command
}

func rmCommand(cmd *cobra.Command, args []string, conf config.Config) error {
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

	//Setup spinner
	conf.Spin.Start()
	defer conf.Spin.Stop()

	dir, err := config.GetStoreDirPath(conf)
	if err != nil {
		return fmt.Errorf("failed get repository path: %s", err)
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

	conf.Spin.Message(" pulling remote changes...")
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	conf.Spin.Message(" remote changes pulled")

	conf.Spin.Message(" removing changes")
	storeFile, storeData, storeBytes, err := loadStore(storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	index, err := getScriptEntryIndex(conf, storeData.Scripts, entryID)
	if err != nil {
		return fmt.Errorf("failed to get script `%s` entry index: %s", entryID, err)
	}

	//Remove script entry
	storeData.Scripts = removeScriptEntry(storeData.Scripts, index)

	tempFile, err := createTempFile(conf, storeBytes)
	if err != nil {
		return fmt.Errorf("failed to create tempFile: %s", err)
	}

	//Save storeData in store
	if err := saveStore(storeData, storeBytes, storeFile, tempFile); err != nil {
		return fmt.Errorf("failed to save store in %s:  %s", storeFile, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}

	conf.Spin.Message(" pushing local changes...")
	err = pushLocalChanges(conf, dir, commitRemoveAction, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message(" local changes pushed")

	fmt.Fprintf(conf.OutWriter, "\nentry was removed successfully\n")
	return nil
}

func removeScriptEntry(scripts []store.Script, index int) []store.Script {
	return append(scripts[:index], scripts[index+1:]...)
}

func getScriptEntryIndex(conf config.Config, scripts []store.Script, entryID string) (int, error) {
	if entryID == "" {
		conf.Spin.Stop()
		index, _, err := printers.SelectScriptEntry(scripts)
		if err != nil {
			return index, fmt.Errorf("failed to select entry: %s", err)
		}
		return index, nil
	}

	for index, s := range scripts {
		if entryID == s.ID {
			return index, nil
		}
	}

	return -1, fmt.Errorf("entry not found")
}
