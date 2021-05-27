package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/files"
	"github.com/elhmn/ckp/internal/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddHistoryCommand(t *testing.T) {
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf, mockedExec := createConfig()

		//set history files to fixtures
		origBashHistoryFile := history.BashHistoryFile
		origZshHistoryFile := history.ZshHistoryFile
		history.BashHistoryFile = "bash_history_test"
		history.ZshHistoryFile = "zsh_history_test"

		//create bash_history fixtures
		err := files.CopyFileToHomeDirectory(history.BashHistoryFile, "../fixtures/history/bash_history_test")
		if err != nil {
			t.Error(err)
		}

		//create zsh_history fixtures
		err = files.CopyFileToHomeDirectory(history.ZshHistoryFile, "../fixtures/history/zsh_history_test")
		if err != nil {
			t.Error(err)
		}

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		commandName := "history"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName})

		if err := command.Execute(); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour history was successfully added!\n"
		assert.Equal(t, exp, got)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Delete history tmp files
		if err := files.DeleteFileFromHomeDirectory(history.BashHistoryFile); err != nil {
			t.Error(err)
		}
		if err := files.DeleteFileFromHomeDirectory(history.ZshHistoryFile); err != nil {
			t.Error(err)
		}

		//function call assert
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "fetch", "origin", "main")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "diff", "origin/main", "--", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "stash", "apply")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "add", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "commit", "-m", "ckp: add entry")
		mockedExec.AssertCalled(t, "DoGitPush", mock.Anything, "origin", "main")

		//Restore history path
		history.BashHistoryFile = origBashHistoryFile
		history.ZshHistoryFile = origZshHistoryFile
	})
}
