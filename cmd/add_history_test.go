package cmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/files"
	"github.com/elhmn/ckp/internal/history"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddHistoryCommand(t *testing.T) {
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

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
			fmt.Printf("Failed to remove file: %s\n", err)
		}

		if err := files.DeleteFileFromHomeDirectory(history.ZshHistoryFile); err != nil {
			fmt.Printf("Failed to remove file: %s\n", err)
		}

		//Restore history path
		history.BashHistoryFile = origBashHistoryFile
		history.ZshHistoryFile = origZshHistoryFile
	})
}
