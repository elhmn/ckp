package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func createConfig(t *testing.T) (config.Config, *mocks.MockIExec) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	conf := config.NewDefaultConfig()
	mockedExec := mocks.NewMockIExec(mockCtrl)
	conf.Exec = mockedExec

	//Think of deleting this file later on
	conf.CKPDir = ".ckp_test"
	return conf, mockedExec
}

func getTempStorageFolder(conf config.Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	return fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder), nil
}

func setupFolder(conf config.Config) error {
	if err := deleteFolder(conf); err != nil {
		return fmt.Errorf("Error: failed to delete folder: %s", err)
	}

	folder, err := getTempStorageFolder(conf)
	if err != nil {
		return fmt.Errorf("failed to get temporary storage folder: %s", err)
	}

	if err = os.MkdirAll(folder, 0777); err != nil {
		return err
	}

	return nil
}

func deleteFolder(conf config.Config) error {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	return os.RemoveAll(fmt.Sprintf("%s/%s", home, conf.CKPDir))
}

func TestAddCodeCommand(t *testing.T) {
	t.Run("make sure that it runs successfully", func(t *testing.T) {
		conf, mockedExec := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "pull", "--rebase", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "add", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "commit", "-m", "ckp: add store"),
		)

		commandName := "code"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName,
			"echo \"je suis con\"",
			"--path", "filepath",
			"--comment", "a_comment",
			"--alias", "an_alias",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour code was successfully added!\n"
		assert.Equal(t, exp, got)

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
