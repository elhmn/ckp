package cmd

import (
	"fmt"
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

	storeFile, storeData, storeBytes, err := loadStore(conf)
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

	fmt.Fprintln(conf.OutWriter, "Your solution was successfully added!")
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