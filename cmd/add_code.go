package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

//NewAddCodeCommand adds everything that written after --code or --solution flag
func NewAddCodeCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "code <your_code>",
		Aliases: []string{"c"},
		Short:   "add code will store your code",
		Long: `add code will store your code in your solution repository

	example: ckp add code 'echo this is my command'
	Will store 'echo this is my command' as a code asset in your solution repository

	example: ckp add code -p ./path_to_you_code.sh
	Will store the code in path_to_you_code.sh as a code asset in your solution repository
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := addCodeCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().StringP("alias", "a", "", `ckp add -a <alias>`)

	return command
}

func addCodeCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	cmd.Flags().Parse(args)
	flags := cmd.Flags()
	code := strings.Join(args, " ")

	storeFile, storeData, storeBytes, err := loadStore(conf)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	tempFile, err := createTempFile(conf, storeBytes)
	if err != nil {
		return fmt.Errorf("failed to create tempFile: %s", err)
	}

	script, err := createNewCodeScriptEntry(code, flags)
	if err != nil {
		return fmt.Errorf("failed to create new script entry: %s", err)
	}

	if storeData.EntryAlreadyExist(script.ID) {
		return fmt.Errorf("An identical record was found in the storage, please try `ckp edit --id %s`", script.ID)
	}

	//Add new script entry in the `Store` struct
	storeData.Scripts = append(storeData.Scripts, script)

	//Save storeData in store
	if err := saveStore(storeData, storeBytes, storeFile, tempFile); err != nil {
		return fmt.Errorf("failed to save store in %s:  %s", storeFile, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}

	fmt.Fprintln(conf.OutWriter, "Your code was successfully added!")
	return nil
}

func createTempFile(conf config.Config, storeBytes []byte) (string, error) {
	tempFile, err := config.GetTempStoreFilePath(conf)
	if err != nil {
		return tempFile, fmt.Errorf("failed to get the store temporary file path: %s", err)
	}

	//Copy the store file to a temporary destination
	if err := ioutil.WriteFile(tempFile, storeBytes, 0666); err != nil {
		return tempFile, fmt.Errorf("failed to write to file %s: %s", tempFile, err)
	}

	return tempFile, nil
}

func loadStore(conf config.Config) (string, *store.Store, []byte, error) {
	storeFile, err := config.GetStoreFilePath(conf)
	if err != nil {
		return storeFile, nil, nil, fmt.Errorf("failed to get the store file path: %s", err)
	}

	storeData, storeBytes, err := store.LoadStore(storeFile)
	if err != nil {
		return storeFile, storeData, storeBytes, fmt.Errorf("failed to load store: %s", err)
	}
	return storeFile, storeData, storeBytes, nil
}

func saveStore(storeData *store.Store, storeBytes []byte, storeFile, tempFile string) error {
	//Save the new `Store` struct in the store file
	if err := storeData.SaveStore(storeFile); err != nil {
		//if we failed to write the store file then we restore the file to its original content
		if err1 := ioutil.WriteFile(storeFile, storeBytes, 0666); err1 != nil {
			return fmt.Errorf("failed to write to file %s: %s", storeFile, err)
		}

		//Delete the temporary file
		if err := os.RemoveAll(tempFile); err != nil {
			return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
		}

		return fmt.Errorf("failed to write to file %s: %s", storeFile, err)
	}

	return nil
}

//createNewCodeScriptEntry return a new code Script entry
func createNewCodeScriptEntry(code string, flags *flag.FlagSet) (store.Script, error) {
	timeNow := time.Now()

	alias, err := flags.GetString("alias")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `alias` flag: %s", err)
	}
	comment, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID(code, comment, alias, "")
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	return store.Script{
		ID:           id,
		Comment:      comment,
		CreationTime: timeNow,
		UpdateTime:   timeNow,
		Code: store.Code{
			Content: code,
			Alias:   alias,
		},
	}, nil
}
