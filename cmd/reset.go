package cmd

import (
	"fmt"
	"os"

	"github.com/elhmn/ckp/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

//NewResetCommand create new cobra command for the reset command
func NewResetCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "reset",
		Short: "removes the current remote storage repository",
		Long: `removes the current remote storage repository

		example: ckp reset
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := resetCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func resetCommand(conf config.Config) error {
	//Setup spinner
	conf.Spin.Start()
	defer conf.Spin.Stop()

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	//Delete the temporary file
	dir := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remote storage %s: %s", dir, err)
	}

	fmt.Fprintf(conf.OutWriter, "ckp was successfully reset\n")
	return nil
}
