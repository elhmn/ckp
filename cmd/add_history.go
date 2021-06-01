package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/history"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
)

//NewAddHistoryCommand adds everything that written after --code or --solution flag
func NewAddHistoryCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "add history will store code from your shell history",
		Long: `add history will store code from your shell history
	it will read your .bash_history and zsh_history files and store
	every script oneliner as a code entry in your store.yaml file

	example: ckp history
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := addHistoryCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().BoolP("skip-secrets", "s", false, `will skip secrets and sensitive informations`)

	return command
}

func addHistoryCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}

	flags := cmd.Flags()

	//Get data from flags
	shouldSkipSecrets, err := flags.GetBool("skip-secrets")
	if err != nil {
		return fmt.Errorf("could not parse `--skip-secrets` flag: %s", err)
	}

	storeFile, storeData, storeBytes, err := loadStore(conf)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	tempFile, err := createTempFile(conf, storeBytes)
	if err != nil {
		return fmt.Errorf("failed to create tempFile: %s", err)
	}

	records, err := history.GetHistoryRecords()
	if err != nil {
		return fmt.Errorf("failed to get history records: %s", err)
	}

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

	conf.Spin.Message("pulling remote changes...")
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	conf.Spin.Message("remote changes pulled")

	//Add history code records
	addScriptsFromRecords(conf, records, storeData, shouldSkipSecrets)

	//Save storeData in store
	if err := saveStore(storeData, storeBytes, storeFile, tempFile); err != nil {
		return fmt.Errorf("failed to save store in %s:  %s", storeFile, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}

	conf.Spin.Message("pushing local changes...")
	err = pushLocalChanges(conf, dir, storeFilePath, commitAddAction)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message("local changes pushed")

	fmt.Fprintln(conf.OutWriter, "\nYour history was successfully added!")
	return nil
}

func addScriptsFromRecords(conf config.Config, records []string, storeData *store.Store, shouldSkipSecrets bool) *store.Store {
	size := len(records)
	for i, record := range records {
		conf.Spin.Message(fmt.Sprintf("%d/%d adding record...", i, size))

		//Check if the code entry contains sensitive data
		if ret, word := store.HasSensitiveData(record); ret {
			if shouldSkipSecrets {
				continue
			}

			fmt.Fprintf(conf.OutWriter, "Found the keyword `%s` in %s\n", word, record)
			fmt.Fprintf(conf.OutWriter, "%d/%d records\n", i, size)
			if !printers.Confirm("Add anyway ?") {
				fmt.Fprintf(conf.OutWriter, "Code entry addition was aborted!\n")
				continue
			}
		}

		//Read history file parse its content and store each entry
		script, err := createNewHistoryScriptEntry(record)
		if err != nil {
			fmt.Fprintf(conf.ErrWriter, "failed to create new script entry: %s", err)
			continue
		}

		if storeData.EntryAlreadyExist(script.ID) {
			continue
		}

		//Add new script entry in the `Store` struct
		storeData.Scripts = append(storeData.Scripts, script)

		conf.Spin.Message(fmt.Sprintf("%d/%d record successfully added!", i, size))
	}

	return storeData
}

//createNewHistoryScriptEntry return a new code Script entry
func createNewHistoryScriptEntry(code string) (store.Script, error) {
	timeNow := time.Now()
	//Generate script entry unique id
	id, err := store.GenereateIdempotentID(code, "", "", "")
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	return store.Script{
		ID:           id,
		CreationTime: timeNow,
		UpdateTime:   timeNow,
		Code: store.Code{
			Content: code,
		},
	}, nil
}
