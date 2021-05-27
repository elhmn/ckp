package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//TestPushCommand test the `ckp push` command
func TestPushCommand(t *testing.T) {
	getMockedExec := func() *mocks.IExec {
		mockedExec := &mocks.IExec{}
		//Setup function calls mocks
		mockedExec.On("DoGit", mock.Anything, mock.Anything).Return(mock.Anything, nil)
		mockedExec.On("DoGit", mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil)
		mockedExec.On("DoGit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil)
		mockedExec.On("DoGit", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil)
		mockedExec.On("DoGitPush", mock.Anything, "origin", mock.Anything).Return(mock.Anything, nil)

		return mockedExec
	}

	t.Run("push successfully", func(t *testing.T) {
		conf := config.NewDefaultConfig()
		mockedExec := getMockedExec()
		conf.Exec = mockedExec
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewPushCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		exp := "\nckp store was pushed successfully\n"
		got := writer.String()
		if !assert.Equal(t, exp, got) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}

		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "fetch", "origin", "main")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "diff", "origin/main", "--", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "stash", "apply")
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "add", mock.Anything)
		mockedExec.AssertCalled(t, "DoGit", mock.Anything, "commit", "-m", "update: update store")
		mockedExec.AssertCalled(t, "DoGitPush", mock.Anything, "origin", "main")
	})
}
