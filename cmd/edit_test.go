package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEditCommand(t *testing.T) {
	t.Run("make sure that it runs successfully for code edition", func(t *testing.T) {
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

		command := cmd.NewEditCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		codeID := "hash-of-file-content"
		command.SetArgs([]string{codeID,
			"--code", "a_code",
			"--comment", "a_comment",
			"--alias", "an_alias",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour entry was successfully edited!\n"
		assert.Contains(t, got, exp)

		//Test that the store was correctly edited
		filePath, _ := config.GetStoreFilePath(conf)
		s, _, err := store.LoadStore(filePath)
		if err != nil {
			t.Errorf("Error: failed to load the store with %s", err)
		}
		script := s.Scripts[1]
		assert.Equal(t, "a_code", script.Code.Content)
		assert.Equal(t, "a_comment", script.Comment)
		assert.Equal(t, "an_alias", script.Code.Alias)

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure that it runs successfully for solution edition", func(t *testing.T) {
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

		command := cmd.NewEditCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		solutionID := "hash-of-file-content-2"
		command.SetArgs([]string{solutionID,
			"--solution", "a_solution",
			"--comment", "a_comment",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour entry was successfully edited!\n"
		assert.Contains(t, got, exp)

		//Test that the store was correctly edited
		filePath, _ := config.GetStoreFilePath(conf)
		s, _, err := store.LoadStore(filePath)
		if err != nil {
			t.Errorf("Error: failed to load the store with %s", err)
		}
		script := s.Scripts[1]
		assert.Equal(t, "a_solution", script.Solution.Content)
		assert.Equal(t, "a_comment", script.Comment)

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
