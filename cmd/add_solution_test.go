package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddSolutionCommand(t *testing.T) {
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
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

		commandName := "solution"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName,
			"our solution",
			"--comment", "a_comment",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour solution was successfully added!\n"
		assert.Contains(t, got, exp)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure that is runs successfully without solution argument", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
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
			mockedExec.EXPECT().OpenEditor(gomock.Any(), gomock.Any()).Return(nil),
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "add", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "commit", "-m", "ckp: add store"),
		)

		commandName := "solution"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName,
			"--comment", "a_comment",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour solution was successfully added!\n"
		assert.Contains(t, got, exp)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
