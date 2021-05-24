package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSolutionCommand(t *testing.T) {
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

		commandName := "solution"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName,
			"our solution",
			"--path", "filepath",
			"--comment", "a_comment",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour solution was successfully added!\n"
		assert.Equal(t, exp, got)
		mockedExec.AssertExpectations(t)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
