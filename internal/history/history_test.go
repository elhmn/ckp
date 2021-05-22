package history_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/elhmn/ckp/internal/history"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func createTempFileFromFixture(filepath, fixtureFilePath string) error {
	fixtureData, err := ioutil.ReadFile(fixtureFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s data: %s", fixtureFilePath, err)
	}

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}
	filepath = fmt.Sprintf("%s/%s", home, filepath)

	//Copy the store file to a temporary destination
	if err := ioutil.WriteFile(filepath, fixtureData, 0666); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", filepath, err)
	}

	return nil
}

func deleteFile(filepath string) error {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	return os.Remove(fmt.Sprintf("%s/%s", home, filepath))
}

func TestGetHistoryRecords(t *testing.T) {
	//set history files to fixtures
	origBashHistoryFile := history.BashHistoryFile
	origZshHistoryFile := history.ZshHistoryFile
	history.BashHistoryFile = "bash_history_test"
	history.ZshHistoryFile = "zsh_history_test"

	err := createTempFileFromFixture(history.BashHistoryFile, "../../fixtures/history/bash_history_test")
	if err != nil {
		t.Error(err)
	}

	err = createTempFileFromFixture(history.ZshHistoryFile, "../../fixtures/history/zsh_history_test")
	if err != nil {
		t.Error(err)
	}

	t.Run("return history records", func(t *testing.T) {
		exp := []string{
			"exit",
			"echo \"je suis con\"; echo \"tu es con\"",
			"pwd",
			"cd",
			"cat /dev/null > ~/.bash_history \\",
			"history",
			"echo \"je suis con\"; echo \"tu es con\"",
			"echo \"je suis con end of line\"",
		}
		records, err := history.GetHistoryRecords()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, exp, records)
	})

	//Delete history tmp files
	if err := deleteFile(history.BashHistoryFile); err != nil {
		t.Error(err)
	}
	if err := deleteFile(history.ZshHistoryFile); err != nil {
		t.Error(err)
	}

	//Restore history path
	history.BashHistoryFile = origBashHistoryFile
	history.ZshHistoryFile = origZshHistoryFile
}
