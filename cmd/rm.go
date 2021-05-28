package cmd

import (
	"fmt"
	"os"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

//NewRmCommand create new cobra command for the rm command
func NewRmCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "rm [code_id | solution_id]",
		Short: "rm removes code or solution entries from the store",
		Long: `rm removes code or solution entries from the store

		example: ckp rm
		Will prompt an interactive UI that will allow you to search and delete
		a code or solution entry

		example: ckp rm <entry_id>
		Will remove the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			var entryID string
			if len(args) >= 1 {
				entryID = args[0]
			}

			if err := rmCommand(conf, entryID); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func rmCommand(conf config.Config, entryID string) error {
	//Setup spinner
	conf.Spin.Start()
	defer conf.Spin.Stop()

	dir, err := config.GetStoreDirPath(conf)
	if err != nil {
		return fmt.Errorf("failed get repository path: %s", err)
	}

	storeFilePath, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed get store file path: %s", err)
	}

	conf.Spin.Message(" pulling remote changes...")
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	conf.Spin.Message(" remote changes pulled")

	conf.Spin.Message(" removing changes")
	storeFile, storeData, storeBytes, err := loadStore(conf)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	if entryID == "" {
		conf.Spin.Stop()
		if entryID, err = selectScriptEntry(storeData.Scripts); err != nil {
			return fmt.Errorf("failed to select entry: %s", err)
		}
	}

	//Remove script entry
	storeData.Scripts, err = removeScriptEntry(storeData.Scripts, entryID)
	if err != nil {
		return fmt.Errorf("failed to remove script entry: %s", err)
	}

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
	err = pushLocalChanges(conf, dir, storeFilePath, commitRemoveAction)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message(" local changes pushed")

	fmt.Fprintf(conf.OutWriter, "\nckp store was pushed successfully\n")
	return nil
}

func removeScriptEntry(scripts []store.Script, entryID string) ([]store.Script, error) {
	for index, s := range scripts {
		if entryID == s.ID {
			return append(scripts[:index], scripts[index+1:]...), nil
		}
	}

	return scripts, fmt.Errorf("`%s` entry not found", entryID)
}

func selectScriptEntry(scripts []store.Script) (string, error) {
	searchScript := func(input string, index int) bool {
		s := scripts[index]
		return doesScriptContain(s, input)
	}

	prompt := promptui.Select{
		Label:             "Enter your search text",
		Items:             scripts,
		Size:              5,
		StartInSearchMode: true,
		Searcher:          searchScript,
		Templates:         getTemplates(),
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed %v", err)
	}

	return scripts[i].ID, nil
}