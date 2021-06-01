package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

//NewAddSolutionCommand adds everything that written after --solution or --solution flag
func NewAddSolutionCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "solution <your_solution>",
		Aliases: []string{"s"},
		Short:   "add solution will store your solution",
		Long: `add solution will store your solution in your solution repository

	example: ckp add solution 'echo this is my command'
	Will store 'echo this is my command' as a solution asset in your solution repository
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := addSolutionCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func addSolutionCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}
	flags := cmd.Flags()
	solution := strings.Join(args, " ")

	//Check if the code entry contains sensitive data
	if ret, word := store.HasSensitiveData(solution); ret {
		fmt.Fprintf(conf.OutWriter, "Found the keyword `%s` in %s\n", word, solution)
		if !printers.Confirm("Add anyway ?") {
			fmt.Fprintf(conf.OutWriter, "Solution entry addition was aborted!\n")
			return nil
		}
	}

	//Setup spinner
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	dir, err := config.GetStoreDirPath(conf)
	if err != nil {
		return fmt.Errorf("failed get repository path: %s", err)
	}

	storeFilePath, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed get store file path: %s", err)
	}

	spin.Suffix = " pulling remote changes..."
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	spin.Suffix = " remote changes pulled"

	spin.Suffix = " adding new solution entry..."
	_, storeData, storeBytes, err := loadStore(storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	tempFile, err := createTempFile(conf, storeBytes)
	if err != nil {
		return fmt.Errorf("failed to create tempFile: %s", err)
	}

	script, err := createNewSolutionScriptEntry(solution, flags)
	if err != nil {
		return fmt.Errorf("failed to create new script entry: %s", err)
	}

	if storeData.EntryAlreadyExist(script.ID) {
		//Delete the temporary file
		if err := os.RemoveAll(tempFile); err != nil {
			return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
		}

		return fmt.Errorf("An identical record was found in the storage, please try `ckp edit --id %s`", script.ID)
	}

	//Add new script entry in the `Store` struct
	storeData.Scripts = append(storeData.Scripts, script)

	//Save storeData in store
	if err := saveStore(storeData, storeBytes, storeFilePath, tempFile); err != nil {
		return fmt.Errorf("failed to save store in %s:  %s", storeFilePath, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}
	spin.Suffix = " new entry successfully added"

	spin.Suffix = " pushing local changes..."
	err = pushLocalChanges(conf, dir, commitAddAction, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	spin.Suffix = " local changes pushed"

	fmt.Fprintln(conf.OutWriter, "\nYour solution was successfully added!")
	return nil
}

//createNewSolutionScriptEntry return a new solution entry
func createNewSolutionScriptEntry(solution string, flags *flag.FlagSet) (store.Script, error) {
	timeNow := time.Now()

	comment, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID("", comment, "", solution)
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	return store.Script{
		ID:           id,
		Comment:      comment,
		CreationTime: timeNow,
		UpdateTime:   timeNow,
		Solution: store.Solution{
			Content: solution,
		},
	}, nil
}
