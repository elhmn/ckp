package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

//NewStoreCodeCommand stores everything that written after --code or --solution flag
func NewStoreCodeCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "code [your_code]",
		Aliases: []string{"c"},
		Short:   "store code will store your code",
		Long: `store code will store your code in your solution repository

	example: ckp store code 'echo this is my command'
	Will store 'echo this is my command' as a code asset in your solution repository

	example: ckp store code -p ./path_to_you_code.sh
	Will store the code in path_to_you_code.sh as a code asset in your solution repository
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := storeCodeCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().StringP("alias", "a", "", `ckp store -a <alias>`)

	return command
}

func storeCodeCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	cmd.Flags().Parse(args)
	flags := cmd.Flags()
	code := strings.Join(args, " ")

	storeFile, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed to get the store file path: %s", err)
	}

	storeData, storeBytes, err := store.LoadStore(storeFile)
	if err != nil {
		return fmt.Errorf("failed to laod store: %s", err)
	}

	tempFile, err := config.GetTempStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed to get the store temporary file path: %s", err)
	}

	//Copy the store file to a temporary destination
	if err := ioutil.WriteFile(tempFile, storeBytes, 0666); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", tempFile, err)
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

	//Save the new `Store` struct in the store file
	if err = storeData.SaveStore(storeFile); err != nil {
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

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}

	fmt.Fprintln(conf.OutWriter, "Your code was successfully stored!")
	return nil
}

//createNewCodeScriptEntry return a new code Script entry
func createNewCodeScriptEntry(code string, flags *flag.FlagSet) (store.Script, error) {
	timeNow := time.Now()

	//Get data from flags
	path, err := flags.GetString("path")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `path` flag: %s", err)
	}
	alias, err := flags.GetString("alias")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `alias` flag: %s", err)
	}
	comment, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID(code, path, comment, alias, "")
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	return store.Script{
		ID:           id,
		Comment:      comment,
		CreationTime: timeNow,
		UpdateTime:   timeNow,
		Code: store.Code{
			Content:  code,
			Alias:    alias,
			FilePath: path,
		},
	}, nil
}
