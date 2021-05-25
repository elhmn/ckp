package history_test

import (
	"testing"

	"github.com/elhmn/ckp/internal/files"
	"github.com/elhmn/ckp/internal/history"
	"github.com/stretchr/testify/assert"
)

func TestGetHistoryRecords(t *testing.T) {
	//set history files to fixtures
	origBashHistoryFile := history.BashHistoryFile
	origZshHistoryFile := history.ZshHistoryFile
	history.BashHistoryFile = "bash_history_test"
	history.ZshHistoryFile = "zsh_history_test"

	//create bash_history fixtures
	err := files.CopyFileToHomeDirectory(history.BashHistoryFile, "../../fixtures/history/bash_history_test")
	if err != nil {
		t.Error(err)
	}

	//create zsh_history fixtures
	err = files.CopyFileToHomeDirectory(history.ZshHistoryFile, "../../fixtures/history/zsh_history_test")
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
	if err := files.DeleteFileFromHomeDirectory(history.BashHistoryFile); err != nil {
		t.Error(err)
	}
	if err := files.DeleteFileFromHomeDirectory(history.ZshHistoryFile); err != nil {
		t.Error(err)
	}

	//Restore history path
	history.BashHistoryFile = origBashHistoryFile
	history.ZshHistoryFile = origZshHistoryFile
}
