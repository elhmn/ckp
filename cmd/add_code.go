package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const editorFileCodeTemplate = `## comment: %s
## alias: %s
%s
`

//NewAddCodeCommand adds everything that written after --code or --solution flag
func NewAddCodeCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "code <your_code>",
		Aliases: []string{"c"},
		Short:   "will store your code",
		Long: `will store your code in your solution repository

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
	command.PersistentFlags().StringP("path", "p", "", `ckp add -p <script_path>`)

	return command
}

func addCodeCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}
	flags := cmd.Flags()
	code := strings.Join(args, " ")

	//Check if the code entry contains sensitive data
	if ret, word := store.HasSensitiveData(code); ret {
		fmt.Fprintf(conf.OutWriter, "Found the keyword `%s` in %s\n", word, code)
		if !printers.Confirm("Add anyway ?") {
			fmt.Fprintf(conf.OutWriter, "Code entry addition was aborted!\n")
			return nil
		}
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

	conf.Spin.Message(" pulling remote changes...")
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	conf.Spin.Message(" remote changes pulled")

	conf.Spin.Message(" adding new code entry...")
	storeFile, storeData, storeBytes, err := loadStore(storeFilePath)
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

	//if it is an interactive update
	if len(args) == 0 {
		conf.Spin.Stop()
		s, err := getNewEntryDataFromFile(conf, script, CodeEntryTemplateType)
		if err != nil {
			return fmt.Errorf("failed to get new entry from the editor %s", err)
		}

		//Generate an ID for the newly added script
		s.ID, err = store.GenereateIdempotentID(s.Code.Content, s.Comment, s.Code.Alias, "")
		if err != nil {
			return fmt.Errorf("failed to generate idem potent ID: %s", err)
		}

		script.ID = s.ID
		script.Code = s.Code
		script.Comment = s.Comment
		conf.Spin.Start()
	}

	//Check if entry already exist
	if storeData.EntryAlreadyExist(script.ID) {
		//Delete the temporary file
		if err := os.RemoveAll(tempFile); err != nil {
			return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
		}

		return fmt.Errorf("An identical record was found in the storage, please try `ckp edit %s`", script.ID)
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
	conf.Spin.Message(" new entry successfully added")

	//Push changes
	conf.Spin.Message(" pushing local changes...")
	err = pushLocalChanges(conf, dir, commitAddAction, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message(" local changes pushed")

	fmt.Fprintln(conf.OutWriter, "\nYour code was successfully added!")
	fmt.Fprintf(conf.OutWriter, "\n%s", script)
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

func createHistoryTempFile(conf config.Config, storeBytes []byte) (string, error) {
	tempFile, err := config.GetTempHistoryStoreFilePath(conf)
	if err != nil {
		return tempFile, fmt.Errorf("failed to get the history store temporary file path: %s", err)
	}

	//Copy the store file to a temporary destination
	if err := ioutil.WriteFile(tempFile, storeBytes, 0666); err != nil {
		return tempFile, fmt.Errorf("failed to write to file %s: %s", tempFile, err)
	}

	return tempFile, nil
}

func loadStore(storeFile string) (string, *store.Store, []byte, error) {
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
	path, err := flags.GetString("path")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `path` flag: %s", err)
	}

	if path != "" {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return store.Script{}, fmt.Errorf("failed to read %s: %s", path, err)
		}
		code = string(bytes)
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
