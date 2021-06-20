package cmd

import (
	"fmt"

	"os"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//NewEditCommand create new cobra command for the edit command
func NewEditCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "edit [code_id | solution_id]",
		Short: "edit code or solution entries from the store",
		Long: `edit code or solution entries from the store

		example: ckp edit
		Will prompt an interactive UI that will allow you to search and delete
		a code or solution entry

		example: ckp edit <entry_id>
		Will edit the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := editCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().Bool("from-history", false, `list code and solution records from history`)
	command.PersistentFlags().StringP("comment", "m", "", `ckp edit -m <comment>`)
	command.PersistentFlags().StringP("alias", "a", "", `ckp edit -a <alias>`)
	command.PersistentFlags().StringP("code", "c", "", `ckp edit -c <code>`)
	command.PersistentFlags().StringP("solution", "s", "", `ckp edit -s <solution>`)
	return command
}

func editCommand(cmd *cobra.Command, args []string, conf config.Config) error {
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
	storeData.Scripts, err = editScriptEntry(flags, storeData.Scripts, index)
	if err != nil {
		return fmt.Errorf("failed to editScriptEntry: %s", err)
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
	err = pushLocalChanges(conf, dir, commitRemoveAction, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message(" local changes pushed")

	fmt.Fprintf(conf.OutWriter, "\nYour entry was successfully edited!\n")
	return nil
}

func editScriptEntry(flags *pflag.FlagSet, scripts []store.Script, index int) ([]store.Script, error) {
	script, err := createNewEntry(flags, scripts[index])
	if err != nil {
		return scripts, fmt.Errorf("failed to create new script entry: %s", err)
	}

	return append(scripts[:index], append(scripts[index+1:], script)...), nil
}

//createNewEntry return a new code Script entry
func createNewEntry(flags *pflag.FlagSet, script store.Script) (store.Script, error) {
	timeNow := time.Now()

	//Get alias
	alias := script.Code.Alias
	aliasTmp, err := flags.GetString("alias")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `alias` flag: %s", err)
	}
	if aliasTmp != "" {
		alias = aliasTmp
	}

	//Get comment
	comment := script.Comment
	commentTmp, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}
	if commentTmp != "" {
		comment = commentTmp
	}

	//Get code
	var code string
	if script.Code.Content != "" {
		code = script.Code.Content
	}
	codeTmp, err := flags.GetString("code")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `code` flag: %s", err)
	}
	if codeTmp != "" {
		code = codeTmp
	}

	//Get Solution
	var solution string
	if script.Solution.Content != "" {
		solution = script.Solution.Content
	}
	solutionTmp, err := flags.GetString("solution")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `solution` flag: %s", err)
	}
	if solutionTmp != "" {
		solution = solutionTmp
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID(code, comment, alias, solution)
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	if solution != "" {
		return store.Script{
			ID:           id,
			Comment:      comment,
			CreationTime: timeNow,
			UpdateTime:   timeNow,
			Solution: store.Solution{
				Content: solution,
			},
			Code: store.Code{},
		}, nil
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
		Solution: store.Solution{},
	}, nil
}
