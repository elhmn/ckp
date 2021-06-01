package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elhmn/ckp/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

//NewInitCommand create new cobra command for the init command
func NewInitCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "init <storage_repository>",
		Short: "initialise ckp storage repository",
		Long: `will initialise a storage repository

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

	return command
}

func initCommand(conf config.Config, remoteStorageFolder string) error {
	//Setup spinner
	conf.Spin.Start()
	defer conf.Spin.Stop()

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
	out, err := conf.Exec.DoGitClone(dir, remoteStorageFolder, conf.CKPStorageFolder)
	if err != nil {
		return fmt.Errorf("failed to clone `%s`: %s\n%s", remoteStorageFolder, err, out)
	}

	//Get local storage folder
	localStorageFolder := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder)

	//add an empty commit if the repository has no commits
	{
		out, _ := conf.Exec.DoGit(localStorageFolder, "log")
		if strings.Contains(out, "does not have any commits yet") {
			storeFilePath, err := config.GetStoreFilePath(conf)
			if err != nil {
				return fmt.Errorf("failed get store file path: %s", err)
			}

			//Create the storage file
			if err = ioutil.WriteFile(storeFilePath, []byte{}, 0666); err != nil {
				return fmt.Errorf("failed to write to file %s: %s", storeFilePath, err)
			}

			historyStoreFilePath, err := config.GetHistoryFilePath(conf)
			if err != nil {
				return fmt.Errorf("failed get history store file path: %s", err)
			}

			//Create the history storage file
			if err = ioutil.WriteFile(historyStoreFilePath, []byte{}, 0666); err != nil {
				return fmt.Errorf("failed to write to file %s: %s", historyStoreFilePath, err)
			}

			//Add storage file
			out, err = conf.Exec.DoGit(localStorageFolder, "add", storeFilePath, historyStoreFilePath)
			if err != nil {
				return fmt.Errorf("failed to add changes: %s: %s", err, out)
			}

			//Create first commit
			out, err := conf.Exec.DoGit(localStorageFolder, "commit", "-m", "first commit")
			if err != nil {
				return fmt.Errorf("failed to rename branch to `%s`: %s:%s", conf.MainBranch, err, out)
			}
		}
	}

	//checkout rename branch
	out, err = conf.Exec.DoGit(localStorageFolder, "branch", "-M", conf.MainBranch)
	if err != nil {
		return fmt.Errorf("failed to rename branch to `%s`: %s:%s", conf.MainBranch, err, out)
	}

	//push renamed branch to remote
	out, err = conf.Exec.DoGitPush(localStorageFolder, "origin", conf.MainBranch, "-f")
	if err != nil {
		return fmt.Errorf("failed to push `%s` branch: %s:%s", conf.MainBranch, err, out)
	}

	fmt.Fprintf(conf.OutWriter, "`%s` remote storage folder, Initialised\n", remoteStorageFolder)

	fmt.Fprintf(conf.OutWriter, "ckp successfully initialised\n")
	return nil
}
