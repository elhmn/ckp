package cmd_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
)

//TestInitCommand test the `ckp init` command
func TestInitCommand(t *testing.T) {
	fakeRemoteFolder := "https://github.com/elhmn/fakefolder"

	t.Run("initialised successfully", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().CreateFolderIfDoesNotExist(gomock.Any()),
			mockedExec.EXPECT().DoGitClone(gomock.Any(), gomock.Any(), gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "log"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "branch", "-M", conf.MainBranch),
			mockedExec.EXPECT().DoGitPush(gomock.Any(), "origin", conf.MainBranch, "-f"),
		)

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		command := cmd.NewInitCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("failed to create folder", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().CreateFolderIfDoesNotExist(gomock.Any()).Return(fmt.Errorf("failed to create folder")),
		)

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		exp := "failed to create folder"

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

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("failed to clone remote repository", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().CreateFolderIfDoesNotExist(gomock.Any()),
			mockedExec.EXPECT().DoGitClone(gomock.Any(), gomock.Any(), gomock.Any()).Return("", fmt.Errorf("failed to clone remote repository")),
		)

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		exp := "failed to clone remote repository"

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

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
