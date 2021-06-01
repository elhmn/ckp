package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

//NewPullCommand create new cobra command for the push command
func NewPullCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "pull",
		Short: "pulls changes from remote storage repository",
		Long: `pulls changes from remote storage repository

		example: ckp pull
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := pullCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func pullCommand(conf config.Config) error {
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

	fmt.Fprintf(conf.OutWriter, "\nckp store was pulled successfully\n")
	return nil
}
