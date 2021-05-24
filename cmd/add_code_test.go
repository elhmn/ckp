package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/mocks"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createConfig() (config.Config, *mocks.IExec) {
	conf := config.NewDefaultConfig()
	mockedExec := &mocks.IExec{}
	conf.Exec = mockedExec
	//Think of deleting this file later on
	conf.CKPDir = ".ckp_test"
	return conf, mockedExec
}

func setupFolder(conf config.Config) error {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder), 0777); err != nil {
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
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf, mockedExec := createConfig()

		//Setup function calls mocks
		mockedExec.On("DoGit", mock.Anything, "diff").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "diff", mock.Anything).Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "stash").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "stash", "apply").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "pull", "--rebase", "origin", "master").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "add", mock.Anything).Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "commit", "-m", mock.Anything).Return(mock.Anything, nil).Once()
		mockedExec.On("DoGitPush", mock.Anything, "origin", "master").Return(mock.Anything, nil).Once()

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

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
		mockedExec.AssertExpectations(t)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
