package cmd_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/mocks"
	"github.com/stretchr/testify/mock"
)

//TestInitCommand test the `ckp init` command
func TestInitCommand(t *testing.T) {
	fakeRemoteFolder := "https://github.com/elhmn/fakefolder"

	t.Run("initialised successfully", func(t *testing.T) {
		conf := config.NewDefaultConfig()

		//Mock functions
		mockedExec := &mocks.IExec{}
		mockedExec.On("CreateFolderIfDoesNotExist", mock.Anything).Return(nil)
		mockedExec.On("DoGitClone", mock.Anything, "https://github.com/elhmn/fakefolder", "repo").Return(mock.Anything, nil)
		conf.Exec = mockedExec

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewInitCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}
		mockedExec.AssertExpectations(t)
	})

	t.Run("failed to create folder", func(t *testing.T) {
		conf := config.NewDefaultConfig()
		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		exp := "failed to create folder"

		//Setup for failure
		mockedExec := &mocks.IExec{}
		mockedExec.On("CreateFolderIfDoesNotExist", mock.Anything).Return(fmt.Errorf(exp))
		conf.Exec = mockedExec

		command := cmd.NewInitCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		if !strings.Contains(got, exp) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}

		mockedExec.AssertExpectations(t)
	})

	t.Run("failed to clone remote repository", func(t *testing.T) {
		conf := config.NewDefaultConfig()
		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		exp := "failed to clone remote repository"

		//Setup for failure
		mockedExec := &mocks.IExec{}
		conf.Exec = mockedExec
		mockedExec.On("CreateFolderIfDoesNotExist", mock.Anything).Return(nil)
		mockedExec.On("DoGitClone", mock.Anything, "https://github.com/elhmn/fakefolder", "repo").Return(mock.Anything, fmt.Errorf(exp))

		command := cmd.NewInitCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		if !strings.Contains(got, exp) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}

		mockedExec.AssertExpectations(t)
	})
}
