package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/elhmn/ckp/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func newInitCommand(conf config.Config) *cobra.Command {
	initCommand := &cobra.Command{
		Use:   "init <storage_repository>",
		Short: "init initialise ckp storage repository",
		Long: `init will initialise a storage repository

		example: ckp init <https://github.com/elhmn/solutions>
		This git repository will be used as your storage
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			storageFolder := args[0]

			if err := initCommand(conf, storageFolder); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return initCommand
}

func initCommand(conf config.Config, remoteStorageFolder string) error {
	//Setup spinner
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	//Create ckp folder if it does not exist
	dir := fmt.Sprintf("%s/%s", home, conf.CKPDir)
	err = conf.Exec.CreateFolderIfDoesNotExist(dir)
	if err != nil {
		return fmt.Errorf("failed to create `%s` directory: %s", dir, err)
	}
	fmt.Fprintf(conf.OutWriter, "`%s` directory was created\n", dir)

	//clone remote storage folder
	fmt.Fprintf(conf.OutWriter, "Initialising `%s` remote storage folder\n", remoteStorageFolder)
	output, err := conf.Exec.DoGitClone(dir, remoteStorageFolder, conf.CKPStorageFolder)
	if err != nil {
		return fmt.Errorf("failed to clone `%s`: %s\n%s", remoteStorageFolder, err, output)
	}
	fmt.Fprintf(conf.OutWriter, "`%s` remote storage folder, Initialised\n", remoteStorageFolder)

	fmt.Fprintf(conf.OutWriter, "ckp successfully initialised\n")
	return nil
}
