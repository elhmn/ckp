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
	cmd.Flags().Parse(args)
	flags := cmd.Flags()
	solution := strings.Join(args, " ")

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

	script, err := createNewSolutionScriptEntry(solution, flags)
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

	fmt.Fprintln(conf.OutWriter, "Your solution was successfully added!")
	return nil
}

//createNewSolutionScriptEntry return a new solution entry
func createNewSolutionScriptEntry(solution string, flags *flag.FlagSet) (store.Script, error) {
	timeNow := time.Now()

	//Get data from flags
	path, err := flags.GetString("path")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `path` flag: %s", err)
	}
	comment, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID("", path, comment, "", solution)
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
