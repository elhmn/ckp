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
		mockedExec.On("DoGit", mock.Anything, "diff").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "diff", mock.Anything).Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "stash").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "stash", "apply").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "pull", "--rebase", "origin", "master").Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "add", mock.Anything).Return(mock.Anything, nil).Once()
		mockedExec.On("DoGit", mock.Anything, "commit", "-m", mock.Anything).Return(mock.Anything, nil).Once()

		mockedExec.On("DoGitPush", mock.Anything, "origin", "master").Return(mock.Anything, nil).Once()
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

		mockedExec.AssertExpectations(t)
	})
}
