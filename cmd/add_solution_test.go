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

		//function call assert
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "fetch", "origin", "main")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "diff", "origin/main", "--", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "stash", "apply")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "add", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "commit", "-m", "ckp: add entry")
		mockedExec.AssertCalled(t, "DoGitPush", mock.Anything, "origin", "main")

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
